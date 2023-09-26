package keeper

import (
	"context"
	"crypto/sha256"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pokt-network/smt"

	"poktroll/x/servicer/types"
)

var errInvalidPathFmt = "invalid path: %x, expected: %x"

func (k msgServer) Proof(goCtx context.Context, msg *types.MsgProof) (*types.MsgProofResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := k.Logger(ctx).With("method", "Proof")
	// INCOMPLETE: (see below)
	//currentBlockHash := ctx.BlockHeader().LastBlockId.Hash

	proof := new(smt.SparseMerkleProof)
	if err := proof.Unmarshal(msg.Proof); err != nil {
		return nil, err
	}

	logger = logger.
		With("servicer_address", msg.ServicerAddress).
		With("smst_root_hash", fmt.Sprintf("%x", msg.SmstRootHash))

	// INCOMPLETE: we need to verify that the closest path matches the last block hash.
	//if proof.VerifyClosest(currentBlockHash) {
	//	err := fmt.Errorf(errInvalidPathFmt, msg.Path, currentBlockHash)
	//	logger.Error(err.Error())
	//	return nil, err
	//}

	// INCOMPLETE: lookup the corresponding claim and verify that it matches.

	if valid := smt.VerifySumProof(
		proof,
		msg.SmstRootHash,
		// INCOMPLETE: this **should not** be provided by the client (see above).
		msg.Path,
		msg.ValueHash,
		msg.SmstSum,
		smt.NoPrehashSpec(sha256.New(), true),
	); !valid {
		errInvalidProof := fmt.Errorf("invalid proof")
		logger.Error(errInvalidProof.Error())
		return nil, errInvalidProof
	}

	logger.Debug("proof verified")

	return &types.MsgProofResponse{}, nil
}
