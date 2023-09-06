package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// ServicerPrefix defines the key to store servicer-related data
const (
	ServicerPrefix = "servicer-"
)

type Servicer struct {
	Address     sdk.ValAddress `json:"val_address"`
	StakeAmount sdk.Coin       `json:"stake_amount"`
}

// Marshal implements codec.ProtoMarshaler.
func (*Servicer) Marshal() ([]byte, error) {
	panic("unimplemented")
}

// MarshalTo implements codec.ProtoMarshaler.
func (*Servicer) MarshalTo(data []byte) (n int, err error) {
	panic("unimplemented")
}

// MarshalToSizedBuffer implements codec.ProtoMarshaler.
func (*Servicer) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	panic("unimplemented")
}

// ProtoMessage implements codec.ProtoMarshaler.
func (*Servicer) ProtoMessage() {
	panic("unimplemented")
}

// Reset implements codec.ProtoMarshaler.
func (*Servicer) Reset() {
	panic("unimplemented")
}

// Size implements codec.ProtoMarshaler.
func (*Servicer) Size() int {
	panic("unimplemented")
}

// String implements codec.ProtoMarshaler.
func (*Servicer) String() string {
	panic("unimplemented")
}

// Unmarshal implements codec.ProtoMarshaler.
func (*Servicer) Unmarshal(data []byte) error {
	panic("unimplemented")
}

func NewServicer(addr sdk.ValAddress, amount sdk.Coin) Servicer {
	return Servicer{
		Address:     addr,
		StakeAmount: amount,
	}
}
