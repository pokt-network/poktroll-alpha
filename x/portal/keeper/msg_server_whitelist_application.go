package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"poktroll/x/portal/types"
)

func (k msgServer) WhitelistApplication(goCtx context.Context, msg *types.MsgWhitelistApplication) (*types.MsgWhitelistApplicationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := k.Logger(ctx).With("method", "WhitelistApplication")
	logger.Info(fmt.Sprintf("About to whitelist application %s to portal %s", msg.AppAddress, msg.PortalAddress))

	// DISCUSS: Do we need to check the application exists too?
	portal, found := k.GetPortal(ctx, msg.PortalAddress)
	if !found {
		logger.Error(fmt.Sprintf("portal not found for address %s", msg.PortalAddress))
		return nil, types.ErrUnstakingNonExistentPortal
	}

	// check the app isn't already whitelisted
	found = false
	for _, a := range portal.WhitelistedApps.AppAddresses {
		if a == msg.AppAddress {
			found = true
		}
	}
	if found {
		logger.Error(fmt.Sprintf("portal [%s] already whitelisted app: %s", msg.PortalAddress, msg.AppAddress))
		return nil, types.ErrAppAlreadyWhitelisted
	}

	// Update the application in the store
	if err := k.WhitelistApp(ctx, msg.PortalAddress, msg.AppAddress); err != nil {
		logger.Error(fmt.Errorf("unable to update portal state: %w", err).Error())
		return nil, err
	}
	logger.Info(fmt.Sprintf("successfully updated portal's [%s] whitelist to include: %s", msg.PortalAddress, msg.AppAddress))

	// QED
	return &types.MsgWhitelistApplicationResponse{}, nil
}
