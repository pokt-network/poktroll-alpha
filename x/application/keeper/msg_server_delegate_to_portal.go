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
	logger.Info(fmt.Sprintf("About to delegate application %v to %+v", msg.Address, msg.PortalPubKey.Value))

	// Update the store
	if err := k.DelegatePortal(ctx, msg.Address, *msg.PortalPubKey); err != nil {
		logger.Error(fmt.Sprintf("could not update store with delegated portal for application: %v", msg.Address))
		return nil, err
	}
	logger.Info(fmt.Sprintf("Successfully delegated application %v to %+v", msg.Address, msg.PortalPubKey.Value))

	return &types.MsgDelegateToPortalResponse{}, nil
}
