package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"poktroll/x/portal/types"
)

func (k Keeper) GetDelegatedPortals(goCtx context.Context, req *types.QueryGetDelegatedPortalsRequest) (*types.QueryGetDelegatedPortalsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	delegatees, found := k.GetDelegatees(ctx, req.AppAddress)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetDelegatedPortalsResponse{Delegatees: delegatees}, nil
}
