package modules

import (
	"context"

	"github.com/pokt-network/smt"

	"poktroll/runtime/di"
	"poktroll/x/poktroll/types"
)

var PocketNetworkClientToken = di.NewInjectionToken[PocketNetworkClient]("pocketNetworkClient")

type PocketNetworkClient interface {
	di.Module
	// Callback listeners anytime a new block is observed
	OnNewBlock() <-chan *types.Block
	Stake(context.Context, *types.Actor, uint64) <-chan types.Maybe[*types.TxResult]
	Unstake(context.Context, *types.Actor, uint64) <-chan types.Maybe[*types.TxResult]
	//SubmitClaim(context.Context, *types.Claim) <-chan types.Maybe[*types.TxResult]
	SubmitClaim(context.Context, []byte) <-chan types.Maybe[*types.TxResult]
	SubmitProof(
		ctx context.Context,
		closestKey []byte,
		closestValueHash []byte,
		closestSum uint64,
		proof *smt.SparseMerkleProof,
		err error,
	) <-chan types.Maybe[*types.TxResult]
}

type ServicerClient interface {
	di.Module
	// HandleRelay is called by end-users' clients to initiate a relay request.
	HandleRelay(context.Context, *types.Relay) <-chan types.Maybe[*types.Relay]
}
