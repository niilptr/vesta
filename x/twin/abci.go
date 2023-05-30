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

	ts, _ := am.keeper.GetTrainingState(ctx)
	isTraining := ts.Value
	isValidating := ts.ValidationState.Value

	if isTraining {
		am.HandleTrainingResults(ctx, ts)
	}

	if isValidating {
		// TODO: controlla quanti hanno broadcastato il best result
		//

		am.HandleValidationPhase(ctx, ts)
	}

}

// EndBlock contains the logic that is automatically triggered at the end of each block
func (am AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

func (am AppModule) HandleTrainingResults(ctx sdk.Context, ts types.TrainingState) {

	p := processor.NewProcessor(am.keeper.GetNodeHome(), keeper.ModuleLogger(ctx))

	vts, err := p.CheckValidatorsTrainingState(ts.TwinName)
	if err != nil {
		p.Logger.Error(err.Error())
	}

	nunComplete := 0
	numTimeout := 0
	for _, v := range vts {
		if v.Complete {
			nunComplete++
		} else {
			if ctx.BlockTime().Sub(ts.StartTime) > am.keeper.GetParams(ctx).MaxWaitingTraining {
				numTimeout++
			}
		}
	}

	// If all complete
	if nunComplete+numTimeout == len(vts) {
		am.keeper.SetTrainingStateValue(ctx, ts, false)
		am.keeper.SetTrainingStateValidationValue(ctx, ts, true)
	}
}

func (am AppModule) HandleValidationPhase(ctx sdk.Context, ts types.TrainingState) {

	p := processor.NewProcessor(am.keeper.GetNodeHome(), keeper.ModuleLogger(ctx))

	vtr, err := p.ReadValidatorsTrainingResults(ts.TwinName)
	if err != nil {
		p.Logger.Error(err.Error())
		return
	}

	isBestResultValid := false
	reasonWhyNotValid := ""
	for !isBestResultValid {
		idx, trainerMoniker, newTwinHash := p.GetBestTrainingResult(vtr)
		isBestResultValid, reasonWhyNotValid, err = p.ValidateBestTrainingResult(ts.TwinName, trainerMoniker, newTwinHash)
		if err != nil {
			p.Logger.Error(err.Error())
			return
		}

		if isBestResultValid {
			p.BroadcastBestResultIsValid(ts, newTwinHash)

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

func (am AppModule) HandleTwinUpdateFromVestaTraining(ctx sdk.Context, ts types.TrainingState, newTwinHash string) {
	am.keeper.UpdateTwinFromVestaTraining(ctx, ts.TwinName, newTwinHash)
}
