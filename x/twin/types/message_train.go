package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgTrain = "train"

var _ sdk.Msg = &MsgTrain{}

func NewMsgTrain(creator string, name string, trainingHash string) *MsgTrain {
	return &MsgTrain{
		Creator:                   creator,
		Name:                      name,
		TrainingConfigurationHash: trainingHash,
	}
}

func (msg *MsgTrain) Route() string {
	return RouterKey
}

func (msg *MsgTrain) Type() string {
	return TypeMsgTrain
}

func (msg *MsgTrain) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgTrain) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgTrain) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
