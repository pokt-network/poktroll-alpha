package client

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"poktroll/utils"
	"poktroll/x/servicer/types"
)

var _ types.Block = &tendermintBlockEvent{}

type tendermintBlockEvent struct {
	Block  tendermintBlock `json:"block"`
	height uint64
	hash   []byte
}

type tendermintBlock struct {
	Header struct {
		Height         uint64 `json:"height"`
		LastCommitHash string `json:"last_commit_hash"`
	} `json:"header"`
}

func (blockEvent *tendermintBlockEvent) Height() uint64 {
	return blockEvent.height
}

func (blockEvent *tendermintBlockEvent) Hash() []byte {
	return blockEvent.hash
}

func newTendermintBlockEvent(blockEventMessage []byte) (_ types.Block, err error) {
	blockEvent := new(tendermintBlockEvent)
	if err := json.Unmarshal(blockEventMessage, blockEvent); err != nil {
		return nil, err
	}

	if blockEvent.Block == (tendermintBlock{}) {
		return nil, nil
	}

	blockEvent.height = blockEvent.Block.Header.Height
	blockEvent.hash, err = hex.DecodeString(blockEvent.Block.Header.LastCommitHash)
	if err != nil {
		return nil, err
	}

	return blockEvent, nil
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
		block, err := newTendermintBlockEvent(msg)
		if err != nil {
			return fmt.Errorf("skipping due to new block event error: %w", err)
		}

		// If msg does not contain data then block is nil, we can ignore it
		if block == nil {
			return fmt.Errorf("skipping because block is nil")
		}

		log.Printf("new block; height: %d, hash: %x\n", block.Height(), block.Hash())
		blocksNotifier <- block

		return nil
	}
}
