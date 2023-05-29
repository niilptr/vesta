package keeper

import (
	"vesta/x/twin/processor"
	"vesta/x/twin/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) StartTraining(ctx sdk.Context, twinName string, creator string, trainConfHash string) error {

	isTraining := k.GetTrainingStateValue(ctx)

	if isTraining {
		return types.ErrTrainingInProgress
	}

	k.SetTrainingState(ctx, types.TrainingState{
		Value:                     true,
		TwinName:                  twinName,
		StartTime:                 ctx.BlockTime(),
		TrainingConfigurationHash: trainConfHash,
	})

	// Run the local script to train the digital twin.
	// TODO: Check that TrainingConfigurationHash saved in the store corresponds to the hash of the file in the central db.
	p := processor.NewProcessor(k.GetNodeHome(), k.Logger(ctx))
	_, err := p.PrepareTraining(ctx, twinName)
	if err == nil {
		if err == nil {
			go p.Train()
		}
	} else {
		p.Logger.Error("Local training not start.")
	}

	return nil
}
