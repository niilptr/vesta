package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	testkeeper "vesta/testutil/keeper"
	"vesta/x/twin/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.TwinKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
