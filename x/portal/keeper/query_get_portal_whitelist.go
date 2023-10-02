package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"poktroll/x/portal/types"
)

func (k Keeper) GetPortalWhitelist(goCtx context.Context, req *types.QueryGetPortalWhitelistRequest) (*types.QueryGetPortalWhitelistResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetWhitelist(
		ctx,
		req.PortalAddress,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetPortalWhitelistResponse{AppAddresses: val}, nil
}
