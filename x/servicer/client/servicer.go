package client

import (
	"context"

	"github.com/pokt-network/smt"

	"poktroll/relayminer/types"
)

type ServicerClient interface {
	NewBlocks() <-chan types.Block
	SubmitClaim(context.Context, []byte) error
	SubmitProof(
		ctx context.Context,
		closestKey []byte,
		closestValueHash []byte,
		closestSum uint64,
		proof *smt.SparseMerkleProof,
	) error
}
