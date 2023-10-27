package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/types/query"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"poktroll/x/portal/types"
)

func (k Keeper) Portals(goCtx context.Context, req *types.QueryAllPortalsRequest) (*types.QueryAllPortalsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var portals []types.Portal
	ctx := sdk.UnwrapSDKContext(goCtx)
	store := ctx.KVStore(k.storeKey)
	portalStore := prefix.NewStore(store, types.KeyPrefix(types.PortalKeyPrefix))

	pageRes, err := query.Paginate(portalStore, req.Pagination, func(key []byte, value []byte) error {
		var portal types.Portal
		if err := k.cdc.Unmarshal(value, &portal); err != nil {
			return err
		}
		portals = append(portals, portal)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllPortalsResponse{Portals: portals, Pagination: pageRes}, nil
}

func (k Keeper) Portal(goCtx context.Context, req *types.QueryGetPortalRequest) (*types.QueryGetPortalResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetPortal(
		ctx,
		req.Address,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetPortalResponse{Portal: val}, nil
}
