package relayhandler

import (
	"net/url"

	ws "github.com/gorilla/websocket"

	svcTypes "poktroll/x/service/types"
	sessionTypes "poktroll/x/session/types"
)

// sessionConnection is a wrapper around a websocket connection
// that also holds the session id when session changes the connection is closed
// and a new connection is established to a potentially new relayer
type sessionConnection struct {
	sessionId string
	conn      *ws.Conn
	// initialMessageBz and initialMessageType is the first message received from the client
	// this may be used when the session changes to resend the subscription message
	initialMessageBz   []byte
	initialMessageType int
}

// captureInitialMessage captures the first message received from the client
// this may be used when the session changes to resend the subscription message
func (sc *sessionConnection) captureInitialMessage(messageBz []byte, messageType int) {
	sc.initialMessageBz = messageBz
	sc.initialMessageType = messageType
}

// getInitialMessage returns the first message received from the client so it can be resent
// to the new relayer when the session changes
func (sc *sessionConnection) getInitialMessage() (msgType int, msgBz []byte) {
	return sc.initialMessageType, sc.initialMessageBz
}

// connectWithSession establishes a new connection to a relayer given a session
// it gets a relayer url from the session, adds the application address as a query param
// since the upgrade request does not have a body to embed the application address into
func (relayHandler *RelayHandler) connectWithSession(
	sessionConn *sessionConnection,
	session *sessionTypes.Session,
	rpcType svcTypes.RPCType,
) error {
	// get a relayer url from the session
	serviceUrlStr := relayHandler.getSessionRelayerUrl(session, rpcType)
	serviceUrl, err := url.Parse(serviceUrlStr)
	if err != nil {
		return err
	}

	// add the application address as a query param for ws connections to identify the application
	// before a websocket RelayRequest is received, otherwise the Relayer has to accept any connection
	// and wait for a RelayRequest to identify the application
	serviceUrlQuery := serviceUrl.Query()
	serviceUrlQuery.Add("app", relayHandler.applicationAddress)
	serviceUrl.RawQuery = serviceUrlQuery.Encode()

	// dial the selected relayer
	conn, _, err := ws.DefaultDialer.Dial(serviceUrl.String(), nil)
	if err != nil {
		return err
	}

	// update the session id and connection to be used by future read/write operations
	sessionConn.sessionId = session.SessionId
	sessionConn.conn = conn

	return nil
}
