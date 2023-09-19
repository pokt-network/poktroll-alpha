package client

import (
	"encoding/hex"
	"encoding/json"
	"poktroll/x/servicer/types"
)

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

func NewTendermintBlockEvent(blockEventMessage []byte) (_ types.Block, err error) {
	blockEvent := new(tendermintBlockEvent)
	json.Unmarshal(blockEventMessage, blockEvent)

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
