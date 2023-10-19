package proxy

import (
	"context"
	"log"
	"net/http"

	ws "github.com/gorilla/websocket"

	"poktroll/relayer/client"
	"poktroll/utils"
	serviceTypes "poktroll/x/service/types"
	servicerTypes "poktroll/x/servicer/types"
	sessionTypes "poktroll/x/session/types"
)

type wsProxy struct {
	ctx                   context.Context
	serviceId             *serviceTypes.ServiceId
	serviceForwardingAddr string
	sessionQueryClient    sessionTypes.QueryClient
	client                client.ServicerClient
	relayNotifier         chan *RelayWithSession
	signResponse          responseSigner
	servicerAddress       string
	upgrader              *ws.Upgrader
}

func NewWsProxy(
	serviceId *serviceTypes.ServiceId,
	serviceForwardingAddr string,
	sessionQueryClient sessionTypes.QueryClient,
	client client.ServicerClient,
	relayNotifier chan *RelayWithSession,
	signResponse responseSigner,
	servicerAddress string,
) *wsProxy {
	return &wsProxy{
		serviceId:             serviceId,
		serviceForwardingAddr: serviceForwardingAddr,
		sessionQueryClient:    sessionQueryClient,
		client:                client,
		relayNotifier:         relayNotifier,
		signResponse:          signResponse,
		servicerAddress:       servicerAddress,
		upgrader:              &ws.Upgrader{},
	}
}

func (wsProxy *wsProxy) Start(ctx context.Context, advertisedEndpointUrl string) error {
	wsProxy.ctx = ctx
	return http.ListenAndServe(mustGetHostAddress(advertisedEndpointUrl), wsProxy)
}

// ServeHTTP implements the http.Handler interface; called by http.ListenAndServe().
// it validates the initial HTTP request before upgrading the connection to a websocket connection.
// websocket messaging is most of the time asymmetric (0-n requests to 0-m responses)
// and should have a different work accounting (account for requests and responses separately)
// websocket connections are also long-lived and may last across multiple sessions, so each message
// should be validated against the session were that message occurred.
func (wsProxy *wsProxy) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	appAddress := req.URL.Query().Get("app")

	// validate the http upgrade request
	//if err := validateSessionRequest(&sessionResult.Session, relayRequest); err != nil {
	//	utils.ReplyWithHTTPError(400, err, wr)
	//	return
	//}

	// upgrade the connection to a websocket connection
	clientConn, err := wsProxy.upgrader.Upgrade(wr, req, nil)
	if err != nil {
		log.Printf("failed upgrading connection: %v", err)
		utils.ReplyWithHTTPError(500, err, wr)
		return
	}

	// establish a connection to the service
	// OPTIMIZE: reuse the connection to the service
	serviceConn, _, err := ws.DefaultDialer.Dial(wsProxy.serviceForwardingAddr, nil)
	if err != nil {
		log.Printf("failed dialing service: %v", err)
		utils.ReplyWithWsError(err, clientConn)
		return
	}

	// TODO: closing one of the connections should close the other
	// TODO: handle connection errors with errgoups
	go wsProxy.handleWsClientMessages(wsProxy.ctx, clientConn, serviceConn)
	go wsProxy.handleWsServiceMessages(wsProxy.ctx, clientConn, serviceConn, appAddress)
}

func (wsProxy *wsProxy) handleWsClientMessages(ctx context.Context, clientConn, serviceConn *ws.Conn) {
	for {
		messageType, messageBz, err := clientConn.ReadMessage()
		if err != nil {
			if ws.IsUnexpectedCloseError(err) {
				clientConn.Close()
				serviceConn.Close()
				return
			}
			log.Printf("failed reading client message: %v", err)
			utils.ReplyWithWsError(err, clientConn)
			clientConn.Close()
			serviceConn.Close()
			return
		}

		if err := wsProxy.handleWsRequestMessage(ctx, serviceConn, clientConn, messageBz, messageType); err != nil {
			log.Printf("failed handling request message: %v", err)
		}
	}
}

func (wsProxy *wsProxy) handleWsServiceMessages(ctx context.Context, clientConn, serviceConn *ws.Conn, appAddress string) {
	for {
		messageType, messageBz, err := serviceConn.ReadMessage()
		if err != nil {
			if ws.IsUnexpectedCloseError(err) {
				return
			}
			log.Printf("failed reading service message: %v", err)
			utils.ReplyWithWsError(err, clientConn)
			return
		}

		if err := wsProxy.handleWsResponseMessage(ctx, clientConn, serviceConn, messageBz, messageType, appAddress); err != nil {
			log.Printf("failed handling response message: %v", err)
		}
	}
}

func (wsProxy *wsProxy) handleWsRequestMessage(
	ctx context.Context,
	serviceConn *ws.Conn,
	clientConn *ws.Conn,
	req []byte,
	messageType int,
) error {
	relayRequest, err := newWsRelayRequest(req)
	if err != nil {
		return utils.ReplyWithWsError(err, clientConn)
	}

	// TODO: make sure to not request for session info if block height did not change
	// or better, only if the session changed
	query := &sessionTypes.QueryGetSessionRequest{
		AppAddress:  relayRequest.ApplicationAddress,
		ServiceId:   wsProxy.serviceId.Id,
		BlockHeight: wsProxy.client.LatestBlock(ctx).Height(),
	}

	// INVESTIGATE: get the context instead of creating a new one?
	sessionResult, err := wsProxy.sessionQueryClient.GetSession(ctx, query)
	if err != nil {
		return utils.ReplyWithWsError(err, clientConn)
	}

	// validate the websocket request
	if err := validateSessionRequest(&sessionResult.Session, relayRequest); err != nil {
		return utils.ReplyWithWsError(err, clientConn)
	}

	// send the request to the service without handling the response.
	// as this implies managing message ordering, which should be done by the requester.
	// if the client sends requests without waiting for the response, the service should do the same.
	// if the messages contain ordering information, the service would just pass it along.
	if serviceConn.WriteMessage(messageType, relayRequest.Payload) != nil {
		return utils.ReplyWithWsError(err, clientConn)
	}

	// account for request relaying work
	wsProxy.relayNotifier <- &RelayWithSession{
		Relay:   &servicerTypes.Relay{Req: relayRequest, Res: nil},
		Session: &sessionResult.Session,
	}
	return nil
}

func (wsProxy *wsProxy) handleWsResponseMessage(
	ctx context.Context,
	clientConn *ws.Conn,
	servicerCon *ws.Conn,
	response []byte,
	messageType int,

	// the appAddress is needed to query for the session info.
	// it should never change for a given connection.
	appAddress string,
) error {
	relayResponse, err := newWsRelayResponse(response)
	if err != nil {
		return utils.ReplyWithWsError(err, clientConn)
	}

	// TODO: make sure to not request for session info if block height did not change
	// or better, only if the session changed
	query := &sessionTypes.QueryGetSessionRequest{
		AppAddress:  appAddress,
		ServiceId:   wsProxy.serviceId.Id,
		BlockHeight: wsProxy.client.LatestBlock(ctx).Height(),
	}

	// INVESTIGATE: get the context instead of creating a new one?
	sessionResult, err := wsProxy.sessionQueryClient.GetSession(ctx, query)
	if err != nil {
		return utils.ReplyWithWsError(err, clientConn)
	}

	signature, err := wsProxy.signResponse(relayResponse)
	if err != nil {
		return utils.ReplyWithWsError(err, clientConn)
	}

	relayResponse.Signature = signature
	relayResponse.SessionId = sessionResult.Session.SessionId
	relayResponse.ServicerAddress = wsProxy.servicerAddress

	relayResponseBz, err := relayResponse.Marshal()
	if err != nil {
		return utils.ReplyWithWsError(err, clientConn)
	}

	if err := clientConn.WriteMessage(messageType, relayResponseBz); err != nil {
		return utils.ReplyWithWsError(err, clientConn)
	}

	// account for reply relaying work
	wsProxy.relayNotifier <- &RelayWithSession{
		Relay:   &servicerTypes.Relay{Req: nil, Res: relayResponse},
		Session: &sessionResult.Session,
	}

	return nil
}

func newWsRelayRequest(req []byte) (*servicerTypes.RelayRequest, error) {
	var relayRequest servicerTypes.RelayRequest
	err := relayRequest.Unmarshal(req)
	if err != nil {
		return nil, err
	}

	return &relayRequest, nil
}

func newWsRelayResponse(req []byte) (*servicerTypes.RelayResponse, error) {
	relayResponse := &servicerTypes.RelayResponse{
		Payload: req,
	}
	return relayResponse, nil
}
