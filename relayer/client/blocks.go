package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

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
	return blockEvent.Block.LastCommitHash.Bytes()
}

func (client *servicerClient) Blocks() utils.Observable[types.Block] {
	return client.blocksNotifee
}

func (client *servicerClient) subscribeToBlocks(ctx context.Context) utils.Observable[types.Block] {
	query := "tm.event='NewBlock'"

	blocksNotifee, blocksNotifier := utils.NewControlledObservable[types.Block](nil)
	msgHandler := handleBlocksFactory(blocksNotifier)
	client.subscribeWithQuery(ctx, query, msgHandler)

	return blocksNotifee
}

func handleBlocksFactory(blocksNotifier chan types.Block) messageHandler {
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

		log.Printf("new blockMsg; height: %d, hash: %x\n", blockMsg.Height(), blockMsg.Hash())
		blocksNotifier <- blockMsg

		return nil
	}
}

func newCometBlockMsg(blockMsgBz []byte) (types.Block, error) {
	blockMsg := new(cometBlockWebsocketMsg)
	if err := json.Unmarshal(blockMsgBz, blockMsg); err != nil {
		return nil, err
	}

	// If msg does not match the expected format then block will be its zero value.
	if blockMsg.Block.Header.Height == 0 {
		return nil, fmt.Errorf(errNotBlockMsg, string(blockMsgBz))
	}

	return blockMsg, nil
}
