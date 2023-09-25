package relayer

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"

	"poktroll/relayer/miner"
	"poktroll/relayer/proxy"
	"poktroll/relayer/sessionmanager"
	"poktroll/x/servicer/types"
)

const WaitGroupContextKey = "relayer_cmd_wait_group"

type Relayer struct {
	proxy          *proxy.Proxy
	miner          *miner.Miner
	sessionManager *sessionmanager.SessionManager
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

// IMPROVE: we tried this pattern because it seemed to be conventional across
// some cosmos-sdk code. In our use case, it turned out to be problematic. In
// the presence of shared and/or nested dependencies, call order starts to
// matter.
// CONSIDERATION: perhaps the `depinject` cosmos-sdk system or a builder
// pattern would be more appropriate.
// see: https://github.com/cosmos/cosmos-sdk/tree/main/depinject#depinject
func (relayer *Relayer) WithKVStorePath(ctx context.Context, storePath string) *Relayer {
	relayer.sessionManager = sessionmanager.NewSessionManager(ctx, storePath, relayer.servicerClient)
	relayer.miner = miner.NewMiner(sha256.New(), relayer.servicerClient, relayer.sessionManager)
	relayer.miner.MineRelays(ctx, relayer.proxy.Relays())

	return relayer
}

func (relayer *Relayer) WithKey(
	ctx context.Context,
	keyring keyring.Keyring,
	keyName string,
	address string,
	clientCtx client.Context,
	client types.ServicerClient,
	serviceEndpoints map[string][]string,
) *Relayer {
	// IMPROVE: separate configuration from subcomponent construction
	var err error
	relayer.proxy, err = proxy.NewProxy(ctx, keyring, keyName, address, clientCtx, client, serviceEndpoints)

	// yet another reason to avoid this pattern: we have to check for errors
	if err != nil {
		panic(fmt.Errorf("failed constructing Proxy: %v", err))
	}

	return relayer
}
