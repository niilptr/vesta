package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "vesta/testutil/keeper"
	"vesta/testutil/nullify"
	"vesta/x/twin/keeper"
	"vesta/x/twin/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNTwin(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Twin {
	items := make([]types.Twin, n)
	for i := range items {
		items[i].Name = strconv.Itoa(i)

		keeper.SetTwin(ctx, items[i])
	}
	return items
}

func TestTwinGet(t *testing.T) {
	keeper, ctx := keepertest.TwinKeeper(t)
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
	keeper, ctx := keepertest.TwinKeeper(t)
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
	keeper, ctx := keepertest.TwinKeeper(t)
	items := createNTwin(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllTwin(ctx)),
	)
}
