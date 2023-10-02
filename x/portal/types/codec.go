package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgStakePortal{}, "portal/StakePortal", nil)
	cdc.RegisterConcrete(&MsgUnstakePortal{}, "portal/UnstakePortal", nil)
	cdc.RegisterConcrete(&MsgWhitelistApplication{}, "portal/WhitelistApplication", nil)
	cdc.RegisterConcrete(&MsgUnwhitelistApplication{}, "portal/UnwhitelistApplication", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgStakePortal{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUnstakePortal{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgWhitelistApplication{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUnwhitelistApplication{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
