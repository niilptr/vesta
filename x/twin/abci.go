// BeginBlock contains the logic that is automatically triggered at the beginning of each block
package twin

import (
	"vesta/x/twin/keeper"
	"vesta/x/twin/processor"
	"vesta/x/twin/types"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (am AppModule) BeginBlock(ctx sdk.Context, _ abci.RequestBeginBlock) {

	trainingState, _ := am.keeper.GetTrainingState(ctx)
	isTraining := trainingState.Value

	if isTraining {
		am.HandleTrainingResults(ctx, trainingState)
	}
}

// EndBlock contains the logic that is automatically triggered at the end of each block
func (am AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

func (am AppModule) HandleTrainingResults(ctx sdk.Context, trainingState types.TrainingState) {
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
		reasonWhyNotValid := ""
		for !isResultValid {
			idx, trainerMoniker, newTwinHash := p.GetBestTrainingResult(vtr)
			isResultValid, reasonWhyNotValid, err = p.ValidateBestTrainingResult(trainingState.TwinName, trainerMoniker, newTwinHash)
			if err != nil {
				p.Logger.Error(err.Error())
				return
			}

			if isResultValid {
				p.BroadcastResultIsValid(trainingState)
				am.keeper.UpdateTwinFromVestaTraining(ctx, trainingState.TwinName, newTwinHash)

			} else {
				if err != nil {
					p.Logger.Error(err.Error())
				} else {
					p.Logger.Error(reasonWhyNotValid)
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
