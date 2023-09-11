package modules

import (
	"poktroll/runtime/di"
	"poktroll/utils"
	"poktroll/x/poktroll/types"
)

var MinerModuleToken = di.NewInjectionToken[MinerModule]("miner")

type MinerModule interface {
	di.Module
	MineRelays(relays utils.Observable[*types.Relay], sessions utils.Observable[*types.Session])
}
