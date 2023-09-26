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

	claimKey := fmt.Sprintf("%s", claim.SessionId)
	store.Set([]byte(claimKey), claimBz)
	return nil
}
