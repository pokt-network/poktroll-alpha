package client

import (
	"context"

	"github.com/pokt-network/smt"

	"poktroll/utils"
	"poktroll/x/servicer/types"
)

type ServicerClient interface {
	// BlocksNotifee returns an observable which emits newly committed blocks.
	BlocksNotifee() (blocksNotifee utils.Observable[types.Block])
	// LatestBlock returns the latest block that has been committed.
	LatestBlock() types.Block
	// SubmitClaim sends a claim message with the given SMT root hash as the
	// commitment.
	SubmitClaim(ctx context.Context, sessionId string, smtRootHash []byte) error
	// SubmitProof sends a proof message with the given parameters, to be validated
	// on-chain in exchange for a reward.
	SubmitProof(
		ctx context.Context,
		smtRootHash []byte,
		closestKey []byte,
		closestValueHash []byte,
		closestSum uint64,
		proof *smt.SparseMerkleProof,
	) error
}
