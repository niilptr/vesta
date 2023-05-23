package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "vesta/testutil/keeper"
	"vesta/testutil/nullify"
	"vesta/x/twin/types"
)

func TestTrainingStateQuery(t *testing.T) {
	keeper, ctx := keepertest.TwinKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	item := createTestTrainingState(keeper, ctx, true)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetTrainingStateRequest
		response *types.QueryGetTrainingStateResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetTrainingStateRequest{},
			response: &types.QueryGetTrainingStateResponse{TrainingState: item},
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.TrainingState(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}
