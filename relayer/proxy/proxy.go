package proxy

import (
	"context"
	"fmt"
	"net/url"
	"regexp"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"

	"poktroll/utils"
	"poktroll/x/service/types"
	svcTypes "poktroll/x/servicer/types"
	sessionTypes "poktroll/x/session/types"
)

var urlSchemePresenceRegex = regexp.MustCompile(`^\w{0,25}://`)

type responseSigner func(*svcTypes.RelayResponse) error

type RelayWithSession struct {
	Relay   *svcTypes.Relay
	Session *sessionTypes.Session
}

type Proxy struct {
	advertisedServices  []*types.ServiceConfig
	keyring             keyring.Keyring
	keyName             string
	client              svcTypes.ServicerClient
	servicerQueryClient svcTypes.QueryClient
	sessionQueryClient  sessionTypes.QueryClient
	relayNotifier       chan *RelayWithSession
	relayNotifee        utils.Observable[*RelayWithSession]
	serviceEndpoints    map[string][]string
}

// IMPROVE: be consistent with component configuration & setup.
// (We got burned by the `WithXXX` pattern and just did this for now).
func NewProxy(
	ctx context.Context,
	keyring keyring.Keyring,
	keyName string,
	address string,
	clientCtx client.Context,
	client svcTypes.ServicerClient,
	serviceEndpoints map[string][]string,
) (*Proxy, error) {
	servicerQueryClient := svcTypes.NewQueryClient(clientCtx)
	servicerInfo, err := servicerQueryClient.Servicers(ctx, &svcTypes.QueryGetServicersRequest{
		Address: address,
	})
	if err != nil {
		return nil, err
	}

	proxy := &Proxy{
		advertisedServices:  servicerInfo.Servicers.Services,
		sessionQueryClient:  sessionTypes.NewQueryClient(clientCtx),
		servicerQueryClient: servicerQueryClient,
		keyring:             keyring,
		keyName:             keyName,
		client:              client,
		serviceEndpoints:    serviceEndpoints,
	}

	proxy.relayNotifee, proxy.relayNotifier = utils.NewControlledObservable[*RelayWithSession](nil)
	if err := proxy.listen(); err != nil {
		return nil, err
	}

	return proxy, nil
}

func (proxy *Proxy) Relays() utils.Observable[*RelayWithSession] {
	return proxy.relayNotifee
}

func (proxy *Proxy) listen() error {
	// create a proxy for each endpoint of each service
	for _, advertisedService := range proxy.advertisedServices {
		for i, advertisedEndpoint := range advertisedService.Endpoints {
			switch advertisedEndpoint.RpcType {
			case types.RPCType_JSON_RPC:
				// TODO: support https
				// httpProxy should support both JSON-RPC and REST endpoints
				httpProxy := NewHttpProxy(
					advertisedService.Id,
					proxy.serviceEndpoints[advertisedService.Id.Id][i],
					proxy.sessionQueryClient,
					proxy.client,
					proxy.relayNotifier,
					proxy.signResponse,
				)
				go httpProxy.Start(advertisedEndpoint.Url)
			case types.RPCType_WEBSOCKET:
				// TODO: support wss
				websocketProxy := NewWsProxy(
					advertisedService.Id,
					proxy.serviceEndpoints[advertisedService.Id.Id][i],
					proxy.sessionQueryClient,
					proxy.client,
					proxy.relayNotifier,
					proxy.signResponse,
				)
				go websocketProxy.Start(advertisedEndpoint.Url)
			default:
				return fmt.Errorf("unsupported rpc type: %v", advertisedEndpoint.RpcType)
			}
		}
	}

	// TODO_CONSIDERATION: we may accumulate errors and return them here
	return nil
}

func (proxy *Proxy) signResponse(relayResponse *svcTypes.RelayResponse) error {
	relayResBz, err := relayResponse.Marshal()
	if err != nil {
		return err
	}

	relayResponse.ServicerSignature, _, err = proxy.keyring.Sign(proxy.keyName, relayResBz)
	return nil
}

func validateSessionRequest(session *sessionTypes.Session, relayRequest *svcTypes.RelayRequest) error {
	// TODO: validate relayRequest signature

	// a similar SessionId means it's been generated from the same params
	//if session.SessionId != relayRequest.SessionId {
	//	return errors.New("invalid session id")
	//}

	return nil
}

// mustGetHostAddress strip the protocol from the url and path to get the host address
func mustGetHostAddress(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)

	// this should not error since th url is validated before being committed when staking
	if err != nil {
		panic(fmt.Errorf("invalid on-chain data: %s", err))
	}

	return fmt.Sprintf("%s:%s", parsedURL.Hostname(), parsedURL.Port())
}

// parseURLWithScheme ensures that endpoint URLs contain a scheme to avoid ambiguity when
// parsing. (See: https://pkg.go.dev/net/url#Parse)
func parseURLWithScheme(rawURL string) (*url.URL, error) {
	if !urlSchemePresenceRegex.Match([]byte(rawURL)) {
		return nil, fmt.Errorf("empty scheme in endpoint URL: %s", rawURL)
	}

	return url.Parse(rawURL)
}
