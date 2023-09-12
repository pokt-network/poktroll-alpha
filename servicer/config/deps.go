package config

import (
	"poktroll/runtime/di"
)

var (
	// TECHDEBT: move this somewhere that makes more sense but doesn't set the
	// stage for an import cycle.
	PoktrollDepInjectorContextKey = "poktroll_di_injector"
	ServicerConfigToken           = di.NewInjectionToken[ServicerConfig]("servicer-config")
)
