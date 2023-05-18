package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateTwin = "create_twin"
	TypeMsgUpdateTwin = "update_twin"
	TypeMsgDeleteTwin = "delete_twin"
)

var _ sdk.Msg = &MsgCreateTwin{}

func NewMsgCreateTwin(
	creator string,
	name string,
	hash string,

) *MsgCreateTwin {
	return &MsgCreateTwin{
		Creator: creator,
		Name:    name,
		Hash:    hash,
	}
}

func (msg *MsgCreateTwin) Route() string {
	return RouterKey
}

func (msg *MsgCreateTwin) Type() string {
	return TypeMsgCreateTwin
}

func (msg *MsgCreateTwin) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateTwin) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateTwin) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateTwin{}

func NewMsgUpdateTwin(
	creator string,
	name string,
	hash string,

) *MsgUpdateTwin {
	return &MsgUpdateTwin{
		Creator: creator,
		Name:    name,
		Hash:    hash,
	}
}

func (msg *MsgUpdateTwin) Route() string {
	return RouterKey
}

func (msg *MsgUpdateTwin) Type() string {
	return TypeMsgUpdateTwin
}

func (msg *MsgUpdateTwin) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateTwin) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateTwin) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeleteTwin{}

func NewMsgDeleteTwin(
	creator string,
	name string,

) *MsgDeleteTwin {
	return &MsgDeleteTwin{
		Creator: creator,
		Name:    name,
	}
}
func (msg *MsgDeleteTwin) Route() string {
	return RouterKey
}

func (msg *MsgDeleteTwin) Type() string {
	return TypeMsgDeleteTwin
}

func (msg *MsgDeleteTwin) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteTwin) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteTwin) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
