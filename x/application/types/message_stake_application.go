package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgStakeApplication = "stake_application"

var _ sdk.Msg = &MsgStakeApplication{}

func NewMsgStakeApplication(address string, stakeAmount types.Coin) *MsgStakeApplication {
	return &MsgStakeApplication{
		Address:     address,
		StakeAmount: &stakeAmount,
	}
}

func (msg *MsgStakeApplication) Route() string {
	return RouterKey
}

func (msg *MsgStakeApplication) Type() string {
	return TypeMsgStakeApplication
}

func (msg *MsgStakeApplication) GetSigners() []sdk.AccAddress {
	address, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{address}
}

func (msg *MsgStakeApplication) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// CLEANUP: Use `errors.Join` after upgrading to a newer version of go
func (msg *MsgStakeApplication) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return sdkerrors.ErrInvalidAddress
	}
	if msg.StakeAmount == nil {
		return ErrNilStakeAmount
	}
	stakeAmount, err := sdk.ParseCoinsNormalized(msg.StakeAmount.String())
	if !stakeAmount.IsValid() {
		return stakeAmount.Validate()
	}
	if err != nil {
		return err
	}
	if stakeAmount.Empty() {
		return ErrEmptyStakeAmount
	}
	return nil
}
