package keeper

import (
	"testing"

	"poktroll/x/session/keeper"
	"poktroll/x/session/types"

	mocks "poktroll/testutil/mocks"
	apptypes "poktroll/x/application/types"
	svctypes "poktroll/x/servicer/types"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func SessionKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	return SessionKeeperWithMocks(t, nil, nil)
}

// INVESTIGATE: Since the keepers injected into other keepers are private, we can't use
// `WithOption` like paradigm (like we do in the original v1 repo). We would need to break pattern
// from cosmos to do this, and it requires some consideration before we go down that route.
func SessionKeeperWithMocks(
	t testing.TB,
	// IMPROVE: Look into creating options to inject these BECAUSE they need to be customized per test
	appKeeper types.ApplicationKeeper,
	svcKeeper types.ServicerKeeper,
) (*keeper.Keeper, sdk.Context) {
	ctrl := gomock.NewController(t)

	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	paramsSubspace := typesparams.NewSubspace(cdc,
		types.Amino,
		storeKey,
		memStoreKey,
		"SessionParams",
	)

	if appKeeper == nil {
		appKeeper = defaultApplicationKeeper(ctrl)
	}
	if svcKeeper == nil {
		svcKeeper = defaultServicerKeeper(ctrl)
	}

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		paramsSubspace,
		appKeeper,
		svcKeeper,
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())

	return k, ctx
}

func defaultApplicationKeeper(ctrl *gomock.Controller) types.ApplicationKeeper {
	mockApplicationKeeper := mocks.NewMockApplicationKeeper(ctrl)
	mockApplicationKeeper.EXPECT().GetApplication(gomock.Any(), gomock.Any()).Return(apptypes.Application{}, true).AnyTimes()
	return mockApplicationKeeper
}

func defaultServicerKeeper(ctrl *gomock.Controller) types.ServicerKeeper {
	mockServicerKeeper := mocks.NewMockServicerKeeper(ctrl)
	mockServicerKeeper.EXPECT().GetAllServicers(gomock.Any()).Return([]svctypes.Servicers{}).AnyTimes()
	return mockServicerKeeper
}
