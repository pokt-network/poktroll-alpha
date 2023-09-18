package client

import (
	"encoding/hex"
	"encoding/json"
	"strconv"

	"poktroll/x/servicer/types"
)

type tendermintBlockEvent struct {
	Result tendermintBlockEventData `json:"result"`
	height uint64
	hash   []byte
}

type tendermintBlockEventData struct {
	Data struct {
		Value struct {
			Block struct {
				Header struct {
					Height         string `json:"height"`
					LastCommitHash string `json:"last_commit_hash"`
				} `json:"header"`
			} `json:"block"`
		} `json:"value"`
		Type string `json:"type"`
	} `json:"data"`
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

	if blockEvent.Result == (tendermintBlockEventData{}) {
		return nil, nil
	}

	blockEvent.height, err = strconv.ParseUint(blockEvent.Result.Data.Value.Block.Header.Height, 10, 64)
	if err != nil {
		return nil, err
	}

	blockEvent.hash, err = hex.DecodeString(blockEvent.Result.Data.Value.Block.Header.LastCommitHash)
	if err != nil {
		return nil, err
	}

	return blockEvent, nil
}
