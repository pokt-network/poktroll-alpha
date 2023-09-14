package servicer

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"poktroll/testutil/sample"
	servicersimulation "poktroll/x/servicer/simulation"
	"poktroll/x/servicer/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = servicersimulation.FindAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
	_ = rand.Rand{}
)

const (
	opWeightMsgStakeServicer = "op_weight_msg_stake_servicer"
	// TODO: Determine the simulation weight value
	defaultWeightMsgStakeServicer int = 100

	opWeightMsgUnstakeServicer = "op_weight_msg_unstake_servicer"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUnstakeServicer int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	servicerGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&servicerGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgStakeServicer int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgStakeServicer, &weightMsgStakeServicer, nil,
		func(_ *rand.Rand) {
			weightMsgStakeServicer = defaultWeightMsgStakeServicer
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgStakeServicer,
		servicersimulation.SimulateMsgStakeServicer(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUnstakeServicer int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUnstakeServicer, &weightMsgUnstakeServicer, nil,
		func(_ *rand.Rand) {
			weightMsgUnstakeServicer = defaultWeightMsgUnstakeServicer
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUnstakeServicer,
		servicersimulation.SimulateMsgUnstakeServicer(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgStakeServicer,
			defaultWeightMsgStakeServicer,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				servicersimulation.SimulateMsgStakeServicer(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUnstakeServicer,
			defaultWeightMsgUnstakeServicer,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				servicersimulation.SimulateMsgUnstakeServicer(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
