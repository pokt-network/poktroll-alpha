package proxy

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"

	"poktroll/utils"
	"poktroll/x/service/types"
	svcTypes "poktroll/x/servicer/types"
	sessionTypes "poktroll/x/session/types"
)

type responseSigner func(*svcTypes.RelayResponse) error

type RelayWithSession struct {
	Relay   *svcTypes.Relay
	Session *sessionTypes.Session
}

type Proxy struct {
	services            []*types.ServiceConfig
	keyring             keyring.Keyring
	keyName             string
	client              svcTypes.ServicerClient
	servicerQueryClient svcTypes.QueryClient
	sessionQueryClient  sessionTypes.QueryClient
	relayNotifier       chan *RelayWithSession
	relayNotifee        utils.Observable[*RelayWithSession]
}

// IMPROVE: be consistent with component configuration & setup.
// (We got burned by the `WithXXX` pattern and just did this for now).
func NewProxy(
	ctx context.Context,
	keyring keyring.Keyring,
	keyName string,
	address string,
	clientCtx client.Context,
) *Proxy {
	servicerQueryClient := svcTypes.NewQueryClient(clientCtx)
	servicerInfo, err := servicerQueryClient.Servicers(ctx, &svcTypes.QueryGetServicersRequest{
		Address: address,
	})
	if err != nil {
		log.Fatal(err)
	}

	proxy := &Proxy{
		services:            servicerInfo.Servicers.Services,
		sessionQueryClient:  sessionTypes.NewQueryClient(clientCtx),
		servicerQueryClient: servicerQueryClient,
		keyring:             keyring,
		keyName:             keyName,
	}

	proxy.relayNotifee, proxy.relayNotifier = utils.NewControlledObservable[*RelayWithSession](nil)

	go proxy.listen()

	return proxy
}

func (proxy *Proxy) Relays() utils.Observable[*RelayWithSession] {
	return proxy.relayNotifee
}

func (proxy *Proxy) listen() {
	// create a proxy for each endpoint of each service
	for _, service := range proxy.services {
		for _, endpoint := range service.Endpoints {
			switch endpoint.RpcType {
			case types.RPCType_JSON_RPC:
				go func(serviceId, url string) {
					// TODO: support https
					// httpProxy should support both JSON-RPC and REST endpoints
					httpProxy := NewHttpProxy(
						// serviceAddr should be sourced from config files/params mapping to the service endpoint
						"localhost:8546",
						proxy.sessionQueryClient,
						proxy.client,
						proxy.relayNotifier,
						proxy.signResponse,
						serviceId,
					)

					if err := http.ListenAndServe(url, httpProxy); err != nil {
						log.Fatal(err)
					}
				}(service.Id.Id, endpoint.Url)
			case types.RPCType_WEBSOCKET:
				go func(serviceId, url string) {
					// TODO: support wss
					websocketProxy := NewWsProxy(
						// serviceAddr should be sourced from config files/params mapping to the service endpoint
						"localhost:8546",
						proxy.sessionQueryClient,
						proxy.client,
						proxy.relayNotifier,
						proxy.signResponse,
						serviceId,
					)

					if err := http.ListenAndServe(url, websocketProxy); err != nil {
						log.Fatal(err)
					}
				}(service.Id.Id, endpoint.Url)
			default:
				log.Fatalf("unsupported rpc type: %v", endpoint.RpcType)
			}
		}
	}
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
	if session.SessionId != relayRequest.SessionId {
		return errors.New("invalid session id")
	}

	return nil
}
