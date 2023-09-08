package modules

import (
	"poktroll/runtime/di"
	"poktroll/types"
)

type SessionManager interface {
	di.Module
	GetSession() *types.Session
	ClosedSessions() (sessions <-chan *types.Session, close func())
}

var SessionManagerToken = di.NewInjectionToken[SessionManager]("sessionManager")
