package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"poktroll/x/portal/types"
)

func (k msgServer) UnallowlistApplication(goCtx context.Context, msg *types.MsgUnallowlistApplication) (*types.MsgUnallowlistApplicationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := k.Logger(ctx).With("method", "UnallowlistApplication")
	logger.Info(fmt.Sprintf("About to unallowlist application %s from portal %s", msg.AppAddress, msg.PortalAddress))

	// DISCUSS: Do we need to check the application exists too?
	portal, found := k.GetPortal(ctx, msg.PortalAddress)
	if !found {
		logger.Error(fmt.Sprintf("portal not found for address %s", msg.PortalAddress))
		return nil, types.ErrUnstakingNonExistentPortal
	}

	// check the app is already allowlisted
	found = false
	for _, a := range portal.AllowlistedApps {
		if a == msg.AppAddress {
			found = true
		}
	}
	if !found {
		logger.Error(fmt.Sprintf("portal [%s] hasn't allowlisted app: %s", msg.PortalAddress, msg.AppAddress))
		return nil, types.ErrAppNotAllowlisted
	}

	// Update the application in the store
	if err := k.UnallowlistApp(ctx, msg.PortalAddress, msg.AppAddress); err != nil {
		logger.Error(fmt.Errorf("unable to update portal state: %w", err).Error())
		return nil, err
	}
	logger.Info(fmt.Sprintf("successfully updated portal's [%s] allowlist to exclude: %s", msg.PortalAddress, msg.AppAddress))

	// QED
	return &types.MsgUnallowlistApplicationResponse{}, nil
}
