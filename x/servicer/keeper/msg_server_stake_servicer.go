package keeper

import (
	"context"
	"fmt"

	"poktroll/x/servicer/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) StakeServicer(goCtx context.Context, msg *types.MsgStakeServicer) (*types.MsgStakeServicerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := k.Logger(ctx).With("method", "StakeServicer")
	logger.Info(fmt.Sprintf("About to stake servicer %v with %v", msg.Address, msg.StakeAmount))

	// Get the address of the staking servicer
	appAddress, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		logger.Error(fmt.Sprintf("could not parse address %v", msg.Address))
		return nil, err
	}

	// Determine the new stake amount
	newServicerStake, err := sdk.ParseCoinNormalized(msg.StakeAmount.String())
	if err != nil {
		logger.Error(fmt.Sprintf("could not parse stake amount %v", msg.StakeAmount))
		return nil, sdkerrors.ErrInvalidAddress
	}

	// Get the previously staked servicer
	var coinsToSend sdk.Coin
	servicer, found := k.GetServicers(ctx, msg.Address)
	if !found {
		logger.Info(fmt.Sprintf("servicer not found, creating new servicer for address %s with stake amount %v", appAddress, newServicerStake))

		// If the servicer is not found, create a new one
		servicer = types.Servicers{
			Address:  msg.Address,
			Stake:    msg.StakeAmount,
			Services: msg.Services,
		}

		// Determine the number of coins to send from the servicer address to the servicer module account
		coinsToSend = newServicerStake
	} else {
		logger.Info(fmt.Sprintf("servicer found, updating servicer stake from %v to %v", servicer.Stake, newServicerStake))

		// If the servicer is found, make sure the stake amount has increased
		if servicer.Stake.IsGTE(newServicerStake) {
			logger.Error(fmt.Sprintf("stake amount must %v be higher than previous stake amount %v", newServicerStake, servicer.Stake))
			return nil, types.ErrStakeAmountMustBeHigher
		}

		// Determine the number of coins to send from the servicer address to the servicer module account
		coinsToSend = newServicerStake.Sub(*servicer.Stake)
		servicer.Stake = &newServicerStake

		// Update the services (just an override operation)
		for _, serviceConfig := range msg.Services {
			if err := serviceConfig.ValidateEndpoints(); err != nil {
				// TODO_THIS_COMMIT: improve error
				return nil, err
			}
		}
		servicer.Services = msg.Services
	}

	// Send coins to the servicer module account
	sdkError := k.bankKeeper.SendCoinsFromAccountToModule(ctx, appAddress, types.ModuleName, []sdk.Coin{coinsToSend})
	if sdkError != nil {
		logger.Error(fmt.Sprintf("could not send coins %v coins from %s to %s module account due to %v", coinsToSend, appAddress, types.ModuleName, sdkError))
		return nil, sdkError
	}
	logger.Info(fmt.Sprintf("successfully sent coins %v from %s to %s module account", coinsToSend, appAddress, types.ModuleName))

	// Update the servicer in the store
	k.SetServicers(ctx, servicer)
	logger.Info(fmt.Sprintf("successfully updated servicer in the store: %v", servicer))

	// QED
	return &types.MsgStakeServicerResponse{}, nil
}
