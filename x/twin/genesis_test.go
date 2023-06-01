package twin_test

import (
	"testing"

	keepertest "vesta/testutil/keeper"
	"vesta/testutil/nullify"
	"vesta/x/twin"
	"vesta/x/twin/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		TwinList: []types.Twin{
			{
				Name: "0",
			},
			{
				Name: "1",
			},
		},
		TrainingState: &types.TrainingState{},
	}

	k, ctx := keepertest.NewTestKeeper(t)
	twin.InitGenesis(ctx, *k, genesisState)
	got := twin.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.TwinList, got.TwinList)
	require.Equal(t, genesisState.TrainingState, got.TrainingState)
}
