package types

import (
	"cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUnallowlistApplication = "unallowlist_application"

var _ sdk.Msg = &MsgUnallowlistApplication{}

func NewMsgUnallowlistApplication(portalAddress string, appAddress string) *MsgUnallowlistApplication {
	return &MsgUnallowlistApplication{
		PortalAddress: portalAddress,
		AppAddress:    appAddress,
	}
}

func (msg *MsgUnallowlistApplication) Route() string {
	return RouterKey
}

func (msg *MsgUnallowlistApplication) Type() string {
	return TypeMsgUnallowlistApplication
}

func (msg *MsgUnallowlistApplication) GetSigners() []sdk.AccAddress {
	portalAddress, err := sdk.AccAddressFromBech32(msg.PortalAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{portalAddress}
}

func (msg *MsgUnallowlistApplication) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUnallowlistApplication) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.PortalAddress); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, fmt.Errorf("invalid portal address: %w", err).Error())
	}
	if _, err := sdk.AccAddressFromBech32(msg.AppAddress); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, fmt.Errorf("invalid application address: %w", err).Error())
	}
	return nil
}
