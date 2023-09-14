package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgStakeServicer = "stake_servicer"

var _ sdk.Msg = &MsgStakeServicer{}

func NewMsgStakeServicer(address string, stakeAmount types.Coin) *MsgStakeServicer {
	return &MsgStakeServicer{
		Address:     address,
		StakeAmount: &stakeAmount,
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

// CLEANUP: Use `errors.Join` after upgrading to a newer version of go
func (msg *MsgStakeServicer) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return sdkerrors.ErrInvalidAddress
	}
	if msg.StakeAmount == nil {
		return ErrNilStakeAmount
	}
	stakeAmount, err := sdk.ParseCoinNormalized(msg.StakeAmount.String())
	if !stakeAmount.IsValid() {
		return stakeAmount.Validate()
	}
	if err != nil {
		return err
	}
	return nil
}
