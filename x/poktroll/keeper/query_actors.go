package keeper

import (
	"context"

	"poktroll/x/poktroll/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Actors(goCtx context.Context, req *types.QueryActorsRequest) (*types.QueryActorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := ctx.Logger()

	// Get the key-value module store using the store key (in our case store key is "chain")
	store := ctx.KVStore(k.storeKey)

	// TODO: Add other actor types here
	var servicers []*types.Servicer
	servicerStore := prefix.NewStore(store, []byte(types.ServicerPrefix))

	// Paginate the recipes store based on PageRequest
	pageRes, err := query.Paginate(servicerStore, req.Pagination, func(key []byte, value []byte) error {
		var servicer types.Servicer
		if err := k.cdc.Unmarshal(value, &servicer); err != nil {
			logger.Error("could not unmarshal servicer")
			return err
		}
		servicers = append(servicers, &servicer)
		return nil
	})

	// Throw an error if pagination failed
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Return a struct containing a list of recipes and pagination info
	return &types.QueryActorsResponse{Servicers: servicers, Pagination: pageRes}, nil
}
