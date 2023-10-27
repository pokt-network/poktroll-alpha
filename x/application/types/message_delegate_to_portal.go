package types

import (
	"cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgDelegateToPortal = "delegate_to_portal"

var _ sdk.Msg = &MsgDelegateToPortal{}

func NewMsgDelegateToPortal(address, portalAddress string) *MsgDelegateToPortal {
	return &MsgDelegateToPortal{
		AppAddress:    address,
		PortalAddress: portalAddress,
	}
}

func (msg *MsgDelegateToPortal) Route() string {
	return RouterKey
}

func (msg *MsgDelegateToPortal) Type() string {
	return TypeMsgDelegateToPortal
}

func (msg *MsgDelegateToPortal) GetSigners() []sdk.AccAddress {
	address, err := sdk.AccAddressFromBech32(msg.AppAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{address}
}

func (msg *MsgDelegateToPortal) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDelegateToPortal) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.AppAddress); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, fmt.Errorf("invalid app address (%s): %w", msg.AppAddress, err).Error())
	}
	if _, err := sdk.AccAddressFromBech32(msg.PortalAddress); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, fmt.Errorf("invalid portal address (%s): %w", msg.PortalAddress, err).Error())
	}
	return nil
}
