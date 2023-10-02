package keeper

//go:generate mockgen -destination ../../../testutil/mocks/portal_keeper_mock.go -package mocks . PortalKeeper

import (
	"fmt"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	apptypes "poktroll/x/application/types"

	"poktroll/x/portal/types"
)

type PortalKeeper interface {
	SetPortal(ctx sdk.Context, portals types.Portal)
	GetPortal(ctx sdk.Context, address string) (val types.Portal, found bool)
	RemovePortal(ctx sdk.Context, address string)
	GetAllPortals(ctx sdk.Context) (list []types.Portal)
	SetDelegatedApplication(ctx sdk.Context, appAddress string, delegatedPortals apptypes.DelegatedPortals)
	GetDelegatedPortals(ctx sdk.Context, appAddress string) (val apptypes.DelegatedPortals, found bool)
	WhitelistApp(ctx sdk.Context, portalAddress, appAddress string) error
	UnwhitelistApp(ctx sdk.Context, portalAddress, appAddress string) error
}

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		paramstore paramtypes.Subspace
		bankKeeper types.BankKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	bk types.BankKeeper,

) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,
		bankKeeper: bk,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
