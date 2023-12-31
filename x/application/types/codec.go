package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgStakeApplication{}, "application/StakeApplication", nil)
	cdc.RegisterConcrete(&MsgUnstakeApplication{}, "application/UnstakeApplication", nil)
	cdc.RegisterConcrete(&MsgDelegateToPortal{}, "application/DelegateToPortal", nil)
	cdc.RegisterConcrete(&MsgUndelegateFromPortal{}, "application/UndelegateFromPortal", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgStakeApplication{},
		&MsgUnstakeApplication{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgDelegateToPortal{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUndelegateFromPortal{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
