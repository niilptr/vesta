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

	// Get how many authorized accounts are there. This is needed to verify if majority
	// agrees.
	numAuthorized := len(am.keeper.GetAuthorizedAccounts(ctx))

	// Get the training state.
	ts, found := am.keeper.GetTrainingState(ctx)
	if !found {
		return
	}

	isTraining := ts.Value

	if isTraining {

		agreement := am.keeper.CheckMajorityAgreesOnTrainingPhaseEnded(ctx, ts, uint32(numAuthorized))

		if agreement {
			ts = am.keeper.MustUpdateTrainingStateValue(ctx, ts, false)
			// TODO: emit event training phase complete

		} else {

			p, err := processor.NewProcessor(am.keeper.GetNodeHome(), am.keeper.Logger(ctx))
			if err != nil {
				return
			}

			confirmed := am.CheckIfTrainerAlreadyConfirmedTrainingPhaseEnded(ctx, ts, p)
			if !confirmed {
				am.HandleTrainingResults(ctx, ts, p)
			}

		}
	}

	ts, found = am.keeper.GetTrainingState(ctx)
	if !found {
		return
	}

	isValidating := ts.ValidationState.Value

	if isValidating {

		agreement, twinHash := am.keeper.CheckMajorityAgreesOnTrainingBestResult(ctx, ts, uint32(numAuthorized))

		if agreement {
			ts = am.keeper.MustUpdateTrainingStateValidationValue(ctx, ts, false)
			am.keeper.UpdateTwinFromVestaTraining(ctx, ts.TwinName, twinHash)
			// TODO: Emit event validation complete

		} else {

			p, err := processor.NewProcessor(am.keeper.GetNodeHome(), am.keeper.Logger(ctx))
			if err != nil {
				return
			}

			confirmed := am.CheckIfTrainerAlreadyConfirmedBestResult(ctx, ts, p)
			if !confirmed {
				am.HandleValidationPhase(ctx, ts)
			}
		}
	}
}

// EndBlock contains the logic that is automatically triggered at the end of each block
func (am AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

func (am AppModule) HandleTrainingResults(ctx sdk.Context, ts types.TrainingState, p processor.Processor) {

	vts, err := p.CheckValidatorsTrainingState(ts.TwinName, ts.TrainingConfigurationHash)
	if err != nil {
		p.Logger.Error(err.Error())
		return
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
		err := p.BroadcastConfirmationTrainingPhaseEnded()
		if err != nil {
			p.Logger.Error("Failed to broadcast confirmation training phase ended")
			return
		}
	}
}

func (am AppModule) HandleValidationPhase(ctx sdk.Context, ts types.TrainingState) {

	p, err := processor.NewProcessor(am.keeper.GetNodeHome(), keeper.ModuleLogger(ctx))
	if err != nil {
		return
	}

	vtr, err := p.ReadValidatorsTrainingResults(ts.TwinName, ts.TrainingConfigurationHash)
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
			err := p.BroadcastConfirmationBestResultIsValid(newTwinHash)
			if err != nil {
				p.Logger.Error("Failed to broadcast confirmation train best result")
				return
			}

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

func (am AppModule) CheckIfTrainerAlreadyConfirmedTrainingPhaseEnded(ctx sdk.Context, ts types.TrainingState, p processor.Processor) bool {

	address := p.GetAddress()
	for addr, confirmed := range ts.TrainingPhaseEndedConfirmations {
		if addr == address {
			if confirmed {
				return true
			}
		}
	}

	return false
}

func (am AppModule) CheckIfTrainerAlreadyConfirmedBestResult(ctx sdk.Context, ts types.TrainingState, p processor.Processor) bool {

	address := p.GetAddress()
	for addr, confirmed := range ts.ValidationState.MapValidatorsBestresulthash {
		if addr == address {
			if len(confirmed) > 0 {
				return true
			}
		}
	}

	return false
}
