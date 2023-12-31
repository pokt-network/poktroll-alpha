package keeper

import (
	"testing"

	"poktroll/x/application/keeper"
	"poktroll/x/application/types"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"

	// "github.com/cosmos/cosmos-sdk/x/group/testutil"
	"github.com/golang/mock/gomock"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtestutil "github.com/cosmos/cosmos-sdk/x/distribution/testutil"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"
)

func ApplicationKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
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

	mockBankKeeper := distrtestutil.NewMockBankKeeper(ctrl)
	mockBankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(), types.ModuleName, gomock.Any(), gomock.Any()).AnyTimes()
	mockBankKeeper.EXPECT().SendCoinsFromAccountToModule(gomock.Any(), gomock.Any(), types.ModuleName, gomock.Any()).AnyTimes()

	// mockBankKeeper.
	paramsSubspace := typesparams.NewSubspace(cdc,
		types.Amino,
		storeKey,
		memStoreKey,
		"ApplicationParams",
	)
	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		paramsSubspace,
		mockBankKeeper,
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())

	return k, ctx
}
