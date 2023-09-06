package modules

import (
	"poktroll/runtime/configs"
	"poktroll/runtime/di"
)

var RuntimeManagerToken = di.NewInjectionToken[RuntimeMgr]("runtimeManager")

type RuntimeMgr interface {
	di.Module

	GetConfig() *configs.Config
}
