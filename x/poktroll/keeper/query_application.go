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

func (k Keeper) Application(goCtx context.Context, req *types.QueryApplicationRequest) (*types.QueryApplicationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := ctx.Logger()

	// Get the key-value module store using the store key (in our case store key is "chain")
	store := ctx.KVStore(k.storeKey)

	// TODO: Add other actor types here
	var applications []*types.Application
	applicationStore := prefix.NewStore(store, []byte(types.ApplicationPrefix))

	// Paginate the recipes store based on PageRequest
	pageRes, err := query.Paginate(applicationStore, req.Pagination, func(key []byte, value []byte) error {
		var application types.Application
		if err := k.cdc.Unmarshal(value, &application); err != nil {
			logger.Error("could not unmarshal application")
			return err
		}
		applications = append(applications, &application)
		return nil
	})

	// Throw an error if pagination failed
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Return a struct containing a list of recipes and pagination info
	return &types.QueryApplicationResponse{Applications: applications, Pagination: pageRes}, nil
}
