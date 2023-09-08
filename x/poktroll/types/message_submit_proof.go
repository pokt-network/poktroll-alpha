package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSubmitProof = "submit_proof"

var _ sdk.Msg = &MsgSubmitProof{}

func NewMsgSubmitProof(creator string, data string) *MsgSubmitProof {
	return &MsgSubmitProof{
		Creator: creator,
		Data:    data,
	}
}

func (msg *MsgSubmitProof) Route() string {
	return RouterKey
}

func (msg *MsgSubmitProof) Type() string {
	return TypeMsgSubmitProof
}

func (msg *MsgSubmitProof) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSubmitProof) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSubmitProof) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
