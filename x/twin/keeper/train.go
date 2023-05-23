package keeper

import (
	"vesta/x/twin/processor"
	"vesta/x/twin/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) StartTraining(ctx sdk.Context, name string, creator string) error {

	isTraining := k.GetTrainingStateValue(ctx)

	if isTraining {
		return types.ErrTrainingInProgress
	}

	k.SetTrainingStateValue(ctx, true)

	// Run the local script to train the digital twin.
	p := processor.NewProcessor(k.GetNodeHome(), k.Logger(ctx))
	vtd, err := p.PrepareTraining(ctx, name)
	if err == nil {
		lr, err := vtd.Lr.Float64()
		if err == nil {
			go p.StartTraining(ctx, lr)
		}
	}

	return nil
}
