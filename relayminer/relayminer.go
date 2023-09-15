package relayminer

import (
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/pokt-network/smt"

	"poktroll/relayminer/miner"
	"poktroll/relayminer/relayer"
	sessiontracker "poktroll/relayminer/session_tracker"
	"poktroll/relayminer/types"
	"poktroll/x/servicer/client"
)

type RelayMiner struct {
	relayer        *relayer.Relayer
	miner          *miner.Miner
	sessionTracker *sessiontracker.SessionTracker
	newBlocks      chan types.Block
}

func NewRelayMiner(client client.ServicerClient) *RelayMiner {
	relayer := relayer.NewRelayer(log.Default())
	// should be sourced somehow form a subscription to the blockchain
	newBlocks := make(chan types.Block)
	sessionTracker := sessiontracker.NewSessionTracker(newBlocks)

	storePath := "/tmp/smt"
	kvStore, err := smt.NewKVStore(storePath)

	if err != nil {
		panic(fmt.Errorf("failed to create KVStore %q: %w", storePath, err))
	}

	miner := miner.NewMiner(sha256.New(), kvStore, client)
	miner.MineRelays(relayer.Relays(), sessionTracker.ClosedSessions())

	return &RelayMiner{
		relayer:        relayer,
		miner:          miner,
		sessionTracker: sessionTracker,
		newBlocks:      newBlocks,
	}
}

func (relayMiner *RelayMiner) Start() error {
	return nil
}
