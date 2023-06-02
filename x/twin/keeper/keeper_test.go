package keeper_test

import (
	"strconv"
	"testing"

	keepertest "vesta/testutil/keeper"
	//processortest "vesta/testutil/processor"

	"vesta/testutil/nullify"
	"vesta/x/twin/keeper"
	"vesta/x/twin/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"
)

// ====================================================================================
// Params
// ====================================================================================

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.NewTestKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}

// ====================================================================================
// Twin
// ====================================================================================

func createNTwin(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Twin {
	items := make([]types.Twin, n)
	for i := range items {
		items[i].Name = strconv.Itoa(i)

		keeper.SetTwin(ctx, items[i])
	}
	return items
}

func TestTwinGet(t *testing.T) {
	keeper, ctx := keepertest.NewTestKeeper(t)
	items := createNTwin(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetTwin(ctx,
			item.Name,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestTwinRemove(t *testing.T) {
	keeper, ctx := keepertest.NewTestKeeper(t)
	items := createNTwin(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveTwin(ctx,
			item.Name,
		)
		_, found := keeper.GetTwin(ctx,
			item.Name,
		)
		require.False(t, found)
	}
}

func TestTwinGetAll(t *testing.T) {
	keeper, ctx := keepertest.NewTestKeeper(t)
	items := createNTwin(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllTwin(ctx)),
	)
}

// ====================================================================================
// Training state
// ====================================================================================

func createTestTrainingState(keeper *keeper.Keeper, ctx sdk.Context, value bool) types.TrainingState {
	item := types.TrainingState{Value: value}
	keeper.SetTrainingState(ctx, item)
	return item
}

func TestTrainingStateGet(t *testing.T) {
	keeper, ctx := keepertest.NewTestKeeper(t)
	item := createTestTrainingState(keeper, ctx, true)
	rst, found := keeper.GetTrainingState(ctx)
	require.True(t, found)
	require.Equal(t, item, rst)
}

func TestTrainingStateRemove(t *testing.T) {
	keeper, ctx := keepertest.NewTestKeeper(t)
	createTestTrainingState(keeper, ctx, true)
	keeper.RemoveTrainingState(ctx)
	_, found := keeper.GetTrainingState(ctx)
	require.False(t, found)
}

func TestTrainingStateSetValue(t *testing.T) {
	keeper, ctx := keepertest.NewTestKeeper(t)
	ts := createTestTrainingState(keeper, ctx, false)
	keeper.MustUpdateTrainingStateValue(ctx, ts, true)
	rst, found := keeper.GetTrainingState(ctx)
	require.True(t, found)
	require.Equal(t, ts, rst)
}

func TestTrainingStateGetValue(t *testing.T) {
	keeper, ctx := keepertest.NewTestKeeper(t)
	_ = createTestTrainingState(keeper, ctx, true)
	value := keeper.GetTrainingStateValue(ctx)
	require.True(t, value)
}

// ====================================================================================
// Train
// ====================================================================================

func TestStartTraining(t *testing.T) {

	k, ctx := keepertest.NewTestKeeper(t)

	twinName := "eva00"
	creator := "testaddr0123456789"
	trainConfHash := "abcd1234efgh567"
	err := k.StartTraining(ctx, twinName, creator, trainConfHash)
	require.NoError(t, err)
}

// ====================================================================================
// Confirm train phase ended
// ====================================================================================
func TestAddTrainingPhaseEndedConfirmation(t *testing.T) {

	k, ctx := keepertest.NewTestKeeper(t)

	ts := types.NewEmptyTrainingState()
	ts.Value = true

	k.SetTrainingState(ctx, ts)

	err := k.AddTrainingPhaseEndedConfirmation(ctx, "vesta1testaddress01234")
	require.NoError(t, err)

	ts, found := k.GetTrainingState(ctx)
	require.True(t, found)
	require.NotNil(t, ts.TrainingPhaseEndedConfirmations)
	require.Equal(t, 1, len(ts.TrainingPhaseEndedConfirmations))

	err = k.AddTrainingPhaseEndedConfirmation(ctx, "vesta1testaddress5678")
	require.NoError(t, err)

	ts, found = k.GetTrainingState(ctx)
	require.True(t, found)
	require.NotNil(t, ts.TrainingPhaseEndedConfirmations)
	require.Equal(t, 2, len(ts.TrainingPhaseEndedConfirmations))

	err = k.AddTrainingPhaseEndedConfirmation(ctx, "vesta1testaddress5678")
	require.NoError(t, err)

	ts, found = k.GetTrainingState(ctx)
	require.True(t, found)
	require.NotNil(t, ts.TrainingPhaseEndedConfirmations)
	require.Equal(t, 2, len(ts.TrainingPhaseEndedConfirmations))

}
