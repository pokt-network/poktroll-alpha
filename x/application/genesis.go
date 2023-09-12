package application

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"poktroll/x/application/keeper"
	"poktroll/x/application/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the application
	for _, elem := range genState.ApplicationList {
		k.SetApplication(ctx, elem)
	}
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.ApplicationList = k.GetAllApplication(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
