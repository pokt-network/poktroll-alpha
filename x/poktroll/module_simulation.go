package poktroll

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"poktroll/testutil/sample"
	poktrollsimulation "poktroll/x/poktroll/simulation"
	"poktroll/x/poktroll/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = poktrollsimulation.FindAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
	_ = rand.Rand{}
)

const (
	opWeightMsgStake = "op_weight_msg_stake"
	// TODO: Determine the simulation weight value
	defaultWeightMsgStake int = 100

	opWeightMsgUnstake = "op_weight_msg_unstake"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUnstake int = 100

	// this line is used by starport scaffolding # simapp/module/const
	opWeightMsgSubmitClaim = "op_weight_msg_submit_claim"
	// TODO: Determine the simulation weight value
	defaultWeightMsgSubmitClaim int = 100

	opWeightMsgSubmitProof = "op_weight_msg_submit_proof"
	// TODO: Determine the simulation weight value
	defaultWeightMsgSubmitProof int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	poktrollGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&poktrollGenesis)
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

	var weightMsgStake int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgStake, &weightMsgStake, nil,
		func(_ *rand.Rand) {
			weightMsgStake = defaultWeightMsgStake
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgStake,
		poktrollsimulation.SimulateMsgStake(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUnstake int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUnstake, &weightMsgUnstake, nil,
		func(_ *rand.Rand) {
			weightMsgUnstake = defaultWeightMsgUnstake
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUnstake,
		poktrollsimulation.SimulateMsgUnstake(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgSubmitClaim int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgSubmitClaim, &weightMsgSubmitClaim, nil,
		func(_ *rand.Rand) {
			weightMsgSubmitClaim = defaultWeightMsgSubmitClaim
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSubmitClaim,
		poktrollsimulation.SimulateMsgSubmitClaim(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgSubmitProof int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgSubmitProof, &weightMsgSubmitProof, nil,
		func(_ *rand.Rand) {
			weightMsgSubmitProof = defaultWeightMsgSubmitProof
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSubmitProof,
		poktrollsimulation.SimulateMsgSubmitProof(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgStake,
			defaultWeightMsgStake,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				poktrollsimulation.SimulateMsgStake(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgUnstake,
			defaultWeightMsgUnstake,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				poktrollsimulation.SimulateMsgUnstake(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgSubmitClaim,
			defaultWeightMsgSubmitClaim,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				poktrollsimulation.SimulateMsgSubmitClaim(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgSubmitProof,
			defaultWeightMsgSubmitProof,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				poktrollsimulation.SimulateMsgSubmitProof(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
