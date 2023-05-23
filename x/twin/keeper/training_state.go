package keeper

import (
	"vesta/x/twin/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetTraining set training in the store
func (k Keeper) SetTrainingState(ctx sdk.Context, training types.TrainingState) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TrainingStateKey))
	b := k.cdc.MustMarshal(&training)
	store.Set([]byte{0}, b)
}

// GetTraining returns training
func (k Keeper) GetTrainingState(ctx sdk.Context) (val types.TrainingState, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TrainingStateKey))

	b := store.Get([]byte{0})
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveTraining removes training from the store
func (k Keeper) RemoveTrainingState(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TrainingStateKey))
	store.Delete([]byte{0})
}
