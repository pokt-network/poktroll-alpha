package types

import (
	"cosmossdk.io/errors"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocdc "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgDelegateToPortal = "delegate_to_portal"

var _ sdk.Msg = &MsgDelegateToPortal{}

func NewMsgDelegateToPortal(address string, portalPubKey string) *MsgDelegateToPortal {
	reg := codectypes.NewInterfaceRegistry()
	cryptocdc.RegisterInterfaces(reg)
	cdc := codec.NewProtoCodec(reg)
	anyPk := new(codectypes.Any)
	if err := cdc.UnmarshalJSON(json.RawMessage(portalPubKey), anyPk); err != nil {
		panic(fmt.Errorf("unable pack portal public key into Any: %w", err))
	}
	return &MsgDelegateToPortal{
		Address:      address,
		PortalPubKey: anyPk,
	}
}

func (msg *MsgDelegateToPortal) Route() string {
	return RouterKey
}

func (msg *MsgDelegateToPortal) Type() string {
	return TypeMsgDelegateToPortal
}

func (msg *MsgDelegateToPortal) GetSigners() []sdk.AccAddress {
	address, err := sdk.AccAddressFromBech32(msg.Address)
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
	_, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address address (%s)", err)
	}
	reg := codectypes.NewInterfaceRegistry()
	cryptocdc.RegisterInterfaces(reg)
	cdc := codec.NewProtoCodec(reg)
	var pubI cryptotypes.PubKey
	if err := cdc.UnpackAny(msg.PortalPubKey, &pubI); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidPubKey, "portal public key is not cryptotypes.PubKey")
	}
	return nil
}
