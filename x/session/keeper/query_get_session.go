package keeper

import (
	"context"
	"poktroll/x/session/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetSession(goCtx context.Context, req *types.QueryGetSessionRequest) (*types.QueryGetSessionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	session, err := k.GetSessionForApp(ctx, req.AppAddress)
	if err != nil {
		return nil, err
	}

	return &types.QueryGetSessionResponse{Session: *session}, nil
}
