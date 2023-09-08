package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"poktroll/x/poktroll/types"
)

func (k msgServer) SubmitProof(goCtx context.Context, msg *types.MsgSubmitProof) (*types.MsgSubmitProofResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgSubmitProofResponse{}, nil
}
