package servicer

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"poktroll/x/servicer/keeper"
	"poktroll/x/servicer/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the servicers
	for _, elem := range genState.ServicersList {
		k.SetServicers(ctx, elem)
	}
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.ServicersList = k.GetAllServicers(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
