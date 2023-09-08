package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSubmitClaim = "submit_claim"

var _ sdk.Msg = &MsgSubmitClaim{}

func NewMsgSubmitClaim(creator string, data string) *MsgSubmitClaim {
	return &MsgSubmitClaim{
		Creator: creator,
		Data:    data,
	}
}

func (msg *MsgSubmitClaim) Route() string {
	return RouterKey
}

func (msg *MsgSubmitClaim) Type() string {
	return TypeMsgSubmitClaim
}

func (msg *MsgSubmitClaim) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSubmitClaim) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSubmitClaim) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
