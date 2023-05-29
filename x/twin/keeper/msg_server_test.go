package keeper_test

import (
	"context"
	"testing"

	keepertest "vesta/testutil/keeper"
	"vesta/x/twin/keeper"
	"vesta/x/twin/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.NewTestKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
