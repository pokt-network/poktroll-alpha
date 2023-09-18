package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgProof = "proof"

var _ sdk.Msg = &MsgProof{}

func NewMsgProof(
	creator string,
	root,
	path,
	valueHash []byte,
	sum int32,
	proofBz []byte,
) *MsgProof {
	return &MsgProof{
		Creator:   creator,
		Root:      root,
		Path:      path,
		ValueHash: valueHash,
		Sum:       sum,
		ProofBz:   proofBz,
	}
}

func (msg *MsgProof) Route() string {
	return RouterKey
}

func (msg *MsgProof) Type() string {
	return TypeMsgProof
}

func (msg *MsgProof) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgProof) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgProof) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
