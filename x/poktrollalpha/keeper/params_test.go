package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	testkeeper "poktroll-alpha/testutil/keeper"
	"poktroll-alpha/x/poktrollalpha/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.PoktrollalphaKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
