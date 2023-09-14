package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgGetSession = "get_session"

var _ sdk.Msg = &MsgGetSession{}

func NewMsgGetSession(address string) *MsgGetSession {
	return &MsgGetSession{
		Address: address,
	}
}

func (msg *MsgGetSession) Route() string {
	return RouterKey
}

func (msg *MsgGetSession) Type() string {
	return TypeMsgGetSession
}

func (msg *MsgGetSession) GetSigners() []sdk.AccAddress {
	address, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{address}
}

func (msg *MsgGetSession) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgGetSession) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address address (%s)", err)
	}
	return nil
}
