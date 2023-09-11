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
	NewBlocks() <-chan *types.Block
	StakeServicer(
		ctx context.Context,
		servicer *types.Servicer,
		amount string,
	) <-chan types.Maybe[*types.TxResult]
	StakeApplication(
		ctx context.Context,
		application *types.Application,
		amount string,
	) <-chan types.Maybe[*types.TxResult]
	UnstakeServicer(
		ctx context.Context,
		servicer *types.Servicer,
		amount string,
	) <-chan types.Maybe[*types.TxResult]
	UnstakeApplication(
		ctx context.Context,
		application *types.Application,
		amount string,
	) <-chan types.Maybe[*types.TxResult]
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
