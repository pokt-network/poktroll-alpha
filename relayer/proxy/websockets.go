package proxy

import (
	"context"
	"log"
	"net/http"

	ws "github.com/gorilla/websocket"

	serviceTypes "poktroll/x/service/types"
	servicerTypes "poktroll/x/servicer/types"
	sessionTypes "poktroll/x/session/types"
)

type wsProxy struct {
	serviceId             *serviceTypes.ServiceId
	serviceForwardingAddr string
	sessionQueryClient    sessionTypes.QueryClient
	client                types.ServicerClient
	relayNotifier         chan *RelayWithSession
	signResponse          responseSigner
	upgrader              *ws.Upgrader
}

func NewWsProxy(
	serviceId *serviceTypes.ServiceId,
	serviceForwardingAddr string,
	sessionQueryClient sessionTypes.QueryClient,
	client types.ServicerClient,
	relayNotifier chan *RelayWithSession,
	signResponse responseSigner,
) *wsProxy {
	return &wsProxy{
		serviceId:             serviceId,
		serviceForwardingAddr: serviceForwardingAddr,
		sessionQueryClient:    sessionQueryClient,
		client:                client,
		relayNotifier:         relayNotifier,
		signResponse:          signResponse,
	}
}

func (wsProxy *wsProxy) Start(advertisedEndpointUrl string) error {
	return http.ListenAndServe(mustGetHostAddress(advertisedEndpointUrl), wsProxy)
}

// ServeHTTP implements the http.Handler interface; called by http.ListenAndServe().
// it validates the initial HTTP request before upgrading the connection to a websocket connection.
// websocket messaging is most of the time asymmetric (0-n requests to 0-m responses)
// and should have a different work accounting (account for requests and responses separately)
// websocket connections are also long-lived and may last across multiple sessions, so each message
// should be validated against the session were that message occurred.
func (wsProxy *wsProxy) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	relayRequest, err := newHTTPRelayRequest(req)
	if err != nil {
		log.Printf("failed creating relay request: %v", err)
		replyWithHTTPError(500, err, wr)
		return
	}

	// query for session info to validate http initial request prior upgrading the connection
	query := &sessionTypes.QueryGetSessionRequest{
		AppAddress:  relayRequest.ApplicationAddress,
		ServiceId:   wsProxy.serviceId.Id,
		BlockHeight: wsProxy.client.LatestBlock().Height(),
	}

	// INVESTIGATE: get the context instead of creating a new one?
	sessionResult, err := wsProxy.sessionQueryClient.GetSession(context.TODO(), query)
	if err != nil {
		log.Printf("failed getting session info: %v", err)
		replyWithHTTPError(500, err, wr)
		return
	}

	// validate the http upgrade request
	if err := validateSessionRequest(&sessionResult.Session, relayRequest); err != nil {
		replyWithHTTPError(400, err, wr)
		return
	}

	// upgrade the connection to a websocket connection
	clientConn, err := wsProxy.upgrader.Upgrade(wr, req, nil)
	if err != nil {
		log.Printf("failed upgrading connection: %v", err)
		replyWithHTTPError(500, err, wr)
		return
	}

	// establish a connection to the service
	// OPTIMIZE: reuse the connection to the service
	serviceConn, _, err := ws.DefaultDialer.Dial(wsProxy.serviceForwardingAddr, nil)
	if err != nil {
		log.Printf("failed dialing service: %v", err)
		replyWithWsError(err, clientConn)
		return
	}

	// TODO: closing one of the connections should close the other
	// TODO: handle connection errors with errgoups
	go wsProxy.handleWsClientMessages(clientConn, serviceConn)
	go wsProxy.handleWsServiceMessages(clientConn, serviceConn, relayRequest.ApplicationAddress)
}

func (wsProxy *wsProxy) handleWsClientMessages(clientConn, serviceConn *ws.Conn) error {
	defer clientConn.Close()
	for {
		messageType, messageBz, err := clientConn.ReadMessage()
		if err != nil {
			log.Printf("failed reading message: %v", err)
			return replyWithWsError(err, clientConn)
		}
		if err := wsProxy.handleWsRequestMessage(serviceConn, clientConn, messageBz, messageType); err != nil {
			log.Printf("failed handling request message: %v", err)
			return err
		}
	}
}

func (wsProxy *wsProxy) handleWsServiceMessages(clientConn, serviceConn *ws.Conn, appAddress string) error {
	defer serviceConn.Close()
	for {
		messageType, messageBz, err := serviceConn.ReadMessage()
		if err != nil {
			log.Printf("failed reading message: %v", err)
			return replyWithWsError(err, clientConn)
		}
		if err := wsProxy.handleWsResponseMessage(clientConn, serviceConn, messageBz, messageType, appAddress); err != nil {
			log.Printf("failed handling response message: %v", err)
			return err
		}
	}
}

func (wsProxy *wsProxy) handleWsRequestMessage(
	serviceConn *ws.Conn,
	clientConn *ws.Conn,
	req []byte,
	messageType int,
) error {
	relayRequest, err := newWsRelayRequest(req)
	if err != nil {
		return replyWithWsError(err, clientConn)
	}

	// TODO: make sure to not request for session info if block height did not change
	// or better, only if the session changed
	query := &sessionTypes.QueryGetSessionRequest{
		AppAddress:  relayRequest.ApplicationAddress,
		ServiceId:   wsProxy.serviceId.Id,
		BlockHeight: wsProxy.client.LatestBlock().Height(),
	}

	// INVESTIGATE: get the context instead of creating a new one?
	sessionResult, err := wsProxy.sessionQueryClient.GetSession(context.TODO(), query)
	if err != nil {
		return replyWithWsError(err, clientConn)
	}

	// validate the websocket request
	if err := validateSessionRequest(&sessionResult.Session, relayRequest); err != nil {
		return replyWithWsError(err, clientConn)
	}

	// send the request to the service without handling the response.
	// as this implies managing message ordering, which should be done by the requester.
	// if the client sends requests without waiting for the response, the service should do the same.
	// if the messages contain ordering information, the service would just pass it along.
	if serviceConn.WriteMessage(messageType, req) != nil {
		return replyWithWsError(err, clientConn)
	}

	// account for request relaying work
	//wsProxy.relayNotifier <- &RelayWithSession{
	//	Relay:   &servicerTypes.Relay{Req: relayRequest, Res: nil},
	//	Session: &sessionResult.Session,
	//}
	return nil
}

func (wsProxy *wsProxy) handleWsResponseMessage(
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
		return replyWithWsError(err, clientConn)
	}

	// TODO: make sure to not request for session info if block height did not change
	// or better, only if the session changed
	query := &sessionTypes.QueryGetSessionRequest{
		AppAddress:  appAddress,
		ServiceId:   wsProxy.serviceId.Id,
		BlockHeight: wsProxy.client.LatestBlock().Height(),
	}

	// INVESTIGATE: get the context instead of creating a new one?
	sessionResult, err := wsProxy.sessionQueryClient.GetSession(context.TODO(), query)
	if err != nil {
		return replyWithWsError(err, clientConn)
	}

	if err := wsProxy.signResponse(relayResponse); err != nil {
		return replyWithWsError(err, clientConn)
	}

	// serialized relay signed response and send it to the client
	relayResponseBz, err := relayResponse.Marshal()
	if err != nil {
		return replyWithWsError(err, clientConn)
	}

	if clientConn.WriteMessage(messageType, relayResponseBz) != nil {
		return replyWithWsError(err, clientConn)
	}

	// account for reply relaying work
	wsProxy.relayNotifier <- &RelayWithSession{
		Relay:   &servicerTypes.Relay{Req: nil, Res: relayResponse},
		Session: &sessionResult.Session,
	}

	return nil
}

func newWsRelayRequest(req []byte) (*servicerTypes.RelayRequest, error) {
	relayRequest := &servicerTypes.RelayRequest{}
	if err := relayRequest.Unmarshal(req); err != nil {
		return nil, err
	}

	// HACK: the application address should be populated by the requesting client
	relayRequest.ApplicationAddress = "pokt1mrqt5f7qh8uxs27cjm9t7v9e74a9vvdnq5jva4"
	return relayRequest, nil
}

func newWsRelayResponse(req []byte) (*servicerTypes.RelayResponse, error) {
	relayResponse := &servicerTypes.RelayResponse{}
	if err := relayResponse.Unmarshal(req); err != nil {
		return nil, err
	}
	return relayResponse, nil
}

// reply to the client with a derived error message then return the original error
// TODO: send appropriate error instead of the original error
func replyWithWsError(err error, clientConn *ws.Conn) error {
	replyError := clientConn.WriteMessage(ws.TextMessage, []byte(err.Error()))
	if replyError != nil {
		log.Printf("failed sending error response: %v", replyError)
	}

	return err
}
