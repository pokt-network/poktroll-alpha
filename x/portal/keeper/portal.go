package keeper

import (
	apptypes "poktroll/x/application/types"
	"poktroll/x/portal/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetPortal set a specific portal in the store from its index
func (k Keeper) SetPortal(ctx sdk.Context, portals types.Portal) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PortalKeyPrefix))
	b := k.cdc.MustMarshal(&portals)
	store.Set(types.PortalKey(
		portals.Address,
	), b)
}

// GetPortal returns a portal from its index
func (k Keeper) GetPortal(
	ctx sdk.Context,
	address string,

) (val types.Portal, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PortalKeyPrefix))
	b := store.Get(types.PortalKey(
		address,
	))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemovePortal removes a portal from the store
func (k Keeper) RemovePortal(
	ctx sdk.Context,
	address string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PortalKeyPrefix))
	store.Delete(types.PortalKey(
		address,
	))
}

// GetAllPortals returns all portals
func (k Keeper) GetAllPortals(ctx sdk.Context) (list []types.Portal) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PortalKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.Portal
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}
	return
}

// SetDelegatedApplication set a specific application's delegated portals in the store from its index
func (k Keeper) SetDelegator(ctx sdk.Context, appAddress string, delegatedPortals apptypes.Delegatees) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PortalDelegationsKeyPrefix))
	b := k.cdc.MustMarshal(&delegatedPortals)
	store.Set(types.PortalDelegationsKey(
		appAddress,
	), b)
}

// GetDelegatedPortals returns a application's delegated portals from its index
func (k Keeper) GetDelegatees(ctx sdk.Context, appAddress string) (val apptypes.Delegatees, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PortalDelegationsKeyPrefix))
	b := store.Get(types.PortalDelegationsKey(
		appAddress,
	))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// WhitelistApp updates the portal state appending the application address to its whitelist
func (k Keeper) WhitelistApp(ctx sdk.Context, portalAddress, appAddress string) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PortalKeyPrefix))
	b := store.Get(types.PortalKey(
		portalAddress,
	))
	if b == nil {
		return types.ErrPortalNotFound
	}
	portal := new(types.Portal)
	k.cdc.MustUnmarshal(b, portal)
	portal.WhitelistedApps.AppAddresses = append(portal.WhitelistedApps.AppAddresses, appAddress)
	b = k.cdc.MustMarshal(portal)
	store.Set(types.PortalKey(
		portalAddress,
	), b)
	return nil
}

// UnwhitelistApp updates the portal state removing the application address from its whitelist
func (k Keeper) UnwhitelistApp(ctx sdk.Context, portalAddress, appAddress string) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PortalKeyPrefix))
	b := store.Get(types.PortalKey(
		portalAddress,
	))
	if b == nil {
		return types.ErrPortalNotFound
	}
	portal := new(types.Portal)
	k.cdc.MustUnmarshal(b, portal)
	idx := -1
	for i, a := range portal.WhitelistedApps.AppAddresses {
		if a == appAddress {
			idx = i
			break
		}
	}
	if idx == -1 {
		return types.ErrAppNotWhitelisted
	}
	portal.WhitelistedApps.AppAddresses = append(portal.WhitelistedApps.AppAddresses[:idx], portal.WhitelistedApps.AppAddresses[idx+1:]...)
	b = k.cdc.MustMarshal(portal)
	store.Set(types.PortalKey(
		portalAddress,
	), b)
	return nil
}

// GetWhitelist returns the portal's whitelist
func (k Keeper) GetWhitelist(ctx sdk.Context, portalAddress string) (val []string, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PortalKeyPrefix))
	b := store.Get(types.PortalKey(
		portalAddress,
	))
	if b == nil {
		return nil, false
	}
	portal := new(types.Portal)
	k.cdc.MustUnmarshal(b, portal)
	return portal.WhitelistedApps.AppAddresses, true
}
