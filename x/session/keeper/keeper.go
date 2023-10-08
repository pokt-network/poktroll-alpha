package keeper

//go:generate mockgen -destination ../../../testutil/mocks/session_keeper_mock.go -package mocks . SessionKeeper

import (
	"fmt"

	"cosmossdk.io/depinject"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"poktroll/x/session/types"
	sharedtypes "poktroll/x/shared/types"
)

type SessionKeeper interface {
	GetSessionForApp(ctx sdk.Context, appAddress string, serviceId string, blockHeight uint64) (*sharedtypes.Session, error)
}

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		paramstore paramtypes.Subspace

		appKeeper types.ApplicationKeeper
		svcKeeper types.ServicerKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,

	appKeeper types.ApplicationKeeper,
	svcKeeper types.ServicerKeeper,

) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	sessionKeeper := &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,

		appKeeper: appKeeper,
		svcKeeper: svcKeeper,
	}

	depConfig := supply(sessionKeeper)
	fmt.Printf("OLSH DEPINJECT CONFIG %+v\n", depConfig)
	if err := svcKeeper.InjectDeps(depConfig); err != nil {
		fmt.Printf("OLSH INJECTION FAILED %+v\n", err)
		panic(err)
	} else {
		fmt.Printf("OLSH INJECTION SUCCESSFUL\n")
	}

	return sessionKeeper
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// func supply(k *Keeper) depinject.Config {
// 	return depinject.Configs(
// 		// *NOTE*: May not need to use `depinject.BindInterface()`; this part
// 		// seems to be working.
// 		//
// 		// Not sure if using `depinject.BindInterface()` correctly:
// 		// - The servicer keeper is referencing
// 		// the interface defined in its expected_keepers.go (i.e. NOT the
// 		// `SessionKeeper` interface defined above).
// 		// - "Supplying" the session keeper here as a (concrete) `*Keeper` type.
// 		// - See: https://github.com/cosmos/cosmos-sdk/tree/main/depinject#bindinterface-api
// 		//depinject.BindInterface(
// 		//	// NB: pkg_path/pkg_name.InterfaceName
// 		//	"poktroll/x/servicer/types/types.SessionKeeper",
// 		//	"poktroll/x/session/keeper/keeper.Keeper",
// 		//),
// 		depinject.Supply(k),
// 	)
// }

func supply(k *Keeper) depinject.Config {
	return depinject.Configs(
		depinject.Supply(SessionKeeper(k)),
	)
	// depinject.BindInterface("types.SessionKeeper", "SessionKeeper"),
	// depinject.Provide(k)) // k.GetSessionForApp,
	// k.GetSession,

}
