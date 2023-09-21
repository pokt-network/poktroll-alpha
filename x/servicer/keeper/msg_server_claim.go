package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"poktroll/x/servicer/types"
)

func (k msgServer) Claim(goCtx context.Context, msg *types.MsgClaim) (*types.MsgClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	k.Logger(ctx).With("method", "Claim").Error("EMITTING EVENT CLAIMED")

	eventClaimed := &types.EventClaimed{
		SmtRootHash:     msg.SmtRootHash,
		ServicerAddress: msg.Creator,
	}
	if err := ctx.EventManager().EmitTypedEvent(eventClaimed); err != nil {
		return nil, err
	}

	event := sdk.NewEvent("EventClaimed",
		sdk.NewAttribute("smt_root_hash", string(msg.SmtRootHash)),
		sdk.NewAttribute("servicer_address", msg.Creator),
	)
	ctx.EventManager().EmitEvent(event)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgClaimResponse{}, nil
}
