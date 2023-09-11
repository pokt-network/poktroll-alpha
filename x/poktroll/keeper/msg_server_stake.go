package keeper

import (
	"context"
	"fmt"

	"poktroll/x/poktroll/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// REFACTOR: Introduce actor specific staking commands

func (k msgServer) Stake(goCtx context.Context, msg *types.MsgStake) (*types.MsgStakeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := ctx.Logger()
	if err := k.StakeActor(ctx, msg); err != nil {
		logger.Error("Error staking", err.Error())
		return &types.MsgStakeResponse{Success: false}, err
	}
	logger.Info(fmt.Sprintf("%s %s Staked %s", msg.ActorType, msg.Creator, msg.Amount))
	return &types.MsgStakeResponse{Success: true}, nil
}
