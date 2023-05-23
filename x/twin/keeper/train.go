package keeper

import (
	"vesta/x/twin/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) StartTraining(ctx sdk.Context, name string, creator string) error {

	k.SetTrainingState(ctx, types.TrainingState{TrainingState: true})

	return nil
}
