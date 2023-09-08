package modules

import (
	"poktroll/runtime/di"
	"poktroll/types"
	"poktroll/utils"
)

var RelayerToken = di.NewInjectionToken[RelayerModule]("relayer")

type RelayerModule interface {
	di.Module
	Relays() utils.Observable[*types.Relay]
}
