package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"poktroll/x/servicer/types"
)

func (k msgServer) Claim(goCtx context.Context, claim *types.MsgClaim) (*types.MsgClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	//logger := k.Logger(ctx).With("method", "Claim")

	// INCOMPLETE: validate that the message is signed by the servicer address
	// in the claim.
	//
	// CONSIDERATION: use `claim.GetSigners()` or the `cosmos.claim.v1.signer`
	// option on `MsgClaim` protobuf type.
	// (see: https://github.com/cosmos/cosmos-sdk/blob/main/proto/cosmos/bank/v1beta1/bank.proto#L34C1-L35C1)

	// TODO_THIS_COMMIT: verify that the session in question is closed.

	if err := k.InsertClaim(ctx, claim); err != nil {
		return nil, err
	}

	if err := ctx.EventManager().EmitTypedEvent(claim.NewClaimEvent()); err != nil {
		return nil, err
	}

	return &types.MsgClaimResponse{}, nil
}
