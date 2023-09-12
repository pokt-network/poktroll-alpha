package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "poktroll/testutil/keeper"
	"poktroll/testutil/nullify"
	"poktroll/x/application/keeper"
	"poktroll/x/application/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNApplication(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Application {
	items := make([]types.Application, n)
	for i := range items {
		items[i].Index = strconv.Itoa(i)

		keeper.SetApplication(ctx, items[i])
	}
	return items
}

func TestApplicationGet(t *testing.T) {
	keeper, ctx := keepertest.ApplicationKeeper(t)
	items := createNApplication(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetApplication(ctx,
			item.Index,
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
			item.Index,
		)
		_, found := keeper.GetApplication(ctx,
			item.Index,
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
