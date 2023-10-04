package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"poktroll/x/portal/types"
)

func (k Keeper) GetPortalAllowlist(goCtx context.Context, req *types.QueryGetPortalAllowlistRequest) (*types.QueryGetPortalAllowlistResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetAllowlist(
		ctx,
		req.PortalAddress,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetPortalAllowlistResponse{AppAddresses: val}, nil
}
