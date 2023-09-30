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
	logger.Info(fmt.Sprintf("About to undelegate application %v from %+v", msg.Address, msg.PortalPubKey.Value))
	// Update the store
	if err := k.UndelegatePortal(ctx, msg.Address, *msg.PortalPubKey); err != nil {
		logger.Error(fmt.Sprintf("could not update store with delegated portal for application: %v", msg.Address))
		return nil, err
	}
	logger.Info(fmt.Sprintf("Successfully undelegated application %v from %+v", msg.Address, msg.PortalPubKey.Value))
	return &types.MsgUndelegateFromPortalResponse{}, nil
}
