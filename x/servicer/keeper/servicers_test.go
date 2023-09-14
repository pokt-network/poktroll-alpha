package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "poktroll/testutil/keeper"
	"poktroll/testutil/nullify"
	"poktroll/x/servicer/keeper"
	"poktroll/x/servicer/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNServicers(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Servicers {
	items := make([]types.Servicers, n)
	for i := range items {
		items[i].Index = strconv.Itoa(i)

		keeper.SetServicers(ctx, items[i])
	}
	return items
}

func TestServicersGet(t *testing.T) {
	keeper, ctx := keepertest.ServicerKeeper(t)
	items := createNServicers(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetServicers(ctx,
			item.Index,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestServicersRemove(t *testing.T) {
	keeper, ctx := keepertest.ServicerKeeper(t)
	items := createNServicers(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveServicers(ctx,
			item.Index,
		)
		_, found := keeper.GetServicers(ctx,
			item.Index,
		)
		require.False(t, found)
	}
}

func TestServicersGetAll(t *testing.T) {
	keeper, ctx := keepertest.ServicerKeeper(t)
	items := createNServicers(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllServicers(ctx)),
	)
}
