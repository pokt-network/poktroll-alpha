package application_test

import (
	"math/rand"
	keepertest "poktroll/testutil/keeper"
	"poktroll/testutil/nullify"
	"poktroll/x/application"
	"poktroll/x/application/types"
	"testing"
	"time"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	accs := simtypes.RandomAccounts(r, 2)

	coin1 := sdk.NewCoin("upokt", math.NewInt(1))
	coin2 := sdk.NewCoin("upokt", math.NewInt(12))

	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		ApplicationList: []types.Application{
			{
				Address: accs[0].Address.String(),
				Stake:   &coin1,
			},
			{
				Address: accs[1].Address.String(),
				Stake:   &coin2,
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.ApplicationKeeper(t)
	application.InitGenesis(ctx, *k, genesisState)
	got := application.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.ApplicationList, got.ApplicationList)
	// this line is used by starport scaffolding # genesis/test/assert
}
