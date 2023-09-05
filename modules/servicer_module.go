package modules

import (
	"poktroll/runtime/di"
	"poktroll/shared/crypto"
)

var (
	ServicerToken            = di.NewInjectionToken[ServicerModule]("servicer")
	PrivateKeyInjectionToken = di.NewInjectionToken[crypto.PrivateKey]("servicerPrivateKey")
)

type ServicerModule interface {
	di.Module
}
