package types

import (
	"testing"

	"vesta/testutil/sample"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateTwin_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreateTwin
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreateTwin{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCreateTwin{
				Creator: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgUpdateTwin_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateTwin
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateTwin{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateTwin{
				Creator: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgDeleteTwin_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgDeleteTwin
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgDeleteTwin{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgDeleteTwin{
				Creator: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

/////////////////////////////

func TestMsgTrain_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgTrain
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgTrain{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgTrain{
				Creator: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

//////////////////////////

func TestMsgConfirmTrainPhaseEnded_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgConfirmTrainPhaseEnded
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgConfirmTrainPhaseEnded{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgConfirmTrainPhaseEnded{
				Creator: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

//////////////////////////////////

func TestMsgConfirmBestTrainResultIs_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgConfirmBestTrainResultIs
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgConfirmBestTrainResultIs{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgConfirmBestTrainResultIs{
				Creator: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
