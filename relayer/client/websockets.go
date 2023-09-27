package client

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// listen blocks on reading messages from a websocket connection.
// IMPORTANT: it is intended to be called from within a go routine.
func (client *servicerClient) listen(ctx context.Context, conn *websocket.Conn, msgHandler messageHandler) {
	wg, haveWaitGroup := ctx.Value(WaitGroupContextKey).(*sync.WaitGroup)
	if haveWaitGroup {
		// Increment the relayer wait group to track this goroutine
		// TODO_CLEANUP: Given that we call this a relayer, we should rename `servicerClient` to `relayerClient`
		wg.Add(1)
	}

	// read and handle messages from the websocket. This loop will exit when the
	// websocket connection is closed and/or returns an error.
	for {
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

// messageHandler is a function that handles a websocket chain-event subscription message.
type messageHandler func(ctx context.Context, msg []byte) error

// TODO_CONSIDERATION: the cosmos-sdk CLI code seems to use a cometbft RPC client
// which includes a `#Subscribe()` method for a similar purpose. Perhaps we could
// replace this custom websocket client with that.
// (see: https://github.com/cometbft/cometbft/blob/main/rpc/client/http/http.go#L110)
// (see: https://github.com/cosmos/cosmos-sdk/blob/main/client/rpc/tx.go#L114)
// subscribeWithQuery subscribes to chain event messages matching the given query,
// via a websocket connection.
// (see: https://pkg.go.dev/github.com/cometbft/cometbft/types#pkg-constants)
// (see: https://docs.cosmos.network/v0.47/core/events#subscribing-to-events)
func (client *servicerClient) subscribeWithQuery(ctx context.Context, query string, msgHandler messageHandler) {
	conn, _, err := websocket.DefaultDialer.Dial(client.wsURL, nil)
	if err != nil {
		panic(fmt.Errorf("failed to connect to websocket: %w", err))
	}

	// TODO_DISCUSS: Should we replace `requestId` with just 
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
	go func() {
		<-ctx.Done()
		log.Println("closing websocket")
		_ = conn.Close()
	}()
}

// getNextRequestId increments and returns the JSON-RPC request ID which should
// be used for the next request. These IDs are expected to be unique (per request).
func (client *servicerClient) getNextRequestId() uint64 {
	client.nextRequestId++
	return client.nextRequestId
}
