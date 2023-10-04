package types

import (
	"cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgAllowlistApplication = "allowlist_application"

var _ sdk.Msg = &MsgAllowlistApplication{}

func NewMsgAllowlistApplication(portalAddress, appAddress string) *MsgAllowlistApplication {
	return &MsgAllowlistApplication{
		PortalAddress: portalAddress,
		AppAddress:    appAddress,
	}
}

func (msg *MsgAllowlistApplication) Route() string {
	return RouterKey
}

func (msg *MsgAllowlistApplication) Type() string {
	return TypeMsgAllowlistApplication
}

func (msg *MsgAllowlistApplication) GetSigners() []sdk.AccAddress {
	address, err := sdk.AccAddressFromBech32(msg.PortalAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{address}
}

func (msg *MsgAllowlistApplication) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAllowlistApplication) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.PortalAddress); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, fmt.Errorf("invalid portal address: %w", err).Error())
	}
	if _, err := sdk.AccAddressFromBech32(msg.AppAddress); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, fmt.Errorf("invalid application address: %w", err).Error())
	}
	return nil
}
