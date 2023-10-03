package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"poktroll/x/application/types"
)

func (k msgServer) UndelegateFromPortal(goCtx context.Context, msg *types.MsgUndelegateFromPortal) (*types.MsgUndelegateFromPortalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := k.Logger(ctx).With("method", "UndelegateFromPortal")
	logger.Info(fmt.Sprintf("About to undelegate application %s from %s", msg.Address, msg.PortalAddress))
	// Update the store
	if err := k.UndelegatePortal(ctx, msg.Address, msg.PortalAddress); err != nil {
		logger.Error(fmt.Sprintf("could not update store with delegated portal for application: %s", msg.Address))
		return nil, err
	}
	logger.Info(fmt.Sprintf("Successfully undelegated application %s from %s", msg.Address, msg.PortalAddress))
	return &types.MsgUndelegateFromPortalResponse{}, nil
}
