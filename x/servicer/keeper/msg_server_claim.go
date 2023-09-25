package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"poktroll/x/servicer/types"
)

func (k msgServer) Claim(goCtx context.Context, msg *types.MsgClaim) (*types.MsgClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	//logger := k.Logger(ctx).With("method", "Claim")

	// INCOMPLETE: validate that the message is signed by the servicer address
	// in the claim.
	//
	// CONSIDERATION: use `msg.GetSigners()` or the `cosmos.msg.v1.signer`
	// option on `MsgClaim` protobuf type.
	// (see: https://github.com/cosmos/cosmos-sdk/blob/main/proto/cosmos/bank/v1beta1/bank.proto#L34C1-L35C1)

	if err := k.InsertClaim(ctx, msg); err != nil {
		return nil, err
	}
	return &types.MsgClaimResponse{}, nil
}
