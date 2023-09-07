package keeper

import (
	"context"

	"poktroll/x/poktroll/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) Unstake(goCtx context.Context, msg *types.MsgUnstake) (*types.MsgUnstakeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Define actor-specific logic here. Assuming just servicers for now
	err := k.UnstakeActor(ctx, msg)

	if err != nil {
		return &types.MsgUnstakeResponse{Success: false}, err
	}

	return &types.MsgUnstakeResponse{Success: true}, nil
}
