package modules

import (
	"poktroll/runtime/di"
	"poktroll/types"
)

type SessionManager interface {
	di.Module
	GetSession() *types.Session
	OnSessionEnd() <-chan *types.Session
}

var SessionManagerToken = di.NewInjectionToken[SessionManager]("sessionManager")
