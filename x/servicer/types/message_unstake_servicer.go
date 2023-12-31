package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUnstakeServicer = "unstake_servicer"

var _ sdk.Msg = &MsgUnstakeServicer{}

func NewMsgUnstakeServicer(address string) *MsgUnstakeServicer {
	return &MsgUnstakeServicer{
		Address: address,
	}
}

func (msg *MsgUnstakeServicer) Route() string {
	return RouterKey
}

func (msg *MsgUnstakeServicer) Type() string {
	return TypeMsgUnstakeServicer
}

func (msg *MsgUnstakeServicer) GetSigners() []sdk.AccAddress {
	address, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{address}
}

func (msg *MsgUnstakeServicer) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUnstakeServicer) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address address (%s)", err)
	}
	return nil
}
