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

	// Define a variable that will store a list of servicers
	var actors []*types.Servicer //TODO: Add other actor types here

	// Get the key-value module store using the store key (in our case store key is "chain")
	store := ctx.KVStore(k.storeKey)

	servStore := prefix.NewStore(store, []byte(types.ServicerPrefix)) //TODO: Add other actor types here

	logger := ctx.Logger()

	// Paginate the recipes store based on PageRequest
	pageRes, err := query.Paginate(servStore, req.Pagination, func(key []byte, value []byte) error {
		var actor types.Servicer
		if err := k.cdc.Unmarshal(value, &actor); err != nil {
			logger.Error("could not unmarshal actor")
			return err
		}

		// Print out staker
		logger.Info(actor.String())

		actors = append(actors, &actor)

		return nil
	})

	// Throw an error if pagination failed
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Return a struct containing a list of recipes and pagination info
	return &types.QueryActorsResponse{Actor: actors, Pagination: pageRes}, nil
}
