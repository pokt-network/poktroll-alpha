package keeper

//go:generate mockgen -destination ../../../testutil/mocks/servicer_keeper_mock.go -package mocks . ServicerKeeper

import (
	"fmt"

	"cosmossdk.io/depinject"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"poktroll/x/servicer/types"
	sharedtypes "poktroll/x/shared/types"
)

type ServicerKeeper interface {
	SetServicers(ctx sdk.Context, servicers sharedtypes.Servicers)
	GetServicers(ctx sdk.Context, address string)
	RemoveServicers(ctx sdk.Context, address string)
	GetAllServicers(ctx sdk.Context) (list []sharedtypes.Servicers)
	Inject(depinject.Config) error
}

type (
	Keeper struct {
		cdc           codec.BinaryCodec
		storeKey      storetypes.StoreKey
		memKey        storetypes.StoreKey
		paramstore    paramtypes.Subspace
		bankKeeper    types.BankKeeper
		sessionKeeper types.SessionKeeper
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
		// NB: sessionKeeper is provided via `depinject`
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) Inject(config depinject.Config) error {
	return depinject.Inject(config, &k.sessionKeeper)
}
