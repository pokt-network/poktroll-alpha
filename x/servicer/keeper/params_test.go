package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	testkeeper "poktroll/testutil/keeper"
	"poktroll/x/servicer/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.ServicerKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
