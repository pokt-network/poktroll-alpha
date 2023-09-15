package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"poktroll/x/servicer/types"
)

func (k msgServer) Proof(goCtx context.Context, msg *types.MsgProof) (*types.MsgProofResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgProofResponse{}, nil
}
