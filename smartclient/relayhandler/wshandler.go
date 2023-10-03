package relayhandler

import (
	"crypto/sha256"
	"log"
	"net/http"

	ws "github.com/gorilla/websocket"

	"poktroll/utils"
	"poktroll/x/servicer/types"
)

func (relayHandler *RelayHandler) handleWsRelays(
	w http.ResponseWriter,
	req *http.Request,
	serviceId string,
) {
	upgrader := ws.Upgrader{}
	clientConn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		utils.ReplyWithHTTPError(500, err, w)
		return
	}

	serviceSession := relayHandler.servicesSessions[serviceId]
	sessionSubscription := serviceSession.Subscribe(relayHandler.ctx)
	go func() {
		<-relayHandler.ctx.Done()
		sessionSubscription.Unsubscribe()
		_ = clientConn.Close()
	}()

	ch := sessionSubscription.Ch()
	var serviceConn *ws.Conn
	for range ch {
		relayHandler.sessionMutex.RLock()

		if serviceConn != nil {
			_ = serviceConn.Close()
		}

		serviceUrl := relayHandler.getSessionConnectionInfo(serviceId).Url
		serviceConn, _, err = ws.DefaultDialer.Dial(serviceUrl, nil)
		if err != nil {
			log.Printf("failed dialing service: %v", err)
		}

		relayHandler.sessionMutex.RUnlock()

		go relayHandler.handleWsRelayRequests(clientConn, serviceConn)
		go relayHandler.handleWsRelayResponses(clientConn, serviceConn)
	}

	_ = serviceConn.Close()
}

func (relayHandler *RelayHandler) handleWsRelayRequests(
	clientConn *ws.Conn,
	serviceConn *ws.Conn,
) {
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

		relayHandler.sessionMutex.RLock()
		relayRequest := &types.RelayRequest{
			Payload:            messageBz,
			SessionId:          relayHandler.currentSession.SessionId,
			ApplicationAddress: relayHandler.currentSession.Application.Address,
		}

		relayRequestBz, err := relayRequest.Marshal()
		if err != nil {
			utils.ReplyWithWsError(err, clientConn)
			relayHandler.sessionMutex.RUnlock()
			continue
		}

		relayRequestHash := sha256.Sum256(relayRequestBz)
		relaySig, err := relayHandler.signer.Sign(relayRequestHash)
		if err != nil {
			utils.ReplyWithWsError(err, clientConn)
			relayHandler.sessionMutex.RUnlock()
			continue
		}

		relayRequest.ApplicationSignature = relaySig
		relayRequestBz, err = relayRequest.Marshal()
		if err != nil {
			utils.ReplyWithWsError(err, clientConn)
			relayHandler.sessionMutex.RUnlock()
			continue
		}

		if serviceConn.WriteMessage(messageType, relayRequestBz) != nil {
			utils.ReplyWithWsError(err, clientConn)
			relayHandler.sessionMutex.RUnlock()
			continue
		}

		relayHandler.sessionMutex.RUnlock()
	}
}

func (relayHandler *RelayHandler) handleWsRelayResponses(
	clientConn *ws.Conn,
	serviceConn *ws.Conn,
) {
	for {
		relayHandler.sessionMutex.RLock()

		messageType, messageBz, err := serviceConn.ReadMessage()
		if err != nil {
			if ws.IsUnexpectedCloseError(err) {
				relayHandler.sessionMutex.RUnlock()
				return
			}
			log.Printf("failed reading message: %v", err)
			utils.ReplyWithWsError(err, clientConn)
			relayHandler.sessionMutex.RUnlock()
			return
		}

		var relayResponse types.RelayResponse
		err = relayResponse.Unmarshal(messageBz)
		if err != nil {
			utils.ReplyWithWsError(err, clientConn)
			relayHandler.sessionMutex.RUnlock()
			continue
		}

		if relayResponse.SessionId != relayHandler.currentSession.SessionId {
			utils.ReplyWithWsError(errSessionMismatch, clientConn)
			relayHandler.sessionMutex.RUnlock()
			continue
		}

		sig := relayResponse.ServicerSignature
		relayResponse.ServicerSignature = nil
		relayResponseBz, err := relayResponse.Marshal()
		relayResponseHash := sha256.Sum256(relayResponseBz)
		if !relayHandler.signer.Verify(relayResponseHash, sig) {
			utils.ReplyWithWsError(err, clientConn)
			relayHandler.sessionMutex.RUnlock()
			continue
		}

		if clientConn.WriteMessage(messageType, relayResponse.Payload) != nil {
			utils.ReplyWithWsError(err, clientConn)
			relayHandler.sessionMutex.RUnlock()
			continue
		}
	}
}
