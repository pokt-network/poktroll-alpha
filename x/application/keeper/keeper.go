package keeper

//go:generate mockgen -destination ../../../testutil/mocks/application_keeper_mock.go -package mocks . ApplicationKeeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"poktroll/x/application/types"
)

type ApplicationKeeper interface {
	SetApplication(ctx sdk.Context, application types.Application)
	GetApplication(ctx sdk.Context, address string) (val types.Application, found bool)
	RemoveApplication(ctx sdk.Context, address string)
	GetAllApplication(ctx sdk.Context) (list []types.Application)
	DelegatePortal(ctx sdk.Context, appAddress string, portalPubKey cryptotypes.PubKey) error
	UndelagatePortal(ctx sdk.Context, appAddress string, portalPubKey cryptotypes.PubKey) error
	BurnCoins(ctx sdk.Context, appAddress sdk.AccAddress, amount sdk.Coin) error
}

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeKey     storetypes.StoreKey
		memKey       storetypes.StoreKey
		paramstore   paramtypes.Subspace
		bankKeeper   types.BankKeeper
		authKeeper   types.AccountKeeper
		portalKeeper types.PortalKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	bk types.BankKeeper,
	ak types.AccountKeeper,
	pk types.PortalKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:          cdc,
		storeKey:     storeKey,
		memKey:       memKey,
		paramstore:   ps,
		bankKeeper:   bk,
		authKeeper:   ak,
		portalKeeper: pk,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
