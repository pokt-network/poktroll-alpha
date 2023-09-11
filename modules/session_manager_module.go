package modules

import (
	"poktroll/runtime/di"
	"poktroll/utils"
	"poktroll/x/poktroll/types"
)

type SessionManager interface {
	di.Module
	GetSession() *types.Session
	ClosedSessions() utils.Observable[*types.Session]
}

var SessionManagerToken = di.NewInjectionToken[SessionManager]("sessionManager")
