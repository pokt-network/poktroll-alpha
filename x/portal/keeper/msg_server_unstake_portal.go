package keeper

import (
	"context"
	"fmt"

	"poktroll/x/portal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// INCOMPLETE: This should start an unbonding period for the portal instead of unstaking it immediately
func (k msgServer) UnstakePortal(goCtx context.Context, msg *types.MsgUnstakePortal) (*types.MsgUnstakePortalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := k.Logger(ctx).With("method", "UnstakePortal")
	logger.Info(fmt.Sprintf("About to unstake portal %v", msg.Address))

	// Get the address of the staking portal
	portalAddress, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		logger.Error(fmt.Sprintf("could not parse address %v", msg.Address))
		return nil, err
	}

	portal, found := k.GetPortal(ctx, msg.Address)
	if !found {
		logger.Error(fmt.Sprintf("portal not found for address %s", msg.Address))
		return nil, types.ErrUnstakingNonExistentPortal
	}

	sdkError := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, portalAddress, []sdk.Coin{*portal.Stake})
	if sdkError != nil {
		logger.Error(fmt.Sprintf("could not send coins %v coins from %s module account to %s due to %v", portal.Stake, types.ModuleName, portalAddress, sdkError))
		return nil, sdkError
	}

	logger.Info(fmt.Sprintf("successfully sent coins %v from %s module account to %s", portal.Stake, types.ModuleName, portalAddress))

	// Update the portal in the store
	k.RemovePortal(ctx, msg.Address)
	logger.Info(fmt.Sprintf("successfully unstaked portal and removed it from the store: %v", portal))

	// QED
	return &types.MsgUnstakePortalResponse{}, nil
}
