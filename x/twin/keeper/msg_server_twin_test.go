package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	keepertest "vesta/testutil/keeper"
	"vesta/x/twin/keeper"
	"vesta/x/twin/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestTwinMsgServerCreate(t *testing.T) {
	k, ctx := keepertest.NewTestKeeper(t)
	srv := keeper.NewMsgServerImpl(*k)
	wctx := sdk.WrapSDKContext(ctx)
	creator := "A"
	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateTwin{Creator: creator,
			Name: strconv.Itoa(i),
		}
		_, err := srv.CreateTwin(wctx, expected)
		require.NoError(t, err)
		rst, found := k.GetTwin(ctx,
			expected.Name,
		)
		require.True(t, found)
		require.Equal(t, expected.Creator, rst.Creator)
	}
}

func TestTwinMsgServerUpdate(t *testing.T) {
	creator := "A"

	for _, tc := range []struct {
		desc    string
		request *types.MsgUpdateTwin
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgUpdateTwin{Creator: creator,
				Name: strconv.Itoa(0),
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgUpdateTwin{Creator: "B",
				Name: strconv.Itoa(0),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgUpdateTwin{Creator: creator,
				Name: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.NewTestKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			wctx := sdk.WrapSDKContext(ctx)
			expected := &types.MsgCreateTwin{Creator: creator,
				Name: strconv.Itoa(0),
			}
			_, err := srv.CreateTwin(wctx, expected)
			require.NoError(t, err)

			_, err = srv.UpdateTwin(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				rst, found := k.GetTwin(ctx,
					expected.Name,
				)
				require.True(t, found)
				require.Equal(t, expected.Creator, rst.Creator)
			}
		})
	}
}

func TestTwinMsgServerDelete(t *testing.T) {
	creator := "A"

	for _, tc := range []struct {
		desc    string
		request *types.MsgDeleteTwin
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgDeleteTwin{Creator: creator,
				Name: strconv.Itoa(0),
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgDeleteTwin{Creator: "B",
				Name: strconv.Itoa(0),
			},
			err: sdkerrors.ErrUnauthorized,
		},
		{
			desc: "KeyNotFound",
			request: &types.MsgDeleteTwin{Creator: creator,
				Name: strconv.Itoa(100000),
			},
			err: sdkerrors.ErrKeyNotFound,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.NewTestKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			wctx := sdk.WrapSDKContext(ctx)

			_, err := srv.CreateTwin(wctx, &types.MsgCreateTwin{Creator: creator,
				Name: strconv.Itoa(0),
			})
			require.NoError(t, err)
			_, err = srv.DeleteTwin(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				_, found := k.GetTwin(ctx,
					tc.request.Name,
				)
				require.False(t, found)
			}
		})
	}
}
