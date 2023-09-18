package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgStakeServicer{}, "servicer/StakeServicer", nil)
	cdc.RegisterConcrete(&MsgUnstakeServicer{}, "servicer/UnstakeServicer", nil)
	cdc.RegisterConcrete(&MsgClaim{}, "servicer/Claim", nil)
	cdc.RegisterConcrete(&MsgProof{}, "servicer/Proof", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgStakeServicer{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUnstakeServicer{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgClaim{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgProof{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
