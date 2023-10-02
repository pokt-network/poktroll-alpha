package portal

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"poktroll/testutil/sample"
	portalsimulation "poktroll/x/portal/simulation"
	"poktroll/x/portal/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = portalsimulation.FindAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
	_ = rand.Rand{}
)

const (
	opWeightMsgStakePortal = "op_weight_msg_stake_portal"
	// TODO: Determine the simulation weight value
	defaultWeightMsgStakePortal int = 100

	opWeightMsgUnstakePortal = "op_weight_msg_unstake_portal"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUnstakePortal int = 100

	opWeightMsgWhitelistApplication = "op_weight_msg_whitelist_application"
	// TODO: Determine the simulation weight value
	defaultWeightMsgWhitelistApplication int = 100

	opWeightMsgUnwhitelistApplication = "op_weight_msg_unwhitelist_application"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUnwhitelistApplication int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	portalGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&portalGenesis)
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

	var weightMsgStakePortal int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgStakePortal, &weightMsgStakePortal, nil,
		func(_ *rand.Rand) {
			weightMsgStakePortal = defaultWeightMsgStakePortal
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgStakePortal,
		portalsimulation.SimulateMsgStakePortal(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUnstakePortal int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUnstakePortal, &weightMsgUnstakePortal, nil,
		func(_ *rand.Rand) {
			weightMsgUnstakePortal = defaultWeightMsgUnstakePortal
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUnstakePortal,
		portalsimulation.SimulateMsgUnstakePortal(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgWhitelistApplication int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgWhitelistApplication, &weightMsgWhitelistApplication, nil,
		func(_ *rand.Rand) {
			weightMsgWhitelistApplication = defaultWeightMsgWhitelistApplication
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgWhitelistApplication,
		portalsimulation.SimulateMsgWhitelistApplication(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUnwhitelistApplication int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUnwhitelistApplication, &weightMsgUnwhitelistApplication, nil,
		func(_ *rand.Rand) {
			weightMsgUnwhitelistApplication = defaultWeightMsgUnwhitelistApplication
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUnwhitelistApplication,
		portalsimulation.SimulateMsgUnwhitelistApplication(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgStakePortal,
			defaultWeightMsgStakePortal,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				portalsimulation.SimulateMsgStakePortal(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUnstakePortal,
			defaultWeightMsgUnstakePortal,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				portalsimulation.SimulateMsgUnstakePortal(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgWhitelistApplication,
			defaultWeightMsgWhitelistApplication,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				portalsimulation.SimulateMsgWhitelistApplication(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUnwhitelistApplication,
			defaultWeightMsgUnwhitelistApplication,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				portalsimulation.SimulateMsgUnwhitelistApplication(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
