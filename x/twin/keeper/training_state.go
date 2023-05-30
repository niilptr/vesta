package keeper

import (
	"vesta/x/twin/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetTrainingStateValue set trainingState value in the store
func (k Keeper) SetTrainingState(ctx sdk.Context, trainingState types.TrainingState) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TrainingStateKey))
	b := k.cdc.MustMarshal(&trainingState)
	store.Set([]byte{0}, b)
}

// SetTrainingStateValue set trainingState value in the store
func (k Keeper) SetTrainingStateValue(ctx sdk.Context, ts types.TrainingState, value bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TrainingStateKey))
	ts.Value = value
	b := k.cdc.MustMarshal(&ts)
	store.Set([]byte{0}, b)
}

// SetTrainingStateValue set trainingState value in the store
func (k Keeper) SetTrainingStateTwinName(ctx sdk.Context, ts types.TrainingState, twinName string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TrainingStateKey))
	ts.TwinName = twinName
	b := k.cdc.MustMarshal(&ts)
	store.Set([]byte{0}, b)
}

// SetTrainingStateValue set trainingState value in the store
func (k Keeper) SetTrainingStateValidationValue(ctx sdk.Context, ts types.TrainingState, value bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TrainingStateKey))
	ts.ValidationState.Value = value
	b := k.cdc.MustMarshal(&ts)
	store.Set([]byte{0}, b)
}

// GetTrainingState returns trainingState
func (k Keeper) GetTrainingState(ctx sdk.Context) (ts types.TrainingState, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TrainingStateKey))

	b := store.Get([]byte{0})
	if b == nil {
		return ts, false
	}

	k.cdc.MustUnmarshal(b, &ts)
	return ts, true
}

// GetTraining returns trainingState
func (k Keeper) GetTrainingStateValue(ctx sdk.Context) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TrainingStateKey))

	b := store.Get([]byte{0})
	if b == nil {
		return false
	}
	ts := types.TrainingState{}
	k.cdc.MustUnmarshal(b, &ts)

	return ts.Value
}

// GetTraining returns trainingState
func (k Keeper) GetTrainingStateTwinName(ctx sdk.Context) string {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TrainingStateKey))

	b := store.Get([]byte{0})
	if b == nil {
		return ""
	}
	ts := types.TrainingState{}
	k.cdc.MustUnmarshal(b, &ts)

	return ts.TwinName
}

// RemoveTrainingState removes trainingState from the store
func (k Keeper) RemoveTrainingState(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TrainingStateKey))
	store.Delete([]byte{0})
}
