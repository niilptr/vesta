package twin

import (
	"vesta/x/twin/keeper"
	"vesta/x/twin/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the twin
	for _, elem := range genState.TwinList {
		k.SetTwin(ctx, elem)
	}
	// Set if defined
	if genState.TrainingState != nil {
		k.SetTrainingState(ctx, *genState.TrainingState)
	}
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.TwinList = k.GetAllTwin(ctx)
	// Get all training
	training, found := k.GetTrainingState(ctx)
	if found {
		genesis.TrainingState = &training
	}
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
