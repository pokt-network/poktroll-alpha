package relayer

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/pokt-network/smt"

	"poktroll/relayer/miner"
	"poktroll/relayer/proxy"
	sessiontracker "poktroll/relayer/session_tracker"
	"poktroll/x/servicer/types"
)

type Relayer struct {
	relayer        *proxy.Proxy
	miner          *miner.Miner
	sessionTracker *sessiontracker.SessionTracker
	newBlocks      chan types.Block
}

func NewRelayer(ctx context.Context, client types.ServicerClient) *Relayer {
	relayer := proxy.NewProxy(log.Default())
	// should be sourced somehow form a subscription to the blockchain
	newBlocks := make(chan types.Block)
	sessionTracker := sessiontracker.NewSessionTracker(ctx, newBlocks)

	storePath := "/tmp/smt"
	kvStore, err := smt.NewKVStore(storePath)

	if err != nil {
		panic(fmt.Errorf("failed to create KVStore %q: %w", storePath, err))
	}

	miner := miner.NewMiner(sha256.New(), kvStore, client)
	miner.MineRelays(relayer.Relays(), sessionTracker.ClosedSessions())

	return &Relayer{
		relayer:        relayer,
		miner:          miner,
		sessionTracker: sessionTracker,
		newBlocks:      newBlocks,
	}
}

func (relayer *Relayer) Start() error {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	<-sigCh

	return nil
}
