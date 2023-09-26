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

func (k Keeper) Claims(goCtx context.Context, req *types.QueryClaimsRequest) (*types.QueryClaimsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	var claims []types.MsgClaim
	// TODO_THIS_COMMIT: move into claims.go
	claimsStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClaimsKeyPrefix))
	pageRes, err := query.FilteredPaginate(claimsStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var claim types.MsgClaim
		if err := k.cdc.Unmarshal(value, &claim); err != nil {
			return false, err
		}

		if claim.ServicerAddress == req.ServicerAddress {
			if accumulate {
				claims = append(claims, claim)
			}
			return true, nil
		}

		return false, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryClaimsResponse{Claims: claims, Pagination: pageRes}, nil
}
