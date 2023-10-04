package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"poktroll/x/portal/types"
)

func (k msgServer) AllowlistApplication(goCtx context.Context, msg *types.MsgAllowlistApplication) (*types.MsgAllowlistApplicationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := k.Logger(ctx).With("method", "AllowlistApplication")
	logger.Info(fmt.Sprintf("About to allowlist application %s to portal %s", msg.AppAddress, msg.PortalAddress))

	// DISCUSS: Do we need to check the application exists too?
	portal, found := k.GetPortal(ctx, msg.PortalAddress)
	if !found {
		logger.Error(fmt.Sprintf("portal not found for address %s", msg.PortalAddress))
		return nil, types.ErrUnstakingNonExistentPortal
	}

	// check the app isn't already allowlisted
	found = false
	for _, a := range portal.AllowlistedApps {
		if a == msg.AppAddress {
			found = true
		}
	}
	if found {
		logger.Error(fmt.Sprintf("portal [%s] already allowlisted app: %s", msg.PortalAddress, msg.AppAddress))
		return nil, types.ErrAppAlreadyAllowlisted
	}

	// Update the application in the store
	if err := k.AllowlistApp(ctx, msg.PortalAddress, msg.AppAddress); err != nil {
		logger.Error(fmt.Errorf("unable to update portal state: %w", err).Error())
		return nil, err
	}
	logger.Info(fmt.Sprintf("successfully updated portal's [%s] allowlist to include: %s", msg.PortalAddress, msg.AppAddress))

	// QED
	return &types.MsgAllowlistApplicationResponse{}, nil
}
