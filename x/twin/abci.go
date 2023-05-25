// BeginBlock contains the logic that is automatically triggered at the beginning of each block
package twin

import (
	"time"

	"vesta/x/twin/keeper"
	"vesta/x/twin/processor"
	"vesta/x/twin/types"

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
		NunComplete := 0
		NumTimeout := 0
		for _, v := range vts {
			if v.Complete {
				NunComplete++
			} else {
				if ctx.BlockTime().Sub(trainingState.StartTime) > am.keeper.GetParams(ctx).MaxWaitingTraining {
					NumTimeout++
				}
			}
		}

		if NunComplete+NumTimeout == len(vts) {
			newHash := processor.SetBestTrainingResults(trainingState.TwinName)
			am.keeper.UpdateTwinFromVestaTraining(ctx, trainingState.TwinName, newHash)
			am.keeper.SetTrainingState(ctx, types.TrainingState{
				Value:    false,
				TwinName: "",
			})
		}
	}
}

// EndBlock contains the logic that is automatically triggered at the end of each block
func (am AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
