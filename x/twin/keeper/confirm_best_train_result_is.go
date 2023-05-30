package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"vesta/x/twin/types"
)

func (k Keeper) AddBestTrainResultToTrainingState(ctx sdk.Context, signer string, twinHash string) error {

	ts, found := k.GetTrainingState(ctx)

	// Cannot be possible that training state is not found because a train request would have
	// initialized it.
	if !found {
		return types.ErrTrainingStateNotFound
	}

	// Cannot be possible that training state value (aka isTraining) is set to true
	// (because validation phase comes after training phase ended).
	if ts.Value {
		return types.ErrTrainingInProgress
	}

	// Cannot be possible that validation state value (aka isValidating) is set to false
	// (because validation phase must be active before agreement on best result is reached).
	if !ts.ValidationState.Value {
		return types.ErrTrainingValidationNotInProgress
	}

	ts.ValidationState.MapValidatorsBestresulthash[signer] = twinHash

	k.SetTrainingState(ctx, ts)

	return nil
}

func (k Keeper) CheckMajorityAgreesOnTrainingBestResult(ctx sdk.Context, ts types.TrainingState, maxConfirmations uint32) (agreement bool, twinHash string) {

	countMap := make(map[string]uint32)

	for key := range ts.ValidationState.MapValidatorsBestresulthash {
		countMap[key] = countMap[key] + 1
	}

	var maxCount uint32 = 0
	mostReputableHash := ""

	for hash, count := range countMap {
		if count > maxCount {
			maxCount = count
			mostReputableHash = hash
		}
	}

	if float32(maxCount) < float32(maxConfirmations*2/3) {
		return false, mostReputableHash
	}

	return true, mostReputableHash
}
