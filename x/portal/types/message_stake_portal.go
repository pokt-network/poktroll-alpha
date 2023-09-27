package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgStakePortal = "stake_portal"

var _ sdk.Msg = &MsgStakePortal{}

func NewMsgStakePortal(address string, stakeAmount types.Coin) *MsgStakePortal {
	return &MsgStakePortal{
		Address:     address,
		StakeAmount: &stakeAmount,
	}
}

func (msg *MsgStakePortal) Route() string {
	return RouterKey
}

func (msg *MsgStakePortal) Type() string {
	return TypeMsgStakePortal
}

func (msg *MsgStakePortal) GetSigners() []sdk.AccAddress {
	address, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{address}
}

func (msg *MsgStakePortal) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// CLEANUP: Use `errors.Join` after upgrading to a newer version of go
// TODO_TEST: So much validation that has to go into this to make it work better
func (msg *MsgStakePortal) ValidateBasic() error {
	fmt.Println("MsgStakePortal.ValidateBasic()", msg)
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
	if len(msg.Services) == 0 {
		return ErrNoServicesToStake
	}

	return nil
}
