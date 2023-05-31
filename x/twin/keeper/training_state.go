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

// UpdateTrainingStateValue set trainingState value in the store
func (k Keeper) MustUpdateTrainingStateValue(ctx sdk.Context, ts types.TrainingState, value bool) types.TrainingState {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TrainingStateKey))
	ts.Value = value
	b := k.cdc.MustMarshal(&ts)
	store.Set([]byte{0}, b)
	ts, found := k.GetTrainingState(ctx)
	if !found {
		panic("Training state not found after its updating.")
	}
	return ts
}

// SetTrainingStateValue set trainingState value in the store
func (k Keeper) MustUpdateTrainingStateTwinName(ctx sdk.Context, ts types.TrainingState, twinName string) types.TrainingState {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TrainingStateKey))
	ts.TwinName = twinName
	b := k.cdc.MustMarshal(&ts)
	store.Set([]byte{0}, b)
	ts, found := k.GetTrainingState(ctx)
	if !found {
		panic("Training state not found after its updating.")
	}
	return ts
}

// SetTrainingStateValue set trainingState value in the store
func (k Keeper) MustUpdateTrainingStateValidationValue(ctx sdk.Context, ts types.TrainingState, value bool) types.TrainingState {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TrainingStateKey))
	ts.ValidationState.Value = value
	b := k.cdc.MustMarshal(&ts)
	store.Set([]byte{0}, b)
	ts, found := k.GetTrainingState(ctx)
	if !found {
		panic("Training state not found after its updating.")
	}
	return ts
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
