package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgStakeServicer = "stake_servicer"

var _ sdk.Msg = &MsgStakeServicer{}

func NewMsgStakeServicer(address string) *MsgStakeServicer {
	return &MsgStakeServicer{
		Address: address,
	}
}

func (msg *MsgStakeServicer) Route() string {
	return RouterKey
}

func (msg *MsgStakeServicer) Type() string {
	return TypeMsgStakeServicer
}

func (msg *MsgStakeServicer) GetSigners() []sdk.AccAddress {
	address, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{address}
}

func (msg *MsgStakeServicer) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgStakeServicer) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address address (%s)", err)
	}
	return nil
}
