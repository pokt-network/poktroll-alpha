package servicer_test

import (
	"math/rand"
	"testing"
	"time"

	keepertest "poktroll/testutil/keeper"
	"poktroll/testutil/nullify"
	"poktroll/x/servicer"
	"poktroll/x/servicer/types"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	accs := simtypes.RandomAccounts(r, 2)

	coin1 := sdk.NewCoin("stake", math.NewInt(1))
	coin2 := sdk.NewCoin("stake", math.NewInt(12))

	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		ServicersList: []types.Servicers{
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

	k, ctx := keepertest.ServicerKeeper(t)
	servicer.InitGenesis(ctx, *k, genesisState)
	got := servicer.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.ServicersList, got.ServicersList)
	// this line is used by starport scaffolding # genesis/test/assert
}
