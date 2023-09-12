package servicer_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "poktroll/testutil/keeper"
	"poktroll/testutil/nullify"
	"poktroll/x/servicer"
	"poktroll/x/servicer/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.ServicerKeeper(t)
	servicer.InitGenesis(ctx, *k, genesisState)
	got := servicer.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
