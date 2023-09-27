package keeper

import (
	"context"
	"fmt"

	"poktroll/x/portal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// INCOMPLETE: This should start an unbonding period for the application instead of unstaking it immediately
func (k msgServer) UnstakePortal(goCtx context.Context, msg *types.MsgUnstakePortal) (*types.MsgUnstakePortalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := k.Logger(ctx).With("method", "UnstakePortal")
	logger.Info(fmt.Sprintf("About to unstake application %v", msg.Address))

	// Get the address of the staking application
	appAddress, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		logger.Error(fmt.Sprintf("could not parse address %v", msg.Address))
		return nil, err
	}

	application, found := k.GetPortal(ctx, msg.Address)
	if !found {
		logger.Error(fmt.Sprintf("application not found for address %s", msg.Address))
		return nil, types.ErrUnstakingNonExistentPortal
	}

	sdkError := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, appAddress, []sdk.Coin{*application.Stake})
	if sdkError != nil {
		logger.Error(fmt.Sprintf("could not send coins %v coins from %s module account to %s due to %v", application.Stake, types.ModuleName, appAddress, sdkError))
		return nil, sdkError
	}

	logger.Info(fmt.Sprintf("successfully sent coins %v from %s module account to %s", application.Stake, types.ModuleName, appAddress))

	// Update the application in the store
	k.RemovePortal(ctx, msg.Address)
	logger.Info(fmt.Sprintf("successfully unstaked application and removed it from the store: %v", application))

	// QED
	return &types.MsgUnstakePortalResponse{}, nil
}
