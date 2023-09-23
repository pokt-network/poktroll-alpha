package types

import (
	"context"

	"github.com/pokt-network/smt"

	"poktroll/utils"
)

type ServicerClient interface {
	Blocks() utils.Observable[Block]
	GetLatestBlock() Block
	SubmitClaim(context.Context, []byte) error
	SubmitProof(
		ctx context.Context,
		smtRootHash []byte,
		closestKey []byte,
		closestValueHash []byte,
		closestSum uint64,
		proof *smt.SparseMerkleProof,
	) error
}
