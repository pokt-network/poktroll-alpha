package client

import (
	"context"
	"encoding/json"
	"fmt"

	cometTypes "github.com/cometbft/cometbft/types"

	"poktroll/utils"
	"poktroll/x/servicer/types"
)

var (
	_              types.Block = &cometBlockWebsocketMsg{}
	errNotBlockMsg             = "expected block websocket msg; got: %s"
)

// cometBlockWebsocketMsg is used to deserialize incoming websocket messages from
// the block subscription. It implements the types.Block interface by loosely
// wrapping cometbft's block type, into which messages are deserialized.
type cometBlockWebsocketMsg struct {
	Block cometTypes.Block `json:"block"`
}

func (blockEvent *cometBlockWebsocketMsg) Height() uint64 {
	return uint64(blockEvent.Block.Height)
}

func (blockEvent *cometBlockWebsocketMsg) Hash() []byte {
	return blockEvent.Block.LastBlockID.Hash.Bytes()
}

// BlocksNotifee implements the respective method on the ServicerClient interface.
func (client *servicerClient) BlocksNotifee() utils.Observable[types.Block] {
	return client.blocksNotifee
}

// LatestBlocks implements the respective method on the ServicerClient interface.
func (client *servicerClient) LatestBlock() types.Block {
	client.latestBlockMutex.RLock()
	defer client.latestBlockMutex.RUnlock()
	// block until we have a block to return
	if client.latestBlock == nil {
		subscription := client.BlocksNotifee().Subscribe()
		<-subscription.Ch()
		subscription.Unsubscribe()
	}
	return client.latestBlock
}

// subscribeToBlocks subscribes to committed blocks using a single websocket
// connection.
func (client *servicerClient) subscribeToBlocks(ctx context.Context) utils.Observable[types.Block] {
	// NB: cometbft event subscription query
	// (see: https://docs.cosmos.network/v0.47/core/events#subscribing-to-events)
	query := "tm.event='NewBlock'"

	blocksNotifee, blocksNotifier := utils.NewControlledObservable[types.Block](nil)
	msgHandler := blocksFactoryHandler(blocksNotifier)
	client.subscribeWithQuery(ctx, query, msgHandler)

	return blocksNotifee
}

// blocksFactoryHandler returns a websocket message handler function which attempts
// to deserialize a block event message & send it over the blocksNotifier channel
// which will cause it to be emitted by the corresponding blocksNotifee observable.
func blocksFactoryHandler(blocksNotifier chan types.Block) messageHandler {
	return func(ctx context.Context, msg []byte) error {
		blockMsg, err := newCometBlockMsg(msg)
		expectedErr := fmt.Errorf(errNotBlockMsg, string(msg))
		switch {
		case err == nil:
		case err.Error() == expectedErr.Error():
			return nil
		case err != nil:
			return fmt.Errorf("failed to parse new block message: %w", err)
		}

		blocksNotifier <- blockMsg
		return nil
	}
}

// newCometBlockMsg attempts to deserialize the given bytes into a comet block.
// if the resulting block has a height of zero, assume the message was not a
// block message and return an errNotBlockMsg error.
func newCometBlockMsg(blockMsgBz []byte) (types.Block, error) {
	blockMsg := new(cometBlockWebsocketMsg)
	if err := json.Unmarshal(blockMsgBz, blockMsg); err != nil {
		return nil, err
	}

	// If msg does not match the expected format then the block's height has a zero value.
	if blockMsg.Block.Header.Height == 0 {
		return nil, fmt.Errorf(errNotBlockMsg, string(blockMsgBz))
	}

	return blockMsg, nil
}
