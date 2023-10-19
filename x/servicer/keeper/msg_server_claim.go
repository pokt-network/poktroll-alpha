package keeper

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/rpc/client/http"

	servicertypes "poktroll/x/servicer/types"
	sessionkeeper "poktroll/x/session/keeper"
)

const (
	// INCOMPLETE/HACK: these constants should all be governance params. They are exported for the client to use.

	// GovEarliestClaimSubmissionBlocksOffset is a constant number of blocks after
	// the session end, before which any claim for that session could not be submitted.
	GovEarliestClaimSubmissionBlocksOffset int64 = 3

	// GovLatestClaimSubmissionBlocksInterval is a constant number of blocks after the
	// GovEarliestClaimSubmissionBlocksOffset, after which any claim for that session could no longer be submitted.
	GovLatestClaimSubmissionBlocksInterval int64 = 10

	// GovClaimSubmissionBlocksWindow is the number of blocks between which a claim
	// can be submitted. This is used to not impose the Relayer to submit the claim
	// at the exact block height.
	GovClaimSubmissionBlocksWindow int64 = 2
)

func (k msgServer) Claim(goCtx context.Context, msg *servicertypes.MsgClaim) (*servicertypes.MsgClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	client, err := http.New("http://localhost:36657")
	if err != nil {
		return nil, err
	}
	logger := k.Logger(ctx).With("method", "Claim")
	// CONSIDERATION: look into using the `cosmos.msg.v1.signer` option on
	// `MsgClaim` protobuf type instead of relying on `msg.GetSigners()`.
	// (see: https://github.com/cosmos/cosmos-sdk/blob/main/proto/cosmos/bank/v1beta1/bank.proto#L34C1-L35C1)

	// TECHDEBT: get `sessionkeeper.NumSessionBlocks` from governance parameters & calculate
	// session number from last height at which `sessionkeeper.NumSessionBlocks` changed (depends
	// on knowing the session number at the time of change).

	// Assert that there is only one signer (until we have rev share and delegation)
	signers := msg.GetSigners()
	if len(signers) != 1 {
		return nil, servicertypes.ErrUnsupportedMultiSig.Wrapf("got: %d", len(signers))
	}

	signer := msg.GetSigners()[0]
	if msg.GetServicerAddress() != signer.String() {
		return nil, servicertypes.ErrProofAndClaimSignerMismatch.Wrapf(
			"expected: %s;got: %s",
			msg.GetServicerAddress(),
			signer,
		)
	}

	numSessionBlocks := int64(sessionkeeper.NumSessionBlocks)
	currentBlockHeight := ctx.BlockHeight()
	lastEndedSessionNumber := currentBlockHeight / numSessionBlocks

	// impossible to submit a valid msg until the first session has ended
	if lastEndedSessionNumber == 0 {
		return nil, servicertypes.ErrActiveFirstSession
	}

	// block#:                   [ 1 ... 10 ][ 11 12 13 14 15 16 17 18 19 20 ][ 21 22 23 24 ... ]
	// session#:                 [ ↑   1  ↑ ][        ↑     ↑ 2   ↑          ][				 ↑ 3      ]
	//                             │   ↑  │           │     │     │                    │
	// lastEndedSessionNumber ─────│───┘  │           │     │     │                    │
	// (claimed session)           │      │           │     │     │                    │
	//                             │      │           │     └─────┴─────────────GovClaimSubmissionBlocksWindow
	// earliestClaimSubmissionBlockHeight─────────────┘     ↑                          │
	// = sessionStartHeight (1)────┘      │           ↑     │                          │
	// + numSessionBlocks (+10)───────────┘           │     │                          │
	// + GovEarliestClaimSubmissionBlocksOffset (3)   │     │                          │
	//                                                │     │                          │
	// GovLatestClaimSubmissionBlocksInterval (10)────┴─────│──────────────────────────┘
	//                                                      │
	// randClaimSubmissionBlockHeightOffset (seed = block[13].hash)
	//

	// TODO_THIS_COMMIT: factor all this out to a library pkg so that it can be
	// reused in the client / relayer.
	// TODO_CONSIDERATION: we can  do this in terms of sessionId instead of
	// sessionStartHeight; however, it would require refactoring the
	// servicer and/or session modules to eliminate a dependency cycle between
	// their protobuf message types.
	sessionStartHeight := int64(msg.SessionNumber) * numSessionBlocks

	// lastEndedSessionStartHeight + sessionkeeper.NumSessionBlocks will land us at the
	// first block of the the current session.
	// earliestClaimSubmissionBlockHeight is the latest block height that could be inferred
	// from gov params given the current block height.
	// we use its hash to deterministically generate a random offset for the claim committed height.
	earliestClaimSubmissionBlockHeight := sessionStartHeight + numSessionBlocks + GovEarliestClaimSubmissionBlocksOffset
	block, err := client.Block(goCtx, &earliestClaimSubmissionBlockHeight)
	if err != nil {
		return nil, err
	}
	earliestClaimSubmissionBlockHash := block.Block.Header.LastBlockID.Hash.Bytes()
	logger.Error("using block and hash", "earliestClaimSubmissionBlockHeight", earliestClaimSubmissionBlockHeight, "earliestClaimSubmissionBlockHash", fmt.Sprintf("%x", earliestClaimSubmissionBlockHash))
	rngSeed, _ := binary.Varint(earliestClaimSubmissionBlockHash)

	// TECHDEBT: ensure use of a "universal" PRNG implementation; i.e. one that
	// is based on a spec and has multiple language implementations and/or bindings.
	// TODO_CONSIDERATION: it would be nice if the random offset component had
	// a normal distribution with respect to the session block range.
	// INVESTIGATE: using "invariants" in cosmos-sdk to ensure that we don't
	// misconfigure  the chain params for this.
	randomNumber := rand.NewSource(rngSeed).Int63()
	randClaimSubmissionBlockHeightOffset := randomNumber % (GovLatestClaimSubmissionBlocksInterval - GovClaimSubmissionBlocksWindow - 1)

	// claim is too early
	// RATIONALE: distribute the load of proofs across the session block range.
	// IMPROVE/INVESTIGATE: if the randClaimSubmissionBlockHeightOffset could be
	// generated in a normal (or alternative) distribution, we can focus the
	// commit heights of the majority of claims while still being random and
	// fair.
	earliestServicerClaimSubmissionBlockHeight := earliestClaimSubmissionBlockHeight + randClaimSubmissionBlockHeightOffset + 1
	if currentBlockHeight < earliestServicerClaimSubmissionBlockHeight {
		return nil, servicertypes.ErrEarlyClaimSubmission.Wrapf(
			"early claim height: %d; got: %d",
			earliestServicerClaimSubmissionBlockHeight,
			currentBlockHeight,
		)
	}

	// claim is too late
	latestServicerClaimSubmissionBlockHeight := earliestServicerClaimSubmissionBlockHeight + GovClaimSubmissionBlocksWindow
	if currentBlockHeight > latestServicerClaimSubmissionBlockHeight {
		return nil, servicertypes.ErrLateClaimSubmission.Wrapf(
			"late claim height: %d; got: %d",
			latestServicerClaimSubmissionBlockHeight,
			currentBlockHeight,
		)
	}

	claim := &servicertypes.Claim{
		// TODO_CONSIDERATION: may not need `SessionId` field, session ID is the
		// key in the servicer/claims store.
		SessionId:       msg.GetSessionId(),
		SessionNumber:   uint64(lastEndedSessionNumber),
		CommittedHeight: uint64(currentBlockHeight),
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
