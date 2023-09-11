package modules

import (
	"poktroll/runtime/di"
	"poktroll/utils"
	"poktroll/x/poktroll/types"
)

var RelayerToken = di.NewInjectionToken[RelayerModule]("relayer")

type RelayerModule interface {
	di.Module
	Relays() utils.Observable[*types.Relay]
}
