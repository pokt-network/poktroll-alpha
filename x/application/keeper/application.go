package keeper

import (
	"fmt"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocdc "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"poktroll/x/application/types"
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
func (k Keeper) DelegatePortal(ctx sdk.Context, appAddress, portalAddress string) error {
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

	// check against max delegated param
	maxPortals := k.GetParams(ctx).MaxDelegatedPortals
	if uint32(len(app.Delegatees.PubKeys)) >= maxPortals {
		return errors.Wrapf(types.ErrMaxDelegatedReached, fmt.Sprintf("delegated portals: %d, max: %d", len(app.Delegatees.PubKeys), k.GetParams(ctx).MaxDelegatedPortals))
	}
	// ensure the portal is not already present
	for _, p := range app.Delegatees.PubKeys {
		pub, err := anyToPubKey(p)
		if err != nil {
			return err
		}
		if portalAddress == publicKeyToAddress(pub) {
			return types.ErrPortalAlreadyDelegated
		}
	}

	// update application's delegated portals
	pub, err := k.addressToPublicKey(ctx, portalAddress)
	if err != nil {
		return err
	}
	anyPub, err := codectypes.NewAnyWithValue(pub)
	if err != nil {
		return err
	}
	app.Delegatees.PubKeys = append(app.Delegatees.PubKeys, *anyPub)
	b = k.cdc.MustMarshal(app)
	store.Set(types.ApplicationKey(
		app.Address,
	), b)
	// index delegated portals per app address for easy lookup for portals
	k.portalKeeper.SetDelegator(ctx, appAddress, app.Delegatees)

	return nil
}

// UndelegatePortal removes a portal from an application's delegated portals
func (k Keeper) UndelegatePortal(ctx sdk.Context, appAddress, portalAddress string) error {
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
	for i, p := range app.Delegatees.PubKeys {
		pub, err := anyToPubKey(p)
		if err != nil {
			return err
		}
		if equal := portalAddress == publicKeyToAddress(pub); equal {
			found = true
			index = i
		}
	}
	if !found {
		return types.ErrPortalNotDelegated
	}

	// remove portal from application's delegated portals
	app.Delegatees.PubKeys = append(
		app.Delegatees.PubKeys[:index],
		app.Delegatees.PubKeys[index+1:]...,
	)
	b = k.cdc.MustMarshal(app)
	store.Set(types.ApplicationKey(
		app.Address,
	), b)
	// reindex delegated portals per app address for easy lookup for portals
	k.portalKeeper.SetDelegator(ctx, appAddress, app.Delegatees)

	return nil
}

// addressToPublicKey converts a bech32 address string to a public key
func (k Keeper) addressToPublicKey(ctx sdk.Context, address string) (cryptotypes.PubKey, error) {
	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}
	return k.authKeeper.GetPubKey(ctx, addr)
}

// anyToPubKey converts an Any type to a cryptotypes.PubKey
func anyToPubKey(any codectypes.Any) (cryptotypes.PubKey, error) {
	reg := codectypes.NewInterfaceRegistry()
	cryptocdc.RegisterInterfaces(reg)
	cdc := codec.NewProtoCodec(reg)
	var pub cryptotypes.PubKey
	if err := cdc.UnpackAny(&any, &pub); err != nil {
		return nil, fmt.Errorf("any type [%+v] is not cryptotypes.PubKey: %w", any, err)
	}
	return pub, nil
}

// publicKeyToAddress converts a cryptotypes.PubKey to a bech32 address string
func publicKeyToAddress(publicKey cryptotypes.PubKey) string {
	return sdk.AccAddress(publicKey.Address()).String()
}

func (k Keeper) BurnCoins(ctx sdk.Context, appAddress string, amount sdk.Coins) error {
	// TODO_INVESTIGATE: k.GetApplication is not returning the application from this method.
	// Even if it is behaving correctly when called from the CLI or msg_server_{stake,unstake}
	//application, found := k.GetApplication(ctx, appAddress)
	//if !found {
	//	return types.ErrApplicationNotFound
	//}

	//coinAmount := amount[0]

	//if application.Stake.IsLT(coinAmount) {
	//	return types.ErrInsufficientStake
	//}

	//newStake := application.Stake.Sub(coinAmount)
	//application.Stake = &newStake

	//if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, amount); err != nil {
	//	return err
	//}

	return nil
}
