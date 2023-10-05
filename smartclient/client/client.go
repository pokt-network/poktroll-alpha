package client

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/polydawn/refmt/json"

	relayerClient "poktroll/relayer/client"
	"poktroll/utils"
	"poktroll/x/servicer/types"
)

// BlocksQueryClient is a client for querying blocks from a relayer
// it is inspired from the relayer's client but does not need the transaction
// and signing features thus reducing the number of its dependencies
type BlocksQueryClient struct {
	conn             *websocket.Conn
	blocksNotifee    utils.Observable[types.Block]
	latestBlockMutex *sync.RWMutex
	latestBlock      types.Block
}

func NewBlocksQueryClient(ctx context.Context, endpoint string) (*BlocksQueryClient, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse endpoint: %w", err)
	}

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to dial endpoint: %w", err)
	}

	client := &BlocksQueryClient{
		conn:             conn,
		latestBlockMutex: &sync.RWMutex{},
	}

	if err := client.subscribeToBlocks(ctx); err != nil {
		return nil, err
	}

	return client, nil
}

func (qc *BlocksQueryClient) BlocksNotifee() utils.Observable[types.Block] {
	return qc.blocksNotifee
}

func (qc *BlocksQueryClient) LatestBlock(ctx context.Context) types.Block {
	qc.latestBlockMutex.RLock()
	defer qc.latestBlockMutex.RUnlock()
	// block until we have a block to return
	if qc.latestBlock == nil {
		subscription := qc.blocksNotifee.Subscribe(ctx)
		<-subscription.Ch()
		subscription.Unsubscribe()
	}
	return qc.latestBlock
}

func (qc *BlocksQueryClient) listen(ctx context.Context, blocksNotifier chan types.Block) {
	for {
		_, msg, err := qc.conn.ReadMessage()
		if err != nil {
			return
		}

		blockMsg, err := relayerClient.NewCometBlockMsg(msg)
		expectedErr := fmt.Errorf(relayerClient.ErrNotBlockMsg, string(msg))
		switch {
		case err == nil:
		case err.Error() == expectedErr.Error():
		case err != nil:
			log.Printf("failed to handle websocket msg: %s\n", err)
		}

		qc.latestBlock = blockMsg
		blocksNotifier <- blockMsg
	}
}

func (qc *BlocksQueryClient) subscribeToBlocks(ctx context.Context) error {
	blocksNotifee, blocksNotifier := utils.NewControlledObservable[types.Block](nil)

	subscribeMsg := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "subscribe",
		"id":      1,
		"params": map[string]interface{}{
			"query": "tm.event='NewBlock'",
		},
	}

	subscribeMsgBz, err := json.Marshal(subscribeMsg)
	if err != nil {
		return err
	}

	if err := qc.conn.WriteMessage(websocket.TextMessage, subscribeMsgBz); err != nil {
		return err
	}

	go qc.listen(ctx, blocksNotifier)
	go func() {
		<-ctx.Done()
		log.Println("closing websocket")
		_ = qc.conn.Close()
	}()
	qc.blocksNotifee = blocksNotifee
	return nil
}
