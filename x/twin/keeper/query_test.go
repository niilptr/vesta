package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "vesta/testutil/keeper"
	"vesta/testutil/nullify"
	"vesta/x/twin/types"
)

func TestParamsQuery(t *testing.T) {
	keeper, ctx := keepertest.NewTestKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	params := types.DefaultParams()
	keeper.SetParams(ctx, params)

	response, err := keeper.Params(wctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}

func TestTwinQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.NewTestKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNTwin(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetTwinRequest
		response *types.QueryGetTwinResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetTwinRequest{
				Name: msgs[0].Name,
			},
			response: &types.QueryGetTwinResponse{Twin: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetTwinRequest{
				Name: msgs[1].Name,
			},
			response: &types.QueryGetTwinResponse{Twin: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetTwinRequest{
				Name: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Twin(wctx, tc.request)
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

func TestTwinQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.NewTestKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNTwin(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllTwinRequest {
		return &types.QueryAllTwinRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.TwinAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Twin), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Twin),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.TwinAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Twin), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Twin),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.TwinAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Twin),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.TwinAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}

func TestTrainingStateQuery(t *testing.T) {
	keeper, ctx := keepertest.NewTestKeeper(t)
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
