package modules

import (
	"poktroll/runtime/di"
	"poktroll/types"
)

var RelayerToken = di.NewInjectionToken[RelayerModule]("relayer")

type RelayerModule interface {
	di.Module
	Relays() <-chan *types.Relay
}
