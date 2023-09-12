package keeper

import (
	"context"
	"fmt"

	"poktroll/x/application/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) StakeApplication(goCtx context.Context, msg *types.MsgStakeApplication) (*types.MsgStakeApplicationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := k.Logger(ctx).With("method", "StakeApplication")
	logger.Info(fmt.Sprintf("About to stake application %v with %v", msg.Address, msg.StakeAmount))

	// Get the address of the staking application
	appAddress, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		logger.Error(fmt.Sprintf("could not parse address %v", msg.Address))
		return nil, err
	}

	// Determine the new stake amount
	newApplicationStake, err := sdk.ParseCoinNormalized(msg.StakeAmount.String())
	if err != nil {
		logger.Error(fmt.Sprintf("could not parse stake amount %v", msg.StakeAmount))
		return nil, sdkerrors.ErrInvalidAddress
	}

	// Get the previously staked application
	var coinsToSend sdk.Coin
	application, found := k.GetApplication(ctx, msg.Address)
	if !found {
		// If the application is not found, create a new one
		application = types.Application{
			Address: msg.Address,
			Stake:   msg.StakeAmount,
		}
		coinsToSend = newApplicationStake
		logger.Info(fmt.Sprintf("application not found, creating new application for address %s with stake amount %v", appAddress, newApplicationStake))
	} else {
		// If the application is found, make sure the stake amount has increased
		if application.Stake.IsGTE(newApplicationStake) {
			logger.Error(fmt.Sprintf("stake amount must %v be higher than previous stake amount %v", newApplicationStake, application.Stake))
			return nil, types.ErrStakeAmountMustBeHigher
		}
		logger.Info(fmt.Sprintf("application found, updating application stake from %v to %v", application.Stake, newApplicationStake))
		coinsToSend = newApplicationStake.Sub(*application.Stake)
		application.Stake = &newApplicationStake
	}

	// Send coins to the application module account
	sdkError := k.bankKeeper.SendCoinsFromAccountToModule(ctx, appAddress, types.ModuleName, []sdk.Coin{coinsToSend})
	if sdkError != nil {
		logger.Error(fmt.Sprintf("could not send coins %v coins from %s to %s module account due to %v", coinsToSend, appAddress, types.ModuleName, sdkError))
		return nil, sdkError
	}
	logger.Info(fmt.Sprintf("successfully sent coins %v from %s to %s module account", coinsToSend, appAddress, types.ModuleName))

	// Update the application in the store
	k.SetApplication(ctx, application)
	logger.Info(fmt.Sprintf("successfully updated application in the store: %v", application))

	// QED
	return &types.MsgStakeApplicationResponse{}, nil
}
