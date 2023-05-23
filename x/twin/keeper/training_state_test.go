package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "vesta/testutil/keeper"
	"vesta/testutil/nullify"
	"vesta/x/twin/keeper"
	"vesta/x/twin/types"
)

func createTestTrainingState(keeper *keeper.Keeper, ctx sdk.Context) types.TrainingState {
	item := types.TrainingState{}
	keeper.SetTrainingState(ctx, item)
	return item
}

func TestTrainingStateGet(t *testing.T) {
	keeper, ctx := keepertest.TwinKeeper(t)
	item := createTestTrainingState(keeper, ctx)
	rst, found := keeper.GetTrainingState(ctx)
	require.True(t, found)
	require.Equal(t,
		nullify.Fill(&item),
		nullify.Fill(&rst),
	)
}

func TestTrainingStateRemove(t *testing.T) {
	keeper, ctx := keepertest.TwinKeeper(t)
	createTestTrainingState(keeper, ctx)
	keeper.RemoveTrainingState(ctx)
	_, found := keeper.GetTrainingState(ctx)
	require.False(t, found)
}
