package modules

import (
	"poktroll/runtime/di"
	"poktroll/types"
)

type SessionManager interface {
	di.Module
	di.Uninjectable
	GetSession() *types.Session
	ClosedSessions() (sessions <-chan *types.Session, close func())
}

var SessionManagerToken = di.NewInjectionToken[SessionManager]("sessionManager")
