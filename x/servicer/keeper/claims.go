package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"poktroll/x/servicer/types"
)

// InsertClaim inserts the given claim into the state tree.
func (k Keeper) InsertClaim(ctx sdk.Context, claim *types.Claim) error {
	// TODO_CONSIDERATION: do we want to re-use the servicer store for claims or
	// create a new "claims store"?
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClaimsKeyPrefix))
	claimBz, err := k.cdc.Marshal(claim)
	if err != nil {
		return err
	}

	claimKey := fmt.Sprintf("%s", claim.GetSessionId())
	store.Set([]byte(claimKey), claimBz)
	return nil
}

func (k Keeper) GetClaim(ctx sdk.Context, sessionId string) (*types.Claim, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClaimsKeyPrefix))
	claimKey := fmt.Sprintf("%s", sessionId)
	claimBz := store.Get([]byte(claimKey))

	if claimBz == nil {
		return nil, fmt.Errorf("claim not found for sessionId: %s", sessionId)
	}

	var claim types.Claim
	if err := claim.Unmarshal(claimBz); err != nil {
		return nil, err
	}
	return &claim, nil
}
