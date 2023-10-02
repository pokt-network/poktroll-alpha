package types

import (
	"cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUnwhitelistApplication = "unwhitelist_application"

var _ sdk.Msg = &MsgUnwhitelistApplication{}

func NewMsgUnwhitelistApplication(portalAddress string, appAddress string) *MsgUnwhitelistApplication {
	return &MsgUnwhitelistApplication{
		PortalAddress: portalAddress,
		AppAddress:    appAddress,
	}
}

func (msg *MsgUnwhitelistApplication) Route() string {
	return RouterKey
}

func (msg *MsgUnwhitelistApplication) Type() string {
	return TypeMsgUnwhitelistApplication
}

func (msg *MsgUnwhitelistApplication) GetSigners() []sdk.AccAddress {
	portalAddress, err := sdk.AccAddressFromBech32(msg.PortalAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{portalAddress}
}

func (msg *MsgUnwhitelistApplication) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUnwhitelistApplication) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.PortalAddress); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, fmt.Errorf("invalid portal address: %w", err).Error())
	}
	if _, err := sdk.AccAddressFromBech32(msg.AppAddress); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, fmt.Errorf("invalid application address: %w", err).Error())
	}
	return nil
}
