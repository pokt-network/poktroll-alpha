package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"poktroll/x/servicer/types"
)

func (k msgServer) Claim(goCtx context.Context, msg *types.MsgClaim) (*types.MsgClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// INCOMPLETE: using an error for high contrast in log output; revert this.
	ctx.Logger().With("method", "Claim").Error("CLAIM SUBMITTED")

	// TODO: Handling the message
	_ = ctx

	return &types.MsgClaimResponse{}, nil
}
