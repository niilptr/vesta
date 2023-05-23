package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "vesta/testutil/keeper"
	"vesta/x/twin/keeper"
	"vesta/x/twin/types"
)

func createTestTrainingState(keeper *keeper.Keeper, ctx sdk.Context, value bool) types.TrainingState {
	item := types.TrainingState{Value: value}
	keeper.SetTrainingState(ctx, item)
	return item
}

func TestTrainingStateGet(t *testing.T) {
	keeper, ctx := keepertest.TwinKeeper(t)
	item := createTestTrainingState(keeper, ctx, true)
	rst, found := keeper.GetTrainingState(ctx)
	require.True(t, found)
	require.Equal(t, item, rst)
}

func TestTrainingStateRemove(t *testing.T) {
	keeper, ctx := keepertest.TwinKeeper(t)
	createTestTrainingState(keeper, ctx, true)
	keeper.RemoveTrainingState(ctx)
	_, found := keeper.GetTrainingState(ctx)
	require.False(t, found)
}

func TestTrainingStateSetValue(t *testing.T) {
	keeper, ctx := keepertest.TwinKeeper(t)
	_ = createTestTrainingState(keeper, ctx, false)

	ts := types.TrainingState{Value: true}
	keeper.SetTrainingStateValue(ctx, true)
	rst, found := keeper.GetTrainingState(ctx)
	require.True(t, found)
	require.Equal(t, ts, rst)
}

func TestTrainingStateGetValue(t *testing.T) {
	keeper, ctx := keepertest.TwinKeeper(t)
	_ = createTestTrainingState(keeper, ctx, true)

	value := keeper.GetTrainingStateValue(ctx)
	require.True(t, value)
}
