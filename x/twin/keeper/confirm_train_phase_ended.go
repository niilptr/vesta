package keeper

import (
	"vesta/x/twin/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) AddTrainingPhaseEndedConfirmation(ctx sdk.Context, signer string) error {

	ts, found := k.GetTrainingState(ctx)

	// Cannot be possible that training state is not found because a train request would have
	// initialized it.
	if !found {
		return types.ErrTrainingStateNotFound
	}

	// Cannot be possible that training state value (aka isTraining) is not set to true
	// (because it will be modified after majority of confirmations is enstablished).
	if !ts.Value {
		return types.ErrTrainingNotInProgress
	}

	ts.TrainingPhaseEndedConfirmations[signer] = true

	k.SetTrainingState(ctx, ts)

	return nil
}

func (k Keeper) CheckMajorityAgreesOnTrainingPhaseEnded(ctx sdk.Context, ts types.TrainingState, maxConfirmations uint32) bool {

	count := 0
	for _, value := range ts.TrainingPhaseEndedConfirmations {
		if value == true {
			count++
		}
	}

	if float32(count) < float32(maxConfirmations*2/3) {
		return false
	}

	return true
}
