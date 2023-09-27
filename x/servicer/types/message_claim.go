package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgClaim = "claim"

var _ sdk.Msg = &MsgClaim{}

func NewMsgClaim(servicer string, smtRootHash []byte) *MsgClaim {
	return &MsgClaim{
		ServicerAddress: servicer,
		SmstRootHash:    smtRootHash,
	}
}

func (msg *MsgClaim) Route() string {
	return RouterKey
}

func (msg *MsgClaim) Type() string {
	return TypeMsgClaim
}

func (msg *MsgClaim) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.ServicerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgClaim) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClaim) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.ServicerAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

func (msg *MsgClaim) NewClaimEvent() *EventClaim {
	return &EventClaim{
		ServicerAddress: msg.ServicerAddress,
		SmstRootHash:    msg.SmstRootHash,
	}
}

func (msg *MsgClaim) hexRootHash() string {
	return fmt.Sprintf("%x", msg.SmstRootHash)
}
