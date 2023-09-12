package keeper_test

import (
	"strconv"
	"testing"

	keepertest "poktroll/testutil/keeper"
	"poktroll/testutil/nullify"
	"poktroll/x/application/keeper"
	"poktroll/x/application/types"

	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNApplication(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Application {
	addresses := simtestutil.CreateRandomAccounts(n)

	applications := make([]types.Application, n)
	for i := range applications {
		coins := sdk.NewCoin("stake", sdk.NewInt(int64(i)))
		applications[i].Address = addresses[i].String()
		applications[i].Stake = &coins

		keeper.SetApplication(ctx, applications[i])
	}
	return applications
}

func TestApplicationGet(t *testing.T) {
	keeper, ctx := keepertest.ApplicationKeeper(t)
	items := createNApplication(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetApplication(ctx,
			item.Address,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestApplicationRemove(t *testing.T) {
	keeper, ctx := keepertest.ApplicationKeeper(t)
	items := createNApplication(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveApplication(ctx,
			item.Address,
		)
		_, found := keeper.GetApplication(ctx,
			item.Address,
		)
		require.False(t, found)
	}
}

func TestApplicationGetAll(t *testing.T) {
	keeper, ctx := keepertest.ApplicationKeeper(t)
	items := createNApplication(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllApplication(ctx)),
	)
}
