package pokt

import (
	"context"

	"github.com/pokt-network/smt"

	"poktroll/modules"
	"poktroll/runtime/di"
	"poktroll/types"
)

var _ modules.PocketNetworkClient = &rollkitPocketNetworkClient{}

type deps struct {
	// TODO_THIS_COMMIT: some cosmos-sdk client API
	//
	// cosmosClient cosmos.Client  // (for example)
}
type rollkitPocketNetworkClient struct {
	di.ModuleInternals[deps]
}

func (client *rollkitPocketNetworkClient) Resolve(injector *di.Injector, path *[]string) {
	client.ResolveDeps(&deps{
		// cosmosClient: di.Resolve(cosmos.ClientToken, injector, path),
	})

	panic("implement me")
}

func (client *rollkitPocketNetworkClient) CascadeStart() error {
	panic("implement me")
}

func (client *rollkitPocketNetworkClient) Start() error {
	panic("implement me")
}

func (client *rollkitPocketNetworkClient) Stop() {
	panic("implement me")
}

func (client *rollkitPocketNetworkClient) Stake(
	ctx context.Context,
	actor *types.Actor,
	amount uint64,
) <-chan types.Maybe[*types.TxResult] {
	panic("implement me")
}

func (client *rollkitPocketNetworkClient) Unstake(
	ctx context.Context,
	actor *types.Actor,
	amount uint64,
) <-chan types.Maybe[*types.TxResult] {
	panic("implement me")
}

func (client *rollkitPocketNetworkClient) OnNewBlock() <-chan *types.Block {
	panic("implement me")
}

func (client *rollkitPocketNetworkClient) SubmitClaim(
	ctx context.Context,
	// TODO: what type should `claim` be?
	claim []byte,
) <-chan types.Maybe[*types.TxResult] {
	panic("implement me")
}

func (client *rollkitPocketNetworkClient) SubmitProof(
	ctx context.Context,
	// TODO: what type should `claim` be?
	proof *smt.SparseMerkleProof,
) <-chan types.Maybe[*types.TxResult] {
	panic("implement me")
}
