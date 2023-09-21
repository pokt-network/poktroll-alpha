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
			if websocket.IsUnexpectedCloseError(err) {
				// NB: stop this goroutine if the websocket connection is closed
				return
			}
			log.Printf("skipping due to websocket error: %s\n", err)
			// TODO: handle other errors (?)
			continue
		}

		if err := msgHandler(ctx, msg); err != nil {
			log.Printf("skipping due to message handler error: %s\n", err)
			continue
		}
	}
}

type messageHandler func(ctx context.Context, msg []byte) error

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
