package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"poktroll/x/servicer/types"
)

// InsertClaim inserts the given claim into the state tree.
func (k Keeper) InsertClaim(ctx sdk.Context, claim *types.MsgClaim) error {
	// TODO_CONSIDERATION: do we want to re-use the servicer store for claims or
	// create a new "claims store"?
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClaimsKeyPrefix))
	claimBz, err := k.cdc.Marshal(claim)
	if err != nil {
		return err
	}

	claimKey := fmt.Sprintf("%s/%s", claim.Servicer, claim.SmtRootHash)
	store.Set([]byte(claimKey), claimBz)

	// TODO_CONSIDERATION: maybe add a `#String()` method to the claim message.
	hexSmstRootHash := fmt.Sprintf("%x", claim.SmtRootHash)
	event := sdk.NewEvent(
		EventTypeClaim,
		sdk.NewAttribute(AttributeKeySmtRootHash, hexSmstRootHash),
	)

	// HACK/IMPROVE: using "legacy" errors to save time; replace with custom error
	// protobuf types. See: https://docs.cosmos.network/v0.47/core/events.
	//
	// emit claim event
	ctx.EventManager().EmitEvent(event)
	return nil
}
