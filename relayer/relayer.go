package relayer

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/pokt-network/smt"

	"poktroll/relayer/miner"
	"poktroll/relayer/proxy"
	sessiontracker "poktroll/relayer/session_tracker"
	"poktroll/x/servicer/types"
)

const WaitGroupContextKey = "relayer_cmd_wait_group"

type Relayer struct {
	proxy          *proxy.Proxy
	miner          *miner.Miner
	sessionTracker *sessiontracker.SessionTracker
	servicerClient types.ServicerClient
}

func NewRelayer() *Relayer {
	return &Relayer{}
}

func (relayer *Relayer) Start() error {
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

// IMPROVE: we tried this pattern because it seemed to be conventional across
// some cosmos-sdk code. In our use case, it turned out to be problematic. In
// the presence of shared and/or nested dependencies, call order starts to
// matter.
// CONSIDERATION: perhaps the `depinject` cosmos-sdk system or a builder
// pattern would be more appropriate.
// see: https://github.com/cosmos/cosmos-sdk/tree/main/depinject#depinject
func (relayer *Relayer) WithKVStorePath(storePath string) *Relayer {
	// IMPROVE: separate configuration from subcomponent construction
	kvStore, err := smt.NewKVStore(storePath)

	if err != nil {
		panic(fmt.Errorf("failed to create KVStore %q: %w", storePath, err))
	}

	miner := miner.NewMiner(sha256.New(), kvStore, relayer.servicerClient)
	miner.MineRelays(relayer.proxy.Relays(), relayer.sessionTracker.ClosedSessions())

	return relayer
}

func (relayer *Relayer) WithKey(keyring keyring.Keyring, keyName string) *Relayer {
	// IMPROVE: separate configuration from subcomponent construction
	relayer.proxy = proxy.NewProxy(log.Default(), keyring, keyName)

	return relayer
}
