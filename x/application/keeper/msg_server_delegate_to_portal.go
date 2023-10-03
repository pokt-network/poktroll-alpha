package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"poktroll/x/application/types"
)

func (k msgServer) DelegateToPortal(goCtx context.Context, msg *types.MsgDelegateToPortal) (*types.MsgDelegateToPortalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := k.Logger(ctx).With("method", "DelegateToPortal")
	logger.Info(fmt.Sprintf("About to delegate application %s to %s", msg.Address, msg.PortalAddress))
	// Update the store
	if err := k.DelegatePortal(ctx, msg.Address, msg.PortalAddress); err != nil {
		logger.Error(fmt.Sprintf("could not update store with delegated portal for application: %s", msg.Address))
		return nil, err
	}
	logger.Info(fmt.Sprintf("Successfully delegated application %s to %s", msg.Address, msg.PortalAddress))
	return &types.MsgDelegateToPortalResponse{}, nil
}
