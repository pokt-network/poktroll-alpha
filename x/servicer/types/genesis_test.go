package types_test

import (
	"math/rand"
	"testing"
	"time"

	"poktroll/x/servicer/types"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/stretchr/testify/require"
)

// TODO: Replace all `stake` denominations with `upokt` once we get it to start up correctly

func TestGenesisState_Validate(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	accs := simtypes.RandomAccounts(r, 2)

	coin1 := sdk.NewCoin("stake", math.NewInt(1))
	coin2 := sdk.NewCoin("stake", math.NewInt(12))

	tests := []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{

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
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicated servicers",
			genState: &types.GenesisState{
				ServicersList: []types.Servicers{
					{
						Address: accs[0].Address.String(),
						Stake:   &coin1,
					},
					{
						Address: accs[0].Address.String(),
						Stake:   &coin2,
					},
				},
			},
			valid: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
