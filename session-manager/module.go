package sessionmanager

import (
	"poktroll/modules"
	"poktroll/runtime/di"
	"poktroll/types"
)

type sessionManager struct {
	pocketNetworkClient modules.PocketNetworkClient
	blocksPerSession    uint64
	session             *types.Session
	sessionTicker       chan *types.Session
	latestSecret        string
	started             bool
	logger              modules.Logger
}

func NewSessionManager() modules.SessionManager {
	return &sessionManager{session: &types.Session{}, sessionTicker: make(chan *types.Session, 1)}
}

func (s *sessionManager) Resolve(injector *di.Injector, path *[]string) {
	s.pocketNetworkClient = di.Resolve(modules.PocketNetworkClientToken, injector, path)
	s.blocksPerSession = di.Resolve(modules.RuntimeManagerToken, injector, path).GetConfig().BlocksPerSession
	s.logger = *di.Resolve(modules.LoggerModuleToken, injector, path).
		CreateLoggerForModule(modules.ServicerToken.Id())
}

func (s *sessionManager) CascadeStart() error {
	s.pocketNetworkClient.CascadeStart()
	return s.Start()
}

func (s *sessionManager) Start() error {
	if s.started {
		return nil
	}

	go func() {
		for block := range s.pocketNetworkClient.OnNewBlock() {
			if block.Height%s.blocksPerSession == 0 {
				s.session = &types.Session{
					SessionNumber:      block.Height / s.blocksPerSession,
					SessionBlockHeight: s.session.SessionNumber * s.blocksPerSession,
					BlockHeight:        block.Height,
					BlockHash:          block.Hash,
				}

				s.latestSecret = block.Hash
				go func() {
					s.sessionTicker <- s.session
				}()
				break
			}
		}
	}()

	s.started = true
	return nil
}

func (s *sessionManager) GetSession() *types.Session {
	return s.session
}

func (s *sessionManager) OnSessionEnd() <-chan *types.Session {
	return s.sessionTicker
}
