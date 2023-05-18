package twin_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "vesta/testutil/keeper"
	"vesta/testutil/nullify"
	"vesta/x/twin"
	"vesta/x/twin/types"
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
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.TwinKeeper(t)
	twin.InitGenesis(ctx, *k, genesisState)
	got := twin.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.TwinList, got.TwinList)
	// this line is used by starport scaffolding # genesis/test/assert
}
