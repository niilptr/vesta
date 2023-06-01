package twin

import (
	"math/rand"

	"vesta/testutil/sample"
	twinsimulation "vesta/x/twin/simulation"
	"vesta/x/twin/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = twinsimulation.FindAccount
	_ = simappparams.StakePerAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
	opWeightMsgCreateTwin = "op_weight_msg_twin"
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreateTwin int = 100

	opWeightMsgUpdateTwin = "op_weight_msg_twin"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateTwin int = 100

	opWeightMsgDeleteTwin = "op_weight_msg_twin"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteTwin int = 100

	opWeightMsgTrain = "op_weight_msg_train"
	// TODO: Determine the simulation weight value
	defaultWeightMsgTrain int = 100

	opWeightMsgConfirmTrainPhaseEnded = "op_weight_msg_confirm_train_phase_ended"
	// TODO: Determine the simulation weight value
	defaultWeightMsgConfirmTrainPhaseEnded int = 100

	opWeightMsgConfirmBestTrainResultIs = "op_weight_msg_confirm_best_train_result_is"
	// TODO: Determine the simulation weight value
	defaultWeightMsgConfirmBestTrainResultIs int = 100
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	twinGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		TwinList: []types.Twin{
			{
				Creator: sample.AccAddress(),
				Name:    "0",
			},
			{
				Creator: sample.AccAddress(),
				Name:    "1",
			},
		},
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&twinGenesis)
}

// ProposalContents doesn't return any content functions for governance proposals
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized  param changes for the simulator
func (am AppModule) RandomizedParams(_ *rand.Rand) []simtypes.ParamChange {

	return []simtypes.ParamChange{}
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCreateTwin int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreateTwin, &weightMsgCreateTwin, nil,
		func(_ *rand.Rand) {
			weightMsgCreateTwin = defaultWeightMsgCreateTwin
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateTwin,
		twinsimulation.SimulateMsgCreateTwin(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateTwin int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateTwin, &weightMsgUpdateTwin, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateTwin = defaultWeightMsgUpdateTwin
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateTwin,
		twinsimulation.SimulateMsgUpdateTwin(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeleteTwin int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgDeleteTwin, &weightMsgDeleteTwin, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteTwin = defaultWeightMsgDeleteTwin
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteTwin,
		twinsimulation.SimulateMsgDeleteTwin(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgTrain int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgTrain, &weightMsgTrain, nil,
		func(_ *rand.Rand) {
			weightMsgTrain = defaultWeightMsgTrain
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgTrain,
		twinsimulation.SimulateMsgTrain(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgConfirmTrainPhaseEnded int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgConfirmTrainPhaseEnded, &weightMsgConfirmTrainPhaseEnded, nil,
		func(_ *rand.Rand) {
			weightMsgConfirmTrainPhaseEnded = defaultWeightMsgConfirmTrainPhaseEnded
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgConfirmTrainPhaseEnded,
		twinsimulation.SimulateMsgConfirmTrainPhaseEnded(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgConfirmBestTrainResultIs int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgConfirmBestTrainResultIs, &weightMsgConfirmBestTrainResultIs, nil,
		func(_ *rand.Rand) {
			weightMsgConfirmBestTrainResultIs = defaultWeightMsgConfirmBestTrainResultIs
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgConfirmBestTrainResultIs,
		twinsimulation.SimulateMsgConfirmBestTrainResultIs(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	return operations
}
