package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"poktroll/x/servicer/types"
)

// InsertProof inserts the given Proof into the state tree.
func (k Keeper) InsertProof(ctx sdk.Context, proof *types.MsgProof) error {
	// TODO_CONSIDERATION: do we want to re-use the servicer store for Proofs or
	// create a new "Proofs store"?
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ProofsKeyPrefix))
	proofBz, err := k.cdc.Marshal(proof)
	if err != nil {
		return err
	}

	ProofKey := fmt.Sprintf("%s", proof.SessionId)
	store.Set([]byte(ProofKey), proofBz)
	return nil
}
