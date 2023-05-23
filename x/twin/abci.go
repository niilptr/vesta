// BeginBlock contains the logic that is automatically triggered at the beginning of each block
package twin

import (
	"time"

	"vesta/x/twin/processor"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ValidatorsTrainingState struct {
	ValidatorsTrainingState []ValidatorTrainingState
}
type ValidatorTrainingState struct {
	HasCompleted bool
	StartTime    time.Time
}

func (am AppModule) BeginBlock(ctx sdk.Context, _ abci.RequestBeginBlock) {

	trainingState, _ := am.keeper.GetTrainingState(ctx)
	isTraining := trainingState.Value
	if isTraining {
		twinName := trainingState.Name
		ValidatorsTrainingState := processor.GetValidatorsTrainingState(twinName)
		completed := 0
		timeout := 0
		for _, vts := range ValidatorsTrainingState {
			if vts.HasCompleted {
				completed++
			} else {
				if ctx.BlockTime().Sub(vts.StartTime) > am.keeper.GetParams(ctx).MaxWaitingTraining {
					timeout++
				}
			}
		}

		if completed+timeout == len(ValidatorsTrainingState) {
			am.keeper.SetTrainingStateValue(ctx, false)
			newHash := processor.SetBestTrainingResults(twinName)
			am.keeper.UpdateTwinFromVestaTraining(ctx, twinName, newHash)
		}
	}
}

// EndBlock contains the logic that is automatically triggered at the end of each block
func (am AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
