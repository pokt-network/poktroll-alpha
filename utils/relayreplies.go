package utils

import (
	"fmt"
	"log"
	"net/http"

	ws "github.com/gorilla/websocket"
)

// TODO: send appropriate error instead of the original error
// CONSIDERATION: receive err message format string so we don't loose the context of the error.
func ReplyWithHTTPError(statusCode int, err error, wr http.ResponseWriter) {
	wr.WriteHeader(statusCode)
	clientError := err
	if statusCode == 500 {
		clientError = fmt.Errorf("internal server error")
		log.Printf("internal server error: %v", err)
	}

	if _, replyError := wr.Write([]byte(clientError.Error())); replyError != nil {
		log.Printf("failed sending error response: %v", replyError)
	}
}

// reply to the client with a derived error message then return the original error
// TODO: send appropriate error instead of the original error
func ReplyWithWsError(err error, clientConn *ws.Conn) error {
	replyError := clientConn.WriteMessage(ws.TextMessage, []byte(err.Error()))
	if replyError != nil {
		log.Printf("failed sending error response: %v", replyError)
	}

	return err
}
