package keeper

//go:generate mockgen -destination ../../../testutil/mocks/servicer_keeper_mock.go -package mocks . ServicerKeeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"poktroll/x/servicer/types"
)

type ServicerKeeper interface {
	SetServicers(ctx sdk.Context, servicers types.Servicers)
	GetServicers(ctx sdk.Context, address string)
	RemoveServicers(ctx sdk.Context, address string)
	GetAllServicers(ctx sdk.Context) (list []types.Servicers)
}

type (
	Keeper struct {
		cdc               codec.BinaryCodec
		storeKey          storetypes.StoreKey
		memKey            storetypes.StoreKey
		paramstore        paramtypes.Subspace
		bankKeeper        types.BankKeeper
		accountKeeper     types.AccountKeeper
		applicationKeeper types.ApplicationKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	bk types.BankKeeper,
	ak types.AccountKeeper,
	appk types.ApplicationKeeper,

) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:               cdc,
		storeKey:          storeKey,
		memKey:            memKey,
		paramstore:        ps,
		bankKeeper:        bk,
		accountKeeper:     ak,
		applicationKeeper: appk,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
