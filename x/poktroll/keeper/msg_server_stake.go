package keeper

import (
	"context"

	"poktroll/x/poktroll/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) Stake(goCtx context.Context, msg *types.MsgStake) (*types.MsgStakeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Define actor-specific logic here. Assuming just servicers for now
	err := k.StakeActor(ctx, msg)

	logger := ctx.Logger()

	if err != nil {
		logger.Info("staking unsuccesful :(")
		logger.Info(err.Error())
		return &types.MsgStakeResponse{Success: false}, err
	}

	logger.Info("staked!")
	return &types.MsgStakeResponse{Success: true}, nil
}
