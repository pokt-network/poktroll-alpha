package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"poktroll/x/servicer/types"
)

func (k Keeper) ServicersAll(goCtx context.Context, req *types.QueryAllServicersRequest) (*types.QueryAllServicersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var servicerss []types.Servicers
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	servicersStore := prefix.NewStore(store, types.KeyPrefix(types.ServicersKeyPrefix))

	pageRes, err := query.Paginate(servicersStore, req.Pagination, func(key []byte, value []byte) error {
		var servicers types.Servicers
		if err := k.cdc.Unmarshal(value, &servicers); err != nil {
			return err
		}

		servicerss = append(servicerss, servicers)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllServicersResponse{Servicers: servicerss, Pagination: pageRes}, nil
}

func (k Keeper) Servicers(goCtx context.Context, req *types.QueryGetServicersRequest) (*types.QueryGetServicersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetServicers(
		ctx,
		req.Index,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetServicersResponse{Servicers: val}, nil
}
