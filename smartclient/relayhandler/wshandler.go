package relayhandler

import (
	"log"
	"net/http"

	ws "github.com/gorilla/websocket"

	"poktroll/utils"
	svcTypes "poktroll/x/service/types"
	"poktroll/x/servicer/types"
	sessionTypes "poktroll/x/session/types"
)

// handleWsRelays handles websocket relay requests
// it takes the http ServeHTTP arguments, the service id, and the rpc type
// to upgrade the connection to a websocket connection and send it to the relayer
func (relayHandler *RelayHandler) handleWsRelays(
	w http.ResponseWriter,
	req *http.Request,
	serviceId string,
	rpcType svcTypes.RPCType,
) {
	// upgrade the connection to a websocket connection or reply with an http error since
	// we don't have a websocket connection yet
	upgrader := ws.Upgrader{}
	clientConn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		utils.ReplyWithHTTPError(500, err, w)
		return
	}

	// create a sessionConnection that holds the websocket connection and the session id
	// the underlying connection is closed and a new connection is established when
	// the session or the relayer changes
	sessionConn := &sessionConnection{}

	// get the current session at the time of connection
	initialSession := relayHandler.getServiceCurrentSession(serviceId)
	// connect to a relayer with the current session that will be used to infer the relayer url
	// and populate the sessionConnection.conn property
	err = relayHandler.connectWithSession(sessionConn, initialSession, rpcType)
	if err != nil {
		// reply to the client with an error if the connection could not be established
		utils.ReplyWithWsError(err, clientConn)
		log.Printf("failed initial dialing service: %v", err)
		return
	}

	// start a goroutine that will switch the relayer when the session changes
	go relayHandler.switchRelayerOnSessionChange(sessionConn, serviceId, rpcType)

	// start a goroutine that will handle requests from the client and forward them to the relayer
	go relayHandler.handleWsRelayRequests(clientConn, sessionConn, serviceId)

	// start a goroutine that will handle responses from the relayer and forward them to the client
	go relayHandler.handleWsRelayResponses(clientConn, sessionConn, serviceId, rpcType)
}

// switchRelayerOnSessionChange switches the relayer when the session changes
// it listens to the session subscription channel and closes the current connection
// before dialing a new relayer
func (relayHandler *RelayHandler) switchRelayerOnSessionChange(sessionConn *sessionConnection, serviceId string, rpcType svcTypes.RPCType) {
	ch := relayHandler.servicesSessions[serviceId].Subscribe(relayHandler.ctx).Ch()
	for session := range ch {
		// do not change the connection if the session id is the same
		// this should not happen but we ensure avoiding closing and dialing again
		if session.SessionId == sessionConn.sessionId {
			continue
		}

		log.Printf("closing old service connection, %s", sessionConn.sessionId)
		if err := sessionConn.conn.Close(); err != nil {
			log.Printf("failed closing old service connection: %v", err)
			return
		}

		log.Printf("dialing new service %s", session.SessionId)
		err := relayHandler.connectWithSession(sessionConn, session, rpcType)
		if err != nil {
			log.Printf("failed dialing service: %v", err)
			return
		}
	}

	// close the connection if the session subscription channel is closed and have no way
	// to be notified of new sessions
	_ = sessionConn.conn.Close()
}

// handleWsRelayRequests handles websocket relay requests
// it reads messages from the client then hands the message to relayHandler.handleMessage
func (relayHandler *RelayHandler) handleWsRelayRequests(
	clientConn *ws.Conn,
	serviceConn *sessionConnection,
	serviceId string,
) {
	defer clientConn.Close()
	// the first message is special because it may contain the subscription message
	// that is re-sent to the new Relayer when the session changes
	firstMessage := true
	for {
		messageType, messageBz, err := clientConn.ReadMessage()
		if err != nil {
			if ws.IsUnexpectedCloseError(err) {
				return
			}
			log.Printf("failed reading message: %v", err)
			utils.ReplyWithWsError(err, clientConn)
			return
		}

		currentSession := relayHandler.getServiceCurrentSession(serviceId)
		err = relayHandler.handleMessage(
			messageType,
			messageBz,
			serviceId,
			clientConn,
			serviceConn,
			currentSession,
			firstMessage,
		)
		firstMessage = false

		if err != nil {
			log.Printf("failed handling message: %v", err)
			utils.ReplyWithWsError(err, clientConn)
			continue
		}
	}
}

// handleWsRelayResponses handles websocket relay responses
// it reads messages from the relayer then verifies them before forwarding them to the client
func (relayHandler *RelayHandler) handleWsRelayResponses(
	clientConn *ws.Conn,
	serviceConn *sessionConnection,
	serviceId string,
	rpcType svcTypes.RPCType,
) {
	for {
		messageType, messageBz, err := serviceConn.conn.ReadMessage()
		if err != nil {
			if ws.IsUnexpectedCloseError(err) {
				return
			}

			// if the connection is closed, get the current session and use it to dial a relayer
			session := relayHandler.getServiceCurrentSession(serviceId)
			err := relayHandler.connectWithSession(serviceConn, session, rpcType)
			if err != nil {
				// abandon if the connection could not be established, connection retries should be
				// handled by the switchRelayerOnSessionChange goroutine
				log.Printf("failed dialing service: %v", err)
				return
			}

			// since the connection is closed, we need to resend the initial message to the new Relayer
			// this way of doing it may not apply to all websocket protocols. As some may have state altering
			// messages that are not idempotent and should not be re-sent. It also may need to stick to the same
			// relayer if the altered state is not propagated to other Relayers.
			// In that case the connection should be closed and the client should be notified of the error
			// for it to handle the disconnection and reconnection as par the protocol/service it is using
			// TODO: this should be guarded by some query param sent to the relayHandler by the client
			// to indicate that it is safe to resend the initial message
			initialMsgType, initialMsgBz := serviceConn.getInitialMessage()

			// handle the message without passing by serviceConn.conn.ReadMessage as it is blocking and
			// waiting for a message from the new connection that has never received the subscription message
			err = relayHandler.handleMessage(
				initialMsgType,
				initialMsgBz,
				serviceId,
				clientConn,
				serviceConn,
				session,
				false,
			)
			if err != nil {
				log.Printf("failed sending initial message: %v", err)
				return
			}
			continue
		}

		// verify and get the RelayResponse from the message bytes
		relayResponse, err := getVerifiedRelayResponse(messageBz)
		if err != nil {
			continue
		}

		// do not send back the response if the session id is not the same as the websocket message exchange
		// may happen across multiple sessions if the session changes while the client is waiting for a response
		// a subscription to the new Relayer should be sent and the response received from this new one, discarding
		// the response from the old Relayer
		if relayResponse.SessionId != serviceConn.sessionId {
			continue
		}

		// send the original response back to the client
		if clientConn.WriteMessage(messageType, relayResponse.Payload) != nil {
			continue
		}
	}
}

// handleMessage handles a message from the client, it may be the first message
// in that case it is captured and used to resend it to the new Relayer when the session changes
// short circuiting the blocking serviceConn.conn.ReadMessage that has never received the subscription message
func (relayHandler *RelayHandler) handleMessage(
	messageType int,
	messageBz []byte,
	serviceId string,
	clientConn *ws.Conn,
	serviceConn *sessionConnection,
	currentSession *sessionTypes.Session,
	firstMessage bool,
) error {
	relayRequest := &types.RelayRequest{
		Payload:            messageBz,
		SessionId:          currentSession.SessionId,
		ApplicationAddress: currentSession.Application.Address,
	}

	// Sign the RelayRequest with the provided signer
	signature, err := signRelayRequest(relayRequest, relayHandler.signer)
	if err != nil {
		return err
	}

	// append the signature to the RelayRequest and marshal it to bytes to be sent to the relayer
	relayRequest.ApplicationSignature = signature
	relayRequestBz, err := relayRequest.Marshal()
	if err != nil {
		return err
	}

	// do not send the message if the session id is not the same as the one expected by the serviceConnection
	// get the new session info and use it to handle the message again.
	// This call is recursive, if in some way the sessionId does not match the currentSession.SessionId
	// we may end up in an infinite loop or a very deep recursion stack.
	// TODO: Ensure a number of retries and return an error if the session id does not match
	// after the retries are exhausted
	if currentSession.SessionId != serviceConn.sessionId {
		currentSession = relayHandler.getServiceCurrentSession(serviceId)
		return relayHandler.handleMessage(
			messageType,
			messageBz,
			serviceId,
			clientConn,
			serviceConn,
			currentSession,
			firstMessage,
		)
	}

	// if it is the first message of that connection, capture it to be resent when the session changes
	if firstMessage {
		serviceConn.captureInitialMessage(messageBz, messageType)
	}

	// send the RelayRequest to the relayer
	return serviceConn.conn.WriteMessage(messageType, relayRequestBz)
}
