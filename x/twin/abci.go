// BeginBlock contains the logic that is automatically triggered at the beginning of each block
package twin

import (
	"time"

	"vesta/x/twin/keeper"
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

		p := processor.NewProcessor(am.keeper.GetNodeHome(), keeper.ModuleLogger(ctx))

		vts, err := p.CheckValidatorsTrainingState(trainingState.TwinName)
		if err != nil {
			p.Logger.Error(err.Error())
		}

		nunComplete := 0
		numTimeout := 0
		for _, v := range vts {
			if v.Complete {
				nunComplete++
			} else {
				if ctx.BlockTime().Sub(trainingState.StartTime) > am.keeper.GetParams(ctx).MaxWaitingTraining {
					numTimeout++
				}
			}
		}

		// If all complete
		if nunComplete+numTimeout == len(vts) {

			am.keeper.SetTrainingStateValue(ctx, false)

			vtr, err := p.ReadValidatorsTrainingResults(trainingState.TwinName)
			if err != nil {
				p.Logger.Error(err.Error())
				return
			}

			isResultValid := false
			for !isResultValid {
				idx, trainerMoniker, newTwinHash := p.GetBestTrainingResult(vtr)
				isResultValid = p.ValidateTrainingResult(trainingState.TwinName, trainerMoniker)

				if isResultValid {
					am.keeper.UpdateTwinFromVestaTraining(ctx, trainingState.TwinName, newTwinHash)

				} else {
					// remove not valid result from result slice
					vtr = append(vtr[:idx], vtr[idx+1:]...)

					// if result slice is empty break
					if len(vtr) == 0 {
						p.Logger.Error("All training results are not valid.")
						break
					}
				}
			}
		}
	}
}

// EndBlock contains the logic that is automatically triggered at the end of each block
func (am AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
