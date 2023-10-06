package keeper

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pokt-network/smt"

	"poktroll/x/servicer/types"
)

const (
	// INCOMPLETE/HACK: these constants should all be governance params. They are exported for the client to use.

	// GovEarliestProofSubmissionBlocksOffset is a constant number of blocks after
	// claim submission, before which a proof for that claim could not be submitted.
	// TODO_IN_THIS_COMMIT: this should be a governance param.
	GovEarliestProofSubmissionBlocksOffset int64 = 0

	// GovLatestProofSubmissionBlocksInterval is a constant number of blocks after the
	// GovEarliestProofSubmissionBlocksOffset, after which a proof for that claim could no longer be submitted.
	GovLatestProofSubmissionBlocksInterval int64 = 30

	// GovProofSubmissionBlocksWindow is the number of blocks between which a proof
	// can be submitted. This is used to not impose the Relayer to submit the proof
	// at the exact block height.
	GovProofSubmissionBlocksWindow int64 = 2
)

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

	// lookup the corresponding claim and verify that it matches.
	claim, err := k.GetClaim(ctx, msg.SessionId)
	if err != nil {
		return nil, err
	}

	// verify the claim is for the same session tree.
	if !bytes.Equal(claim.SmstRootHash, msg.SmstRootHash) {
		return nil, types.ErrSMSTRootHashMismatch.Wrapf(
			"expected: %x; got: %x",
			claim.SmstRootHash,
			msg.SmstRootHash,
		)
	}

	// ASSUMPTION: the first signer is the servicer address.
	// Assert that there is only one signer (until we have rev share and delegation)
	signers := msg.GetSigners()
	if len(signers) != 1 {
		return nil, types.ErrUnsupportedMultiSig.Wrapf("got: %d", len(signers))
	}

	signer := msg.GetSigners()[0]
	if claim.ServicerAddress != signer.String() {
		return nil, types.ErrProofAndClaimSignerMismatch.Wrapf(
			"expected: %s;got: %s",
			claim.ServicerAddress,
			signer,
		)
	}

	currentBlockHeight := uint64(ctx.BlockHeight())

	// earliestProofSubmissionBlockHeight is the earliest block height at which any Servicer has
	// to submit a proof for a claim.
	// it is the latest block height that could be inferred from past commitments (claims) and governance params.
	// we use its hash to seed a PRNG to generate a random offset.
	earliestProofSubmissionBlockHeight := int64(claim.GetCommittedHeight()) + GovEarliestProofSubmissionBlocksOffset

	// TODO_THIS_COMMIT: factor all this out to a library pkg so that it can be
	// reused in the client / relayer.
	earliestProofSubmissionBlockCtx := ctx.WithBlockHeight(int64(earliestProofSubmissionBlockHeight))
	earliestProofSubmissionBlockHash := earliestProofSubmissionBlockCtx.BlockHeader().LastBlockId.Hash
	rngSeed, _ := binary.Varint(earliestProofSubmissionBlockHash)

	// TECHDEBT: ensure use of a "universal" PRNG implementation; i.e. one that
	// is based on a spec and has multiple language implementations and/or bindings.
	// TODO_CONSIDERATION: it would be nice if the random offset component had
	// a normal distribution with respect to the session block range.
	// INVESTIGATE: using "invariants" in cosmos-sdk to ensure that we don't
	// misconfigure  the chain params for this.
	randomNumber := rand.NewSource(rngSeed).Int63()
	randProofSubmissionBlockHeightOffset := randomNumber % (GovLatestProofSubmissionBlocksInterval - GovProofSubmissionBlocksWindow)

	// proof is too early
	// RATIONALE: distribute the load of proofs across the session block range.
	// IMPROVE/INVESTIGATE: if the randClaimCommittedHeightOffsets could be
	// generated in a normal (or alternative) distribution, we can focus the
	// commit heights of the majority of claims while still being random and
	// fair.
	earliestServicerProofSubmissionBlockHeight := earliestProofSubmissionBlockHeight + randProofSubmissionBlockHeightOffset
	if currentBlockHeight < uint64(earliestServicerProofSubmissionBlockHeight) {
		return nil, types.ErrEarlyProofSubmission.Wrapf(
			"earliest proof height: %d; got: %d",
			earliestServicerProofSubmissionBlockHeight,
			currentBlockHeight,
		)
	}

	// proof is too late
	// RATIONALE: only rewarding proofs committed before some threshold
	// This allows us to set an upper bound on application unstake delay.
	latestServicerClaimSubmissionBlockHeight := earliestProofSubmissionBlockHeight + GovProofSubmissionBlocksWindow
	if currentBlockHeight > uint64(latestServicerClaimSubmissionBlockHeight) {
		return nil, types.ErrLateProofSubmission.Wrapf(
			"current session end height: %d; got: %d",
			latestServicerClaimSubmissionBlockHeight,
			currentBlockHeight,
		)
	}

	// use smt.VerifyClosestSumProof(
	if valid := smt.VerifySumProof(
		proof,
		msg.SmstRootHash,
		// INCOMPLETE: this **should not** be provided by the client (see above).
		msg.Path,
		msg.ValueHash,
		msg.SmstSum,
		smt.NoPrehashSpec(sha256.New(), true),
	); !valid {

		// TECHDEBT: remove this error logs; they're intended for development use only.
		logger.Error(types.ErrInvalidProof.Error())

		return nil, types.ErrInvalidProof
	}

	if err := k.InsertProof(ctx, msg); err != nil {
		// TECHDEBT: remove this error logs; they're intended for development use only.
		logger.Error("failed to insert proof")

		return nil, err
	}

	return &types.MsgProofResponse{}, nil
}
