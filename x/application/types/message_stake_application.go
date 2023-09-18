package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgStakeApplication = "stake_application"

var _ sdk.Msg = &MsgStakeApplication{}

func NewMsgStakeApplication(
	address string,
	stakeAmount types.Coin,
	serviceIds []string,
) *MsgStakeApplication {
	return &MsgStakeApplication{
		Address:     address,
		StakeAmount: &stakeAmount,
		ServiceIds:  serviceIds,
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
// TODO_TEST: So much validation that has to go into this to make it work better
func (msg *MsgStakeApplication) ValidateBasic() error {
	// Validate the address
	_, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return sdkerrors.ErrInvalidAddress
	}

	// Validate the stake amount
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

	// Validate the services
	if len(msg.ServiceIds) == 0 {
		return ErrNoServicesToStake
	}

	return nil
}
