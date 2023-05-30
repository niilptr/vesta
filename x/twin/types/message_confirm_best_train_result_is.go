package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgConfirmBestTrainResultIs = "confirm_best_train_result_is"

var _ sdk.Msg = &MsgConfirmBestTrainResultIs{}

func NewMsgConfirmBestTrainResultIs(creator string, hash string) *MsgConfirmBestTrainResultIs {
	return &MsgConfirmBestTrainResultIs{
		Creator: creator,
		Hash:    hash,
	}
}

func (msg *MsgConfirmBestTrainResultIs) Route() string {
	return RouterKey
}

func (msg *MsgConfirmBestTrainResultIs) Type() string {
	return TypeMsgConfirmBestTrainResultIs
}

func (msg *MsgConfirmBestTrainResultIs) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgConfirmBestTrainResultIs) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgConfirmBestTrainResultIs) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
