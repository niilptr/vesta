package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgConfirmTrainPhaseEnded = "confirm_train_phase_ended"

var _ sdk.Msg = &MsgConfirmTrainPhaseEnded{}

func NewMsgConfirmTrainPhaseEnded(creator string) *MsgConfirmTrainPhaseEnded {
	return &MsgConfirmTrainPhaseEnded{
		Creator: creator,
	}
}

func (msg *MsgConfirmTrainPhaseEnded) Route() string {
	return RouterKey
}

func (msg *MsgConfirmTrainPhaseEnded) Type() string {
	return TypeMsgConfirmTrainPhaseEnded
}

func (msg *MsgConfirmTrainPhaseEnded) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgConfirmTrainPhaseEnded) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgConfirmTrainPhaseEnded) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
