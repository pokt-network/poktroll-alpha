package poktrollalpha_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "poktroll-alpha/testutil/keeper"
	"poktroll-alpha/testutil/nullify"
	"poktroll-alpha/x/poktrollalpha"
	"poktroll-alpha/x/poktrollalpha/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.PoktrollalphaKeeper(t)
	poktrollalpha.InitGenesis(ctx, *k, genesisState)
	got := poktrollalpha.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
