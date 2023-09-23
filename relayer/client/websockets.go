package client

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"poktroll/relayer"
	"sync"
)

// listen blocks on reading messages from a websocket connection, it is intended
// to be called from within a go routine.
func (client *servicerClient) listen(ctx context.Context, conn *websocket.Conn, msgHandler messageHandler) {
	wg, haveWaitGroup := ctx.Value(relayer.WaitGroupContextKey).(*sync.WaitGroup)
	if haveWaitGroup {
		// Increment the relayer wait group to track this goroutine
		wg.Add(1)
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("closing websocket")
			_ = conn.Close()
			if haveWaitGroup {
				// Decrement the wait group as this goroutine stops
				wg.Done()
			}
			return
		default:
		}

		_, msg, err := conn.ReadMessage()
		if err != nil {
			if haveWaitGroup {
				// Decrement the wait group as this goroutine stops
				wg.Done()
			}

			// Stop this goroutine if there's an error.
			//
			// See gorilla websocket `Conn#NextReader()` docs:
			// | Applications must break out of the application's read loop when this method
			// | returns a non-nil error value. Errors returned from this method are
			// | permanent. Once this method returns a non-nil error, all subsequent calls to
			// | this method return the same error.
			return
		}

		if err := msgHandler(ctx, msg); err != nil {
			log.Printf("failed to handle websocket msg: %s\n", err)
			continue
		}
	}
}

type messageHandler func(ctx context.Context, msg []byte) error

// TODO_CONSIDERATION: the cosmos-sdk CLI code seems to use
// subscribeWithQuery subscribes to a websocket connection with the given query,
func (client *servicerClient) subscribeWithQuery(ctx context.Context, query string, msgHandler messageHandler) {
	conn, _, err := websocket.DefaultDialer.Dial(client.wsURL, nil)
	if err != nil {
		panic(fmt.Errorf("failed to connect to websocket: %w", err))
	}

	requestId := client.getNextRequestId()
	conn.WriteJSON(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "subscribe",
		"id":      requestId,
		"params": map[string]interface{}{
			"query": query,
		},
	})

	go client.listen(ctx, conn, msgHandler)
}

func (client *servicerClient) getNextRequestId() uint64 {
	client.nextRequestId++
	return client.nextRequestId
}
