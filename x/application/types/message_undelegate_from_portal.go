package types

import (
	"cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUndelegateFromPortal = "undelegate_from_portal"

var _ sdk.Msg = &MsgUndelegateFromPortal{}

func NewMsgUndelegateFromPortal(address, portalAddress string) *MsgUndelegateFromPortal {
	return &MsgUndelegateFromPortal{
		AppAddress:    address,
		PortalAddress: portalAddress,
	}
}

func (msg *MsgUndelegateFromPortal) Route() string {
	return RouterKey
}

func (msg *MsgUndelegateFromPortal) Type() string {
	return TypeMsgUndelegateFromPortal
}

func (msg *MsgUndelegateFromPortal) GetSigners() []sdk.AccAddress {
	address, err := sdk.AccAddressFromBech32(msg.AppAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{address}
}

func (msg *MsgUndelegateFromPortal) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUndelegateFromPortal) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.AppAddress); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, fmt.Errorf("invalid app address (%s): %w", msg.AppAddress, err).Error())
	}
	if _, err := sdk.AccAddressFromBech32(msg.PortalAddress); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, fmt.Errorf("invalid portal address (%s): %w", msg.PortalAddress, err).Error())
	}
	return nil
}
