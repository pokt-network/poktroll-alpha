package keeper

import (
	"cosmossdk.io/errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocdc "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
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
	app := new(types.Application)
	k.cdc.MustUnmarshal(b, app)

	// ensure the app is whitelisted by the portal
	portalAddr, err := anyPkToAddr(portalPubKey)
	if err != nil {
		return err
	}
	portalWhitelist, found := k.portalKeeper.GetWhitelist(ctx, portalAddr)
	if !found {
		return errors.Wrapf(types.ErrPortalNotFound, fmt.Sprintf("portal [%s] not found", portalAddr))
	}
	if len(portalWhitelist) > 0 {
		whitelisted := false
		for _, p := range portalWhitelist {
			if p == appAddress {
				whitelisted = true
				break
			}
		}
		if !whitelisted {
			return errors.Wrapf(types.ErrAppNotWhitelisted, fmt.Sprintf("app [%s] is not whitelisted for portal [%s]", appAddress, portalAddr))
		}
	}

	// check against max delegated param
	maxPortals := k.GetParams(ctx).MaxDelegatedPortals
	if uint32(len(app.DelegatedPortals.PortalPubKeys)) >= maxPortals {
		return errors.Wrapf(types.ErrMaxDelegatedReached, fmt.Sprintf("delegated portals: %d, max: %d", len(app.DelegatedPortals.PortalPubKeys), k.GetParams(ctx).MaxDelegatedPortals))
	}
	// ensure the portal is not already present
	for _, p := range app.DelegatedPortals.PortalPubKeys {
		equal, err := anyPkEquality(p, portalPubKey)
		if err != nil {
			return err
		}
		if equal {
			return types.ErrPortalAlreadyDelegated
		}
	}

	// update application's delegated portals
	app.DelegatedPortals.PortalPubKeys = append(app.DelegatedPortals.PortalPubKeys, portalPubKey)
	b = k.cdc.MustMarshal(app)
	store.Set(types.ApplicationKey(
		app.Address,
	), b)
	// index delegated portals per app address for easy lookup for portals
	k.portalKeeper.SetDelegatedApplication(ctx, appAddress, app.DelegatedPortals)

	return nil
}

// UndelegatePortal removes a portal from an application's delegated portals
func (k Keeper) UndelegatePortal(ctx sdk.Context, appAddress string, portalPubKey codectypes.Any) error {
	// update current application's value in store
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ApplicationKeyPrefix))
	b := store.Get(types.ApplicationKey(
		appAddress,
	))
	if b == nil {
		// if app doesn't exist it cannot be staked
		return types.ErrApplicationNotFound
	}
	app := new(types.Application)
	k.cdc.MustUnmarshal(b, app)

	// ensure the portal is already present
	found := false
	index := -1
	for i, p := range app.DelegatedPortals.PortalPubKeys {
		equal, err := anyPkEquality(p, portalPubKey)
		if err != nil {
			return err
		}
		if equal {
			found = true
			index = i
		}
	}
	if !found {
		return types.ErrPortalNotDelegated
	}

	// remove portal from application's delegated portals
	app.DelegatedPortals.PortalPubKeys = append(
		app.DelegatedPortals.PortalPubKeys[:index],
		app.DelegatedPortals.PortalPubKeys[index+1:]...,
	)
	b = k.cdc.MustMarshal(app)
	store.Set(types.ApplicationKey(
		app.Address,
	), b)
	// reindex delegated portals per app address for easy lookup for portals
	k.portalKeeper.SetDelegatedApplication(ctx, appAddress, app.DelegatedPortals)

	return nil
}

// anyPkEquality checks if two Any types are cosmos.crypto.PubKey interfaces and whether they are equal
func anyPkEquality(pk1, pk2 codectypes.Any) (equal bool, err error) {
	reg := codectypes.NewInterfaceRegistry()
	cryptocdc.RegisterInterfaces(reg)
	cdc := codec.NewProtoCodec(reg)
	var pubI1 cryptotypes.PubKey
	if err := cdc.UnpackAny(&pk1, &pubI1); err != nil {
		return false, fmt.Errorf("portal public key [%+v] is not cryptotypes.PubKey: %w", pk1.GetValue(), err)
	}
	var pubI2 cryptotypes.PubKey
	if err := cdc.UnpackAny(&pk2, &pubI2); err != nil {
		return false, fmt.Errorf("portal public key [%+v] is not cryptotypes.PubKey: %w", pk1.GetValue(), err)
	}
	return pubI1.Equals(pubI2), nil
}

// pkToAddr converts a public key to a bech32 address string
func pkToAddr(pk cryptotypes.PubKey) string {
	return sdk.AccAddress(pk.Address()).String()
}

// anyPkToAddr converts an Any type to a bech32 address string
func anyPkToAddr(ak codectypes.Any) (string, error) {
	reg := codectypes.NewInterfaceRegistry()
	cryptocdc.RegisterInterfaces(reg)
	cdc := codec.NewProtoCodec(reg)
	var pub cryptotypes.PubKey
	if err := cdc.UnpackAny(&ak, &pub); err != nil {
		return "", fmt.Errorf("Any type [%+v] is not cryptotypes.PubKey: %w", ak.GetValue(), err)
	}
	return pkToAddr(pub), nil
}
