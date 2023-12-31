package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	testkeeper "poktroll/testutil/keeper"
	"poktroll/x/application/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.ApplicationKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
