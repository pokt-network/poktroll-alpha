package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"poktroll/x/poktroll/types"
)

func (k msgServer) SubmitClaim(goCtx context.Context, msg *types.MsgSubmitClaim) (*types.MsgSubmitClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgSubmitClaimResponse{}, nil
}
