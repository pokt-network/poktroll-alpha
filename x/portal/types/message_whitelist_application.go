package types

import (
	"cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgWhitelistApplication = "whitelist_application"

var _ sdk.Msg = &MsgWhitelistApplication{}

func NewMsgWhitelistApplication(portalAddress, appAddress string) *MsgWhitelistApplication {
	return &MsgWhitelistApplication{
		PortalAddress: portalAddress,
		AppAddress:    appAddress,
	}
}

func (msg *MsgWhitelistApplication) Route() string {
	return RouterKey
}

func (msg *MsgWhitelistApplication) Type() string {
	return TypeMsgWhitelistApplication
}

func (msg *MsgWhitelistApplication) GetSigners() []sdk.AccAddress {
	address, err := sdk.AccAddressFromBech32(msg.PortalAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{address}
}

func (msg *MsgWhitelistApplication) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgWhitelistApplication) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.PortalAddress); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, fmt.Errorf("invalid portal address: %w", err).Error())
	}
	if _, err := sdk.AccAddressFromBech32(msg.AppAddress); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, fmt.Errorf("invalid application address: %w", err).Error())
	}
	return nil
}
