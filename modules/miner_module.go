package modules

import (
	"poktroll/runtime/di"
	"poktroll/types"
	"poktroll/utils"
)

var MinerModuleToken = di.NewInjectionToken[MinerModule]("miner")

type MinerModule interface {
	di.Module
	MineRelays(relays utils.Observable[*types.Relay], sessions utils.Observable[*types.Session])
}
