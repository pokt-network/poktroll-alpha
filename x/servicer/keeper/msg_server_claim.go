package keeper

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"

	servicertypes "poktroll/x/servicer/types"
	sessionkeeper "poktroll/x/session/keeper"
)

const (
	// INCOMPLETE/HACK: this should be a governance param.
	// govSessionEndHeightOffset is a constant number of blocks after the end of
	// a session, after which a claim for that session can be submitted.
	govSessionEndHeightOffset = sessionkeeper.NumSessionBlocks / 2
)

func (k msgServer) Claim(goCtx context.Context, msg *servicertypes.MsgClaim) (*servicertypes.MsgClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	//logger := k.Logger(ctx).With("method", "Claim")

	// CONSIDERATION: look into using the `cosmos.msg.v1.signer` option on
	// `MsgClaim` protobuf type instead of relying on `msg.GetSigners()`.
	// (see: https://github.com/cosmos/cosmos-sdk/blob/main/proto/cosmos/bank/v1beta1/bank.proto#L34C1-L35C1)

	// TECHDEBT: get `sessionkeeper.NumSessionBlocks` from governance parameters & calculate
	// session number from last height at which `sessionkeeper.NumSessionBlocks` changed (depends
	// on knowing the session number at the time of change).

	// impossible to submit a valid msg until after the first session has ended
	lastEndedSessionNumber := uint64(ctx.BlockHeight()) / sessionkeeper.NumSessionBlocks
	if lastEndedSessionNumber == 0 {
		// TODO_THIS_COMMIT: make a cosmos-sdk error for this
		return nil, fmt.Errorf("first session has not ended yet")
	}

	// TECHDEBT: If we're considering invalidation height... unclear why this would be useful.
	// (NB: this is carryover from V0)
	//if ctx.BlockHeight() > msg.GetInvalidationHeight() {
	//	// TODO_THIS_COMMIT: make a cosmos-sdk error for this
	//	return nil, fmt.Errorf(
	//		"msg session number is in the future, current session number: %d; got: %d",
	//		lastEndedSessionNumber,
	//		msg.Session.SessionNumber,
	//	)
	//}

	// claim must be for a session that has ended
	if lastEndedSessionNumber < msg.GetSessionNumber() {
		// TODO_THIS_COMMIT: make a cosmos-sdk error for this.
		return nil, fmt.Errorf("claim session number has not ended: %d", msg.GetSessionNumber())
	}

	// block#:                     [ 1 2 3 4 5 ][ 6 7 8 9 10 ]
	// session#:                   [ ↑   1     ][	↑   2      ]
	//  lastEndedSessionStartHeight ─┘   ↑──────┬───┤
	//           lastEndedSessionNumber ─┘	    │   claimCommitedHeight
	//

	// TODO_CONSIDERATION: we can  do this in terms of sessionId instead of
	// lastEndedSessionBlockHash; however, it would require refactoring the
	// servicer and/or session modules to eliminate a dependency cycle between
	// their protobuf message types.
	lastEndedSessionEndHeight := lastEndedSessionNumber * sessionkeeper.NumSessionBlocks
	//lastEndedSessionStartCtx := ctx.WithBlockHeight(int64(lastEndedSessionStartHeight))
	//lastEndedSessionBlockHash := lastEndedSessionStartCtx.BlockHeader().LastBlockId.Hash
	// TODO_THIS_COMMIT: seed should be the claim's sessionId.
	earliestClaimHeight := getPseudoRandomHeightOffset(
		//lastEndedSessionBlockHash,
		// *NOTE*: the session ID has not yet been verified! This will happen when
		// the corresponding proof is submitted along with the input arguments to
		// get the session (i.e. app address, service ID, block height / session
		// number)
		msg.GetSessionId(),
		lastEndedSessionEndHeight,
		govSessionEndHeightOffset,
	)

	// claim is too early
	// RATIONALE: distribute the load of proofs across the session block range.
	// IMPROVE/INVESTIGATE: if the randClaimCommittedHeightOffsets could be
	// generated in a normal (or alternative) distribution, we can focus the
	// commit heights of the majority of claims while still being random and
	// fair.
	if uint64(ctx.BlockHeight()) < earliestClaimHeight {
		// TODO_THIS_COMMIT: uncomment - currently debugging depinject
		// & servicer/session module dep cycle
		//
		// TODO_THIS_COMMIT: make a cosmos-sdk error for this.
		//return nil, fmt.Errorf(
		//	"claim submitted too early, earliest claim height: %d; got: %d",
		//	earliestClaimHeight,
		//	ctx.BlockHeight(),
		//)
	}

	claim := &servicertypes.Claim{
		// TECHDEBT: Insert the SessionHeader in the claim instead
		SessionId:       msg.GetSessionId(),
		SessionNumber:   lastEndedSessionNumber + 1,
		CommittedHeight: uint64(ctx.BlockHeight()),
		SmstRootHash:    msg.GetSmstRootHash(),
		ServicerAddress: msg.GetServicerAddress(),
	}

	if err := k.InsertClaim(ctx, claim); err != nil {
		return nil, err
	}

	if err := ctx.EventManager().EmitTypedEvent(msg.NewClaimEvent()); err != nil {
		return nil, err
	}

	return &servicertypes.MsgClaimResponse{}, nil
}

// TECHDEBT: factor this out to a library pkg so that it can be reused by the clien / relayer.
func getPseudoRandomHeightOffset(sessionId string, startHeight uint64, offset uint64) uint64 {
	rngSeed, _ := binary.Varint([]byte(sessionId))
	maxRandomSessionEndHeightOffset := sessionkeeper.NumSessionBlocks - govSessionEndHeightOffset
	// TECHDEBT: ensure use of a "universal" PRNG implementation; i.e. one that
	// is based on a spec and has multiple language implementations and/or bindings.
	// TODO_CONSIDERATION: it would be nice if the random offset component had
	// a normal distribution with respect to the session block range.
	// TODO_THIS_COMMIT: should take govClaimHeightOffset into account to avoid
	// proof submission in wrong (next) session.
	// INVESTIGATE: using "invariants" in cosmos-sdk to ensure that we don't
	// misconfigure  the chain params for this.
	randSessionEndHeightOffset := uint64(rand.NewSource(rngSeed).Int63()) % maxRandomSessionEndHeightOffset
	return startHeight + offset + randSessionEndHeightOffset
}
