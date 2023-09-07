package keeper

import (
	"context"

	"poktroll/x/poktroll/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) Unstake(goCtx context.Context, msg *types.MsgUnstake) (*types.MsgUnstakeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Convert msg.Amount string to a Coin value
	amount, err := sdk.ParseDecCoin(msg.Amount)
	if err != nil {
		return &types.MsgUnstakeResponse{Success: false}, err
	}
	stake_amt, _ := amount.TruncateDecimal() // TODO: NEed tp deal with better parsing

	// TODO: Define actor-specific logic here. Assuming just servicers for now
	err = k.UnstakeActor(ctx, sdk.ValAddress(msg.Creator), stake_amt)

	if err != nil {
		return &types.MsgUnstakeResponse{Success: false}, err
	}

	return &types.MsgUnstakeResponse{Success: true}, nil
}
