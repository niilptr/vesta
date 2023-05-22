package keeper_test

import (
	"testing"

	keepertest "vesta/testutil/keeper"

	"github.com/stretchr/testify/require"
)

func TestGetAccessToken(t *testing.T) {
	k, ctx := keepertest.TwinKeeper(t)
	_ = ctx
	acctoken, err := k.GetAccessToken()
	require.NoError(t, err)
	require.NotEmpty(t, acctoken)
}

func TestReadTrainConfiguration(t *testing.T) {
	k, ctx := keepertest.TwinKeeper(t)
	_ = ctx
	acctoken, err := k.GetAccessToken()
	require.NoError(t, err)

	content, err := k.ReadTrainConfiguration(acctoken)
	require.NoError(t, err)
	require.NotEmpty(t, content)
}
