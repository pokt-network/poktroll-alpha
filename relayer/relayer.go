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
	proxy          *proxy.Proxy
	miner          *miner.Miner
	sessionTracker *sessiontracker.SessionTracker
	servicerClient types.ServicerClient
}

func NewRelayer() *Relayer {
	return &Relayer{proxy: proxy.NewProxy(log.Default())}
}

func (relayer *Relayer) Start() error {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	<-sigCh

	return nil
}

func (relayer *Relayer) WithServicerClient(client types.ServicerClient) *Relayer {
	relayer.servicerClient = client

	return relayer
}

func (relayer *Relayer) WithBlocksPerSession(ctx context.Context, blocksPerSession uint32) *Relayer {
	sessionTracker := sessiontracker.NewSessionTracker(ctx, blocksPerSession, relayer.servicerClient.NewBlocks())
	relayer.sessionTracker = sessionTracker

	return relayer
}

func (relayer *Relayer) WithKVStorePath(storePath string) *Relayer {
	kvStore, err := smt.NewKVStore(storePath)

	if err != nil {
		panic(fmt.Errorf("failed to create KVStore %q: %w", storePath, err))
	}

	miner := miner.NewMiner(sha256.New(), kvStore, relayer.servicerClient)
	miner.MineRelays(relayer.proxy.Relays(), relayer.sessionTracker.ClosedSessions())

	return relayer
}
