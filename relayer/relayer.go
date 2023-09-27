package relayer

import (
	"context"
	"crypto/sha256"
	"fmt"
	"poktroll/relayer/client"
	"poktroll/relayer/miner"
	"poktroll/relayer/proxy"
	"poktroll/relayer/sessionmanager"

	cosmosClient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
)

// TODO: Need to add some comments related to what each field is responsible for
type Relayer struct {
	proxy          *proxy.Proxy
	miner          *miner.Miner
	sessionManager *sessionmanager.SessionManager
	servicerClient client.ServicerClient
}

func NewRelayer() *Relayer {
	return &Relayer{}
}

func (relayer *Relayer) Start() error {
	return nil
}

// IMPROVE: we tried this pattern because it seemed to be conventional across
// some cosmos-sdk code. In our use case, it turned out to be problematic. In
// the presence of shared and/or nested dependencies, call order starts to
// matter.
func (relayer *Relayer) WithServicerClient(client client.ServicerClient) *Relayer {
	relayer.servicerClient = client

	return relayer
}

// TODO_CONSIDERATION: perhaps the `depinject` cosmos-sdk system or a builder
// pattern would be more appropriate.
// see: https://github.com/cosmos/cosmos-sdk/tree/main/depinject#depinject
func (relayer *Relayer) WithKVStorePath(ctx context.Context, storePath string) *Relayer {
	relayer.sessionManager = sessionmanager.NewSessionManager(ctx, storePath, relayer.servicerClient)
	// TODO_REFACTOR: `WithKVStorePath` has a side effect of starting the relay mining process. This
	// should happen in a separate `Start` command while this only deals with options.
	relayer.miner = miner.NewMiner(sha256.New(), relayer.servicerClient, relayer.sessionManager)
	relayer.miner.MineRelays(ctx, relayer.proxy.Relays())

	return relayer
}

func (relayer *Relayer) WithKey(
	ctx context.Context,
	keyring keyring.Keyring,
	keyName string,
	address string,
	clientCtx cosmosClient.Context,
	client client.ServicerClient,
	serviceEndpoints map[string][]string,
) *Relayer {
	// TODO_IMPROVE: separate configuration from subcomponent construction. Starting the proxy should
	// probably happen in `Start` while `withKey` simply updates the state.
	var err error
	relayer.proxy, err = proxy.NewProxy(ctx, keyring, keyName, address, clientCtx, client, serviceEndpoints)

	// TODO_IMPROVE: yet another reason to avoid this pattern: we have to check for errors
	if err != nil {
		panic(fmt.Errorf("failed constructing Proxy: %v", err))
	}

	return relayer
}
