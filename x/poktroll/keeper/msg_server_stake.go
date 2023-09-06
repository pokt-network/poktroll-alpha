package keeper

import (
	"context"

	"poktroll/x/poktroll/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) Stake(goCtx context.Context, msg *types.MsgStake) (*types.MsgStakeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Define actor-specific logic here. Assuming just servicers for now
	err := k.StakeActor(ctx, sdk.ValAddress(msg.Creator), msg.amount)

	if err != nil {
		return &types.MsgStakeResponse{Success: false}, err
	}

	return &types.MsgStakeResponse{Success: true}, nil
}
