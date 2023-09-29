package keeper

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pokt-network/smt"

	"poktroll/x/servicer/types"
	sessionkeeper "poktroll/x/session/keeper"
)

const (
	// INCOMPLETE/HACK: this should be a governance param.
	// govSessionEndHeightOffset is a constant number of blocks after the end of
	// a session, after which a claim for that session can be submitted.
	govClaimCommittedHeightOffset = sessionkeeper.NumSessionBlocks / 2
)

var errInvalidPathFmt = "invalid path: %x, expected: %x"

// TODO_INCOMPLETE: Just some placeholder implementation for the proof on the server side for now.
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

	// lookup the corresponding claim and verify that it matches.
	claim, err := k.GetClaim(ctx, msg.GetSession().GetSessionId())
	if err != nil {
		return nil, err
	}

	// verify the claim is for the same session tree.
	if !bytes.Equal(claim.SmstRootHash, msg.SmstRootHash) {
		return nil, fmt.Errorf("smst root hash mismatch, expected: %x; got: %x", claim.SmstRootHash, msg.SmstRootHash)
	}

	// ASSUMPTION: the first signer is the servicer address.
	firstSignerAddress := msg.GetSigners()[0]
	if claim.ServicerAddress != firstSignerAddress.String() {
		// TODO_THIS_COMMIT: make a cosmos-sdk error for this.
		return nil, fmt.Errorf("first proof signer doesn't match claim's servicer address, expected: %s; got: %s", claim.ServicerAddress, firstSignerAddress)
	}

	// WIP WIP WIP // WIP WIP WIP // WIP WIP WIP // WIP WIP WIP
	// HELP: expecting `k.sessionKeeper != nil`
	session, err := k.sessionKeeper.GetSessionForApp(
		ctx, msg.GetSession().GetApplication().GetAddress(),
		msg.GetSession().GetService().GetId(),
		// TODO_THIS_COMMIT: replace with `Session#GetSessionBlockStartHeight()`
		msg.GetSession().GetSessionNumber(),
	)
	if err != nil {
		// TODO_THIS_COMMIT: make a cosmos-sdk error for this.
		return nil, fmt.Errorf("failed to get session for app: %w", err)
	}

	// TODO_CONSIDERATION: we can  do this in terms of sessionId instead of
	// claimCommittedBlockHash; however, it would require refactoring the
	//claimCommittedHeightCtx := ctx.WithBlockHeight(int64(claim.GetCommittedHeight()))
	//claimCommittedBlockHash := claimCommittedHeightCtx.BlockHeader().LastBlockId.Hash
	// TODO_THIS_COMMIT: seed should be the claim's sessionId.
	earliestProofHeight := getPseudoRandomHeightOffset(
		//claimCommittedBlockHash,
		//claim.GetSessionId(),
		session.GetSessionId(),
		claim.GetCommittedHeight(),
		govClaimCommittedHeightOffset,
	)

	// proof is too early
	// RATIONALE: distribute the load of proofs across the session block range.
	// IMPROVE/INVESTIGATE: if the randClaimCommittedHeightOffsets could be
	// generated in a normal (or alternative) distribution, we can focus the
	// commit heights of the majority of claims while still being random and
	// fair.
	if uint64(ctx.BlockHeight()) < earliestProofHeight {
		// TODO_THIS_COMMIT: uncomment - currently debugging depinject
		// & servicer/session module dep cycle
		//
		// TODO_THIS_COMMIT: make a cosmos-sdk error for this.
		//return nil, fmt.Errorf(
		//	"proof submitted too early, earliest proof height: %d; got: %d",
		//	earliestProofHeight,
		//	ctx.BlockHeight(),
		//)
	}

	lastEndedSessionNumber := uint64(ctx.BlockHeight()) / sessionkeeper.NumSessionBlocks
	currentSessionEndHeight := (lastEndedSessionNumber + 1) * sessionkeeper.NumSessionBlocks

	// proof is too late
	// RATIONALE: only rewarding proofs committed before some threshold
	// This allows us to set an upper bound on application unstake delay.
	if uint64(ctx.BlockHeight()) > currentSessionEndHeight {
		// TODO_THIS_COMMIT: make a cosmos-sdk error for this.
		return nil, fmt.Errorf(
			"proof submitted too late, current session end height: %d; got: %d",
			currentSessionEndHeight,
			ctx.BlockHeight(),
		)
	}

	// two parts to earliestProofHeight (offsets); one is constant (gov param) &
	// the other is pseudo-random.

	// INCOMPLETE: we need to verify that the proof height is greater than
	// earliestProofHeight and less than currentSessionEndHeight.
	//
	// latestProofHeight should be calculated from a governance parameter and
	// substituted for `currentSessionEndHeight` above.

	if valid := smt.VerifySumProof(
		proof,
		msg.SmstRootHash,
		// INCOMPLETE: this **should not** be provided by the client (see above).
		msg.Path,
		msg.ValueHash,
		msg.SmstSum,
		smt.NoPrehashSpec(sha256.New(), true),
	); !valid {
		// TODO_THIS_COMMIT: make a cosmos-sdk error for this.
		errInvalidProof := fmt.Errorf("failed to validate proof")

		// TECHDEBT: remove this error logs; they're intended for development use only.
		logger.Error(errInvalidProof.Error())

		return nil, errInvalidProof
	}

	if err := k.InsertProof(ctx, msg); err != nil {
		// TECHDEBT: remove this error logs; they're intended for development use only.
		logger.Error("failed to insert proof")

		return nil, err
	}

	return &types.MsgProofResponse{}, nil
}
