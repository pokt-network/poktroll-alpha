package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/types/query"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"poktroll/x/servicer/types"
)

func (k Keeper) Proofs(goCtx context.Context, req *types.QueryProofsRequest) (*types.QueryProofsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	var Proofs []types.MsgProof
	// TODO_THIS_COMMIT: move into Proofs.go
	ProofsStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ProofsKeyPrefix))
	pageRes, err := query.FilteredPaginate(ProofsStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var Proof types.MsgProof
		if err := k.cdc.Unmarshal(value, &Proof); err != nil {
			return false, err
		}

		if Proof.ServicerAddress == req.ServicerAddress {
			if accumulate {
				Proofs = append(Proofs, Proof)
			}
			return true, nil
		}

		return false, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryProofsResponse{Proofs: Proofs, Pagination: pageRes}, nil
}
