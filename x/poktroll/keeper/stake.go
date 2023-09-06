package keeper

import (
	"poktroll/x/poktroll/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) StakeActor(ctx sdk.Context, addr sdk.ValAddress, amount sdk.Coin) error {
	// TODO: sends coins to the staking module's pool!

	// Update store with new staking amount
	store := ctx.KVStore(k.storeKey)
	servStore := prefix.NewStore(store, []byte(types.ServicerPrefix))

	// Convert the valAddr to bytes as it will be the key to store validator info
	byteKey := addr.Bytes()

	// Get existing validator data from store
	bz := servStore.Get(byteKey)

	var servicer types.Servicer
	if bz == nil {
		// Create a new Validator object if not found
		servicer = types.NewServicer(addr, amount)
	} else {
		// Deserialize the byte array into a Validator object
		k.cdc.Unmarshal(bz, &servicer)

		// Update staking amount
		servicer.StakeAmount = servicer.StakeAmount.Add(amount)
	}

	// Serialize the Validator object back to bytes
	bz, err := k.cdc.Marshal(&servicer)
	if err != nil {
		return err
	}

	// Save the updated Validator object back to the store
	servStore.Set(byteKey, bz)

	return nil
}
