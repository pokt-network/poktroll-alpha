package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"poktroll/x/session/types"
)

func (k msgServer) GetSession(goCtx context.Context, msg *types.MsgGetSession) (*types.MsgGetSessionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgGetSessionResponse{}, nil
}
