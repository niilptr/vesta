package keeper_test

import (
	"testing"

	keepertest "vesta/testutil/keeper"
	//processortest "vesta/testutil/processor"

	"github.com/stretchr/testify/require"
)

func TestStartTraining(t *testing.T) {

	k, ctx := keepertest.NewTestKeeper(t)

	twinName := "eva00"
	creator := "testaddr0123456789"
	trainConfHash := "abcd1234efgh567"
	err := k.StartTraining(ctx, twinName, creator, trainConfHash)
	require.NoError(t, err)
}
