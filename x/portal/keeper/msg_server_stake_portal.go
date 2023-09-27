package keeper

import (
	"context"
	"fmt"

	"poktroll/x/portal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) StakePortal(goCtx context.Context, msg *types.MsgStakePortal) (*types.MsgStakePortalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := k.Logger(ctx).With("method", "StakePortal")
	logger.Info(fmt.Sprintf("About to stake portal %v with %v", msg.Address, msg.StakeAmount))

	// Get the address of the staking portal
	appAddress, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		logger.Error(fmt.Sprintf("could not parse address %v", msg.Address))
		return nil, err
	}

	// Determine the new stake amount
	newPortalStake, err := sdk.ParseCoinNormalized(msg.StakeAmount.String())
	if err != nil {
		logger.Error(fmt.Sprintf("could not parse stake amount %v", msg.StakeAmount))
		return nil, sdkerrors.ErrInvalidAddress
	}

	// Get the previously staked portal
	var coinsToSend sdk.Coin
	portal, found := k.GetPortal(ctx, msg.Address)
	if !found {
		logger.Info(fmt.Sprintf("portal not found, creating new portal for address %s with stake amount %v", appAddress, newPortalStake))

		// If the portal is not found, create a new one
		portal = types.Portal{
			Address:  msg.Address,
			Stake:    msg.StakeAmount,
			Services: msg.Services,
		}

		// Determine the number of coins to send from the portal address to the portal module account
		coinsToSend = newPortalStake
	} else {
		logger.Info(fmt.Sprintf("portal found, updating portal stake from %v to %v", portal.Stake, newPortalStake))

		// If the portal is found, make sure the stake amount has increased
		if portal.Stake.IsGTE(newPortalStake) {
			logger.Error(fmt.Sprintf("stake amount must %v be higher than previous stake amount %v", newPortalStake, portal.Stake))
			return nil, types.ErrStakeAmountMustBeHigher
		}

		// Determine the number of coins to send from the portal address to the portal module account
		coinsToSend = newPortalStake.Sub(*portal.Stake)
		portal.Stake = &newPortalStake

		// Update the services (just an override operation)
		portal.Services = msg.Services
	}

	// Send coins to the portal module account
	sdkError := k.bankKeeper.SendCoinsFromAccountToModule(ctx, appAddress, types.ModuleName, []sdk.Coin{coinsToSend})
	if sdkError != nil {
		logger.Error(fmt.Sprintf("could not send coins %v coins from %s to %s module account due to %v", coinsToSend, appAddress, types.ModuleName, sdkError))
		return nil, sdkError
	}
	logger.Info(fmt.Sprintf("successfully sent coins %v from %s to %s module account", coinsToSend, appAddress, types.ModuleName))

	// Update the portal in the store
	k.SetPortal(ctx, portal)
	logger.Info(fmt.Sprintf("successfully updated portal in the store: %v", portal))

	// QED
	return &types.MsgStakePortalResponse{}, nil
}
