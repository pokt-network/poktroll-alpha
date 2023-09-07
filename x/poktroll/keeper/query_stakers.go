package keeper

import (
	"context"
	"fmt"

	"poktroll/x/poktroll/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Stakers(goCtx context.Context, req *types.QueryStakersRequest) (*types.QueryStakersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Define a variable that will store a list of servicers
	var stakers []*types.Staker

	// Get the key-value module store using the store key (in our case store key is "chain")
	store := ctx.KVStore(k.storeKey)

	servStore := prefix.NewStore(store, []byte(types.StakerPrefix))

	// Paginate the recipes store based on PageRequest
	pageRes, err := query.Paginate(servStore, req.Pagination, func(key []byte, value []byte) error {
		var staker types.Staker
		if err := k.cdc.Unmarshal(value, &staker); err != nil {
			return err
		}

		// Print out staker
		fmt.Println(staker)

		stakers = append(stakers, &staker)

		return nil
	})

	// Throw an error if pagination failed
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Return a struct containing a list of recipes and pagination info
	return &types.QueryStakersResponse{Staker: stakers, Pagination: pageRes}, nil
}
