package relayhandler

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocdc "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/noot/ring-go"
	"log"
	"net/http"
	"poktroll/smartclient"
	portalTypes "poktroll/x/portal/types"
	"strings"

	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	ring_secp256k1 "github.com/athanorlabs/go-dleq/secp256k1"
	ring_types "github.com/athanorlabs/go-dleq/types"
	"poktroll/smartclient/client"
	"poktroll/utils"
	applicationTypes "poktroll/x/application/types"
	svcTypes "poktroll/x/service/types"
	"poktroll/x/servicer/types"
	sessionTypes "poktroll/x/session/types"
)

var (
	errInvalidProtocol = errors.New("invalid protocol")
	errSessionMismatch = errors.New("session mismatch")
	errNoRelayerUrl    = errors.New("no relayer url")
)

// RelayHandler is a http.Handler that handles relays for a given application/gateway
// it is responsible for:
// - listening for new blocks
// - emitting new sessions when a new block is received
// - listening for new relays
// - use the appropriate protocol to relay the requests
type RelayHandler struct {
	ctx context.Context
	// listenAddr is the address where the relay handler will listen for new relays
	// it handles all types of requests (http, websockets, etc...)
	listenAddr string

	// the current application address
	applicationAddress string

	// applicationQueryClient is the query client for the application module
	// it is used to fetch the services provided by the application
	// in order to listen and serve the appropriate relays
	applicationQueryClient applicationTypes.QueryClient

	// portalQueryclient is the query client for the portal module
	// it is used to fetch the delegated pubkeys for a given application
	// in order to create a ring for signing delegated relays
	portalQueryClient portalTypes.QueryClient

	// sessionQueryClient is the query client for the session module
	// it is used to fetch the session info for a given service, block height and application
	sessionQueryClient sessionTypes.QueryClient

	// accountQueryClient is the query client for the auth module (accounts) to fetch servicers
	// pubkeys given their addresses to verify the servicer signature on the relay response
	accountQueryClient authTypes.QueryClient

	// blockQueryClient is the query client for the block module
	// it is used to fetch the latest block heights to fetch the latest session info
	blockQueryClient *client.BlocksQueryClient

	// currentSessions is a map of service id to the current session info
	currentSessions map[string]*sessionTypes.Session

	// servicesSessions is a map of service id to an observable of the current session info
	// used to notify the relay handler when a new session is created so it can rebind to
	// the corresponding relayers
	servicesSessions map[string]utils.Observable[*sessionTypes.Session]

	// endpointSelectionStrategy is the strategy used to select the relayer endpoint given
	// a list of valid endpoints for the session and the RPCType used by the client
	endpointSelectionStrategy EndpointSelectionStrategy

	// Signer is the signer used to sign the relay request
	signer smartclient.Signer

	// signingKey is the private key scalar on the secp256k1 curve used to sign the relay
	// request when using the ring siganture provided the portal is a delegatee of the app
	signingKey ring_types.Scalar
}

func NewRelayHandler(
	listenAddr string,
	applicationQueryClient applicationTypes.QueryClient,
	portalQueryClient portalTypes.QueryClient,
	sessionQueryClient sessionTypes.QueryClient,
	accountQueryClient authTypes.QueryClient,
	blockQueryClient *client.BlocksQueryClient,
	applicationAddress string,
	endpointSelectionStrategy EndpointSelectionStrategy,
	signer Signer,
	signingKey ring_types.Scalar,
) *RelayHandler {
	return &RelayHandler{
		listenAddr:                listenAddr,
		applicationQueryClient:    applicationQueryClient,
		portalQueryClient:         portalQueryClient,
		sessionQueryClient:        sessionQueryClient,
		accountQueryClient:        accountQueryClient,
		blockQueryClient:          blockQueryClient,
		applicationAddress:        applicationAddress,
		currentSessions:           make(map[string]*sessionTypes.Session),
		servicesSessions:          make(map[string]utils.Observable[*sessionTypes.Session]),
		endpointSelectionStrategy: endpointSelectionStrategy,
		signer:                    signer,
		signingKey:                signingKey,
	}
}

func (relayHandler *RelayHandler) Start(ctx context.Context) error {
	relayHandler.ctx = ctx

	// get services supported by the application
	applicationRequest := &applicationTypes.QueryGetApplicationRequest{
		Address: relayHandler.applicationAddress,
	}
	appInfo, err := relayHandler.applicationQueryClient.Application(ctx, applicationRequest)
	if err != nil {
		return err
	}

	// create a session notifier for each service
	// fetches the initial session info for them so the relay handler can start serving relays
	servicesActiveSession := map[string]*sessionNotifier{}
	for _, service := range appInfo.Application.Services {
		sessionsNotifee, sessionsNotifier := utils.NewControlledObservable[*sessionTypes.Session](nil)
		relayHandler.servicesSessions[service.Id] = sessionsNotifee

		sessionInfo, err := relayHandler.fetchCurrentSession(ctx, service.Id)
		if err != nil {
			log.Printf("could not create a session notifier for service %s: %v", service.Id, err)
			continue
		}

		servicesActiveSession[service.Id] = &sessionNotifier{
			Notifier: sessionsNotifier,
			Session:  sessionInfo,
		}

		relayHandler.currentSessions[service.Id] = sessionInfo
	}

	go relayHandler.providedServicesSessionsListener(ctx, servicesActiveSession)

	// start listening for relays, this is a regular http server that uses ServeHTTP
	// implemented by the relayHandler to handle http/ws relays
	go http.ListenAndServe(relayHandler.listenAddr, relayHandler)

	return nil
}

// providedServicesSessionsListener listens for new blocks and updates the session info for each service
// if session last block is past the current block, it notifies the corresponding session notifier
func (relayHandler *RelayHandler) providedServicesSessionsListener(
	ctx context.Context,
	sessionsNotifiers map[string]*sessionNotifier,
) {
	newBlocks := relayHandler.blockQueryClient.BlocksNotifee().Subscribe(ctx).Ch()
	for block := range newBlocks {
		for serviceId, activeSession := range sessionsNotifiers {
			if block.Height() <= getLastSessionBlock(activeSession.Session) {
				continue
			}

			sessionInfo, err := relayHandler.fetchCurrentSession(ctx, serviceId)
			// if we cant fetch the session info, we stop relaying for this service
			// but do not stop the whole the relay handler
			if err != nil {
				log.Printf("could not update session notifier for service %s: %v", serviceId, err)
				sessionsNotifiers[serviceId] = nil
				continue
			}

			activeSession.Session = sessionInfo
			activeSession.Notifier <- sessionInfo
			relayHandler.currentSessions[serviceId] = sessionInfo
		}
	}
}

// ServeHTTP is the http.Handler implementation for the relay handler
// it infers the service and the protocol from the request path and scheme
// requests should be in the form of <protocol>://host:port/serviceId
func (relayHandler *RelayHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	protocol := req.URL.Scheme
	serviceId := strings.Split(path, "/")[1]

	// protocol seems to be always empty, so we infer it from the request headers
	if protocol == "" {
		if req.Header.Get("Upgrade") == "websocket" {
			protocol = "ws"
		} else {
			protocol = "http"
		}
	}

	if protocol == "http" {
		relayHandler.handleHTTPRelays(w, req, serviceId, svcTypes.RPCType_JSON_RPC)
	} else if protocol == "ws" {
		relayHandler.handleWsRelays(w, req, serviceId, svcTypes.RPCType_WEBSOCKET)
	} else {
		// we inform the client about his bad request that assumes an unsupported protocol
		utils.ReplyWithHTTPError(400, errInvalidProtocol, w)
		return
	}
}

// fetchCurrentSession fetches the current session info for a given service application address and latest block
func (relayHandler *RelayHandler) fetchCurrentSession(ctx context.Context, serviceId string) (*sessionTypes.Session, error) {
	query := &sessionTypes.QueryGetSessionRequest{
		BlockHeight: relayHandler.blockQueryClient.LatestBlock(ctx).Height(),
		AppAddress:  relayHandler.applicationAddress,
		ServiceId:   serviceId,
	}
	currentSession, err := relayHandler.sessionQueryClient.GetSession(ctx, query)
	if err != nil {
		return nil, err
	}

	return &currentSession.Session, nil
}

// getSessionRelayerUrl returns the relayer url for a given service and rpc type
// it waits for the session info to be available if it is not already but does not fetch it
func (relayHandler *RelayHandler) getServiceCurrentSession(serviceId string) *sessionTypes.Session {
	// if we dont have a session for this service, we wait for it
	if relayHandler.currentSessions[serviceId] == nil {
		subscription := relayHandler.servicesSessions[serviceId].Subscribe(relayHandler.ctx)
		// block until we get a session
		session := <-subscription.Ch()
		// now we have a session, we unsubscribe from the observable
		subscription.Unsubscribe()
		return session
	}

	return relayHandler.currentSessions[serviceId]
}

// getSessionRelayerUrl returns the relayer url for a given service and rpc type
// the available endpoints are passed to the endpoint selection strategy to select the appropriate one
func (relayHandler *RelayHandler) getSessionRelayerUrl(session *sessionTypes.Session, rpcType svcTypes.RPCType) string {
	endpoints := getSessionEndpoints(session, rpcType)
	endpoint := relayHandler.endpointSelectionStrategy.SelectEndpoint(endpoints)
	return endpoint.Url
}

// UpdateSinger returns the RingSinger implementation of the Signer interface
// used to sign delegated relays on behalf of an application
func (relayHandler *RelayHandler) UpdateSinger() error {
	ring, err := relayHandler.getRingForAddress(relayHandler.applicationAddress)
	if err != nil {
		return nil, err
	}
	relayHandler.signer = smartclient.NewRingSigner(ring, relayHandler.signingKey)
	return nil
}

// getRingForAddress returns the ring used to sign a message for the given application
// address, by querying the portal module for it's delegated pubkeys
func (relayerHandler *RelayHandler) getRingForAddress(address string) (*ring.Ring, error) {
	// get application public key
	appPubKeyReq := &authTypes.QueryAccountRequest{Address: address}
	appPubKeyRes, err := relayerHandler.accountQueryClient.Account(relayerHandler.ctx, appPubKeyReq)
	if err != nil {
		return nil, fmt.Errorf("unable to get applications account: %s [%w]", address, err)
	}
	acc := new(authTypes.BaseAccount)
	if err := acc.Unmarshal(appPubKeyRes.Account.Value); err != nil {
		return nil, fmt.Errorf("unable to deserialise applications account: %s [%w]", address, err)
	}
	appPubKey := acc.GetPubKey()
	// get delegated pubkeys
	req := &portalTypes.QueryGetDelegatedPortalsRequest{AppAddress: address}
	res, err := relayerHandler.portalQueryClient.GetDelegatedPortals(relayerHandler.ctx, req)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve delegated portals for application: %s [%w]", address, err)
	}
	// convert all delegated portals pub keys and app pub key into a slice
	// where the app pub key is index 0
	pubKeys := make([]cryptotypes.PubKey, len(res.Delegatees.PubKeys)+1) // +1 for app pub key
	pubKeys[0] = appPubKey
	for i, anyKey := range res.Delegatees.PubKeys {
		pubKeys[i+1], err = anyToPubKey(anyKey)
		if err != nil {
			return nil, fmt.Errorf("unable to convert codectypes.Any into a cosmos.crypto.PubKey: %w", err)
		}
	}
	// convert the pubkeys to points on the secp256k1 curve
	points, err := pubKeysToPoints(pubKeys)
	if err != nil {
		return nil, fmt.Errorf("unable to convert public keys to points on the secp256k1 curve: %w", err)
	}
	// return the ring for these pubkeys
	return ring.NewFixedKeyRingFromPublicKeys(ring_secp256k1.NewCurve(), points)
}

// pubKeysToPoints converts a slice of cosmos.crypto.PubKey to a slice of points on the secp256k1 curve
// NOTE: Assumes the public keys are secp256k1 public keys unexpected behaviour if not
func pubKeysToPoints(keys []cryptotypes.PubKey) ([]ring_types.Point, error) {
	curve := ring_secp256k1.NewCurve()
	points := make([]ring_types.Point, len(keys))
	for i, key := range keys {
		point, err := curve.DecodeToPoint(key.Bytes())
		if err != nil {
			return nil, err
		}
		points[i] = point
	}
	return points, nil
}

// anyToPubKey unmarshals a serialised Any into a cosmos.crypto.PubKey
func anyToPubKey(any codectypes.Any) (cryptotypes.PubKey, error) {
	reg := codectypes.NewInterfaceRegistry()
	cryptocdc.RegisterInterfaces(reg)
	cdc := codec.NewProtoCodec(reg)
	var pub cryptotypes.PubKey
	if err := cdc.UnpackAny(&any, &pub); err != nil {
		return nil, fmt.Errorf("Any type [%+v] is not cryptotypes.PubKey: %w", any, err)
	}
	return pub, nil
}

// getSessionEndpoints returns a slice of valid endpoints for a given session and rpc type
func getSessionEndpoints(session *sessionTypes.Session, rpcType svcTypes.RPCType) []svcTypes.Endpoint {
	serviceEndpoints := []svcTypes.Endpoint{}
	// loop over servicers (service providers)
	for _, servicer := range session.Servicers {
		// get their provided services
		// TODO: only collect services that are supported by the application
		for _, service := range servicer.Services {
			// servicers may provide multiple endpoints for a given service and even of the same service
			for _, endpoint := range service.Endpoints {
				if endpoint.RpcType == rpcType {
					serviceEndpoints = append(serviceEndpoints, endpoint)
				}
			}
		}
	}

	return serviceEndpoints
}

func getLastSessionBlock(session *sessionTypes.Session) uint64 {
	return session.SessionBlockStartHeight + session.NumBlocksPerSession
}

// signRelayRequest signs a relay request using the given signer
// it ensures that the ApplicationSignature field is nil before:
// - marshaling the relay request
// - hashing the marshaled relay request
// - signing the hash
// - returning the signature that may be added back to the relay request by the caller
func signRelayRequest(
	relayRequest *types.RelayRequest,
	signer Signer,
) (signature []byte, err error) {
	relayRequest.ApplicationSignature = nil
	unsignedRelayRequestBz, err := relayRequest.Marshal()
	if err != nil {
		return nil, err
	}

	return signer.Sign(sha256.Sum256(unsignedRelayRequestBz))
}

// getVerifiedRelayResponse verifies the relay response signature against
// the relay response hash it takes the relay response bytes as received from the relayer,
// extracts the signature, re-marshals the relay response, hashes it then verifies the signature
// it puts the signature back in the relay response before returning it properly formed
func getVerifiedRelayResponse(relayResponseBz []byte) (*types.RelayResponse, error) {
	var relayResponse types.RelayResponse
	err := relayResponse.Unmarshal(relayResponseBz)
	if err != nil {
		return nil, err
	}

	signature := relayResponse.ServicerSignature
	relayResponse.ServicerSignature = nil
	relayResponseBz, err = relayResponse.Marshal()
	if err != nil {
		return nil, err
	}

	// TODO: In the current state, the relayHandler do not have access to the servicer public key
	// that should be fetched when new sessions are created and stored in the relayHandler
	//relayResponseHash := sha256.Sum256(relayResponseBz)
	_ = sha256.Sum256(relayResponseBz)

	// verify signature against relayResponseHash

	relayResponse.ServicerSignature = signature

	return &relayResponse, nil
}

// Naive implementation of the endpoint selection strategy that always selects the first endpoint
// We may have other strategies that select endpoints. Round robin or based on their latency, their location, etc.
// It is up to the strategy to maintain the needed data to select the appropriate endpoint
type ChooseFirstEndpoint struct{}

func (c *ChooseFirstEndpoint) SelectEndpoint(endpoints []svcTypes.Endpoint) *svcTypes.Endpoint {
	return &endpoints[0]
}

type sessionNotifier struct {
	Notifier chan *sessionTypes.Session
	Session  *sessionTypes.Session
}
