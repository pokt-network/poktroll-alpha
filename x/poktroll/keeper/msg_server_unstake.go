package keeper

import (
	"context"
	"fmt"

	"poktroll/x/poktroll/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// REFACTOR: Introduce actor specific staking commands

func (k msgServer) Unstake(goCtx context.Context, msg *types.MsgUnstake) (*types.MsgUnstakeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := ctx.Logger().With("message", "Unstake")
	if err := k.UnstakeActor(ctx, msg); err != nil {
		logger.Error("Error unstaking", err.Error())
		return &types.MsgUnstakeResponse{Success: false}, err
	}
	logger.Info(fmt.Sprintf("%s Unstaked %s %s", msg.Creator, msg.Amount, msg.ActorType))
	return &types.MsgUnstakeResponse{Success: true}, nil
}
