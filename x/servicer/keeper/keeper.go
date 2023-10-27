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

// TODO: We need to add this for all modules
var _ ServicerKeeper = &Keeper{}

// TODO: We need to keep this updated to know which functions are exported
type ServicerKeeper interface {
	SetServicers(ctx sdk.Context, servicers sharedtypes.Servicers)
	GetServicers(ctx sdk.Context, address string) (servicers sharedtypes.Servicers, found bool)
	RemoveServicers(ctx sdk.Context, address string)
	GetAllServicers(ctx sdk.Context) (list []sharedtypes.Servicers)
	InjectDeps(depinject.Config) error
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
	// sk types.SessionKeeper,

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
		// sessionKeeper: sk, // NB: sessionKeeper is provided via `depinject`
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) InjectDeps(config depinject.Config) error {
	fmt.Printf("OLSH injecting deps before. svcKeeper: %p; sessionKeeper %p; sessionKeeperIsNil: %t \n\n", k, &k.sessionKeeper, k.sessionKeeper == nil)
	err := depinject.Inject(config, &k.sessionKeeper)
	fmt.Printf("OLSH injecting deps after. svcKeeper: %p; sessionKeeper %p; sessionKeeperIsNil: %t \n\n", k, &k.sessionKeeper, k.sessionKeeper == nil)
	return err
}
