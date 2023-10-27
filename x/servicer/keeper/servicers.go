package keeper

import (
	"poktroll/x/servicer/types"
	sharedtypes "poktroll/x/shared/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetServicers set a specific servicers in the store from its index
func (k Keeper) SetServicers(ctx sdk.Context, servicers sharedtypes.Servicers) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ServicersKeyPrefix))
	b := k.cdc.MustMarshal(&servicers)
	store.Set(types.ServicersKey(
		servicers.Address,
	), b)
}

// GetServicers returns a servicers from its index
func (k Keeper) GetServicers(
	ctx sdk.Context,
	address string,

) (val sharedtypes.Servicers, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ServicersKeyPrefix))

	b := store.Get(types.ServicersKey(
		address,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveServicers removes a servicers from the store
func (k Keeper) RemoveServicers(
	ctx sdk.Context,
	address string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ServicersKeyPrefix))
	store.Delete(types.ServicersKey(
		address,
	))
}

// GetAllServicers returns all servicers
func (k Keeper) GetAllServicers(ctx sdk.Context) (list []sharedtypes.Servicers) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ServicersKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val sharedtypes.Servicers
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
