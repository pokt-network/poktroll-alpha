package keeper

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"poktroll/x/application/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetApplication set a specific application in the store from its index
func (k Keeper) SetApplication(ctx sdk.Context, application types.Application) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ApplicationKeyPrefix))
	b := k.cdc.MustMarshal(&application)
	store.Set(types.ApplicationKey(
		application.Address,
	), b)
}

// GetApplication returns an application from its index
func (k Keeper) GetApplication(
	ctx sdk.Context,
	address string,

) (val types.Application, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ApplicationKeyPrefix))

	b := store.Get(types.ApplicationKey(
		address,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveApplication removes a application from the store
func (k Keeper) RemoveApplication(
	ctx sdk.Context,
	address string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ApplicationKeyPrefix))
	store.Delete(types.ApplicationKey(
		address,
	))
}

// GetAllApplication returns all application
func (k Keeper) GetAllApplication(ctx sdk.Context) (list []types.Application) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ApplicationKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Application
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// DelegatePortal delegates an application to a portal
func (k Keeper) DelegatePortal(ctx sdk.Context, appAddress string, portalPubKey codectypes.Any) error {
	// update current application's value in store
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ApplicationKeyPrefix))
	b := store.Get(types.ApplicationKey(
		appAddress,
	))
	if b == nil {
		// if app doesn't exist it cannot be staked
		return types.ErrApplicationNotFound
	}
	// update application's delegated portals
	app := new(types.Application)
	k.cdc.MustUnmarshal(b, app)
	app.DelegatedPortals.PortalPubKeys = append(app.DelegatedPortals.PortalPubKeys, portalPubKey)
	b = k.cdc.MustMarshal(app)
	store.Set(types.ApplicationKey(
		app.Address,
	), b)

	// index delegated portals per app address for easy lookup for portals
	k.portalKeeper.SetDelegatedApplication(ctx, appAddress, app.DelegatedPortals)

	return nil
}
