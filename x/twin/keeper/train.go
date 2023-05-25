package keeper

import (
	"vesta/x/twin/processor"
	"vesta/x/twin/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) StartTraining(ctx sdk.Context, twinName string, creator string, trainHash string) error {

	isTraining := k.GetTrainingStateValue(ctx)

	if isTraining {
		return types.ErrTrainingInProgress
	}

	k.SetTrainingState(ctx, types.TrainingState{
		Value:                     true,
		TwinName:                  twinName,
		StartTime:                 ctx.BlockTime(),
		TrainingConfigurationHash: trainHash,
	})

	// Run the local script to train the digital twin.
	// TODO: Check that TrainingConfigurationHash saved in the store corresponds to the hash of the file in the central db.
	p := processor.NewProcessor(k.GetNodeHome(), k.Logger(ctx))
	vtd, err := p.PrepareTraining(ctx, twinName)
	if err == nil {
		lr, err := vtd.Lr.Float64()
		if err == nil {
			go p.StartTraining(ctx, lr)
		}
	}

	return nil
}
