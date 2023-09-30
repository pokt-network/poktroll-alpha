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

const TypeMsgUndelegateFromPortal = "undelegate_from_portal"

var _ sdk.Msg = &MsgUndelegateFromPortal{}

func NewMsgUndelegateFromPortal(address string, portalPubKey string) *MsgUndelegateFromPortal {
	reg := codectypes.NewInterfaceRegistry()
	cryptocdc.RegisterInterfaces(reg)
	cdc := codec.NewProtoCodec(reg)
	anyPk := new(codectypes.Any)
	if err := cdc.UnmarshalJSON(json.RawMessage(portalPubKey), anyPk); err != nil {
		panic(fmt.Errorf("unable pack portal public key into Any: %w", err))
	}
	return &MsgUndelegateFromPortal{
		Address:      address,
		PortalPubKey: anyPk,
	}
}

func (msg *MsgUndelegateFromPortal) Route() string {
	return RouterKey
}

func (msg *MsgUndelegateFromPortal) Type() string {
	return TypeMsgUndelegateFromPortal
}

func (msg *MsgUndelegateFromPortal) GetSigners() []sdk.AccAddress {
	address, err := sdk.AccAddressFromBech32(msg.Address)
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
