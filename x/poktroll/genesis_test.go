package poktroll_test

import (
	"testing"

	"poktroll/testutil/nullify"
	"poktroll/x/poktroll"
	"poktroll/x/poktroll/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.poktrollKeeper(t)
	poktroll.InitGenesis(ctx, *k, genesisState)
	got := poktroll.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
