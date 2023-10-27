package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUnstakePortal = "unstake_portal"

var _ sdk.Msg = &MsgUnstakePortal{}

func NewMsgUnstakePortal(address string) *MsgUnstakePortal {
	return &MsgUnstakePortal{
		Address: address,
	}
}

func (msg *MsgUnstakePortal) Route() string {
	return RouterKey
}

func (msg *MsgUnstakePortal) Type() string {
	return TypeMsgUnstakePortal
}

func (msg *MsgUnstakePortal) GetSigners() []sdk.AccAddress {
	address, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{address}
}

func (msg *MsgUnstakePortal) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUnstakePortal) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address address (%s)", err)
	}
	return nil
}
