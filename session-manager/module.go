package sessionmanager

import (
	"poktroll/modules"
	"poktroll/runtime/di"
	"poktroll/types"
	"poktroll/utils"
)

type sessionManager struct {
	pocketNetworkClient modules.PocketNetworkClient
	blocksPerSession    uint64
	session             *types.Session
	sessionTicker       utils.Observable[*types.Session]
	latestSecret        string
	started             bool
	logger              modules.Logger
}

func NewSessionManager() modules.SessionManager {
	return &sessionManager{session: &types.Session{}}
}

func (s *sessionManager) Hydrate(injector *di.Injector, path *[]string) {
	s.pocketNetworkClient = di.Hydrate(modules.PocketNetworkClientToken, injector, path)
	s.blocksPerSession = di.Hydrate(modules.RuntimeManagerToken, injector, path).GetConfig().BlocksPerSession
	s.logger = *di.Hydrate(modules.LoggerModuleToken, injector, path).
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

	observable, ticker := utils.NewControlledObservable[*types.Session]()
	s.sessionTicker = observable

	go func() {
		// tick sessions along as new blocks are received
		for block := range s.pocketNetworkClient.NewBlocks() {
			// discover a new session every `blocksPerSession` blocks
			if block.Height%s.blocksPerSession == 0 {
				s.session = &types.Session{
					SessionNumber:      block.Height / s.blocksPerSession,
					SessionBlockHeight: s.session.SessionNumber * s.blocksPerSession,
					BlockHeight:        block.Height,
					BlockHash:          block.Hash,
				}

				// set the latest secret for claim and proof use
				s.latestSecret = block.Hash
				go func() {
					ticker <- s.session
				}()
			}
		}
	}()

	s.started = true
	return nil
}

func (s *sessionManager) GetSession() *types.Session {
	return s.session
}

func (s *sessionManager) ClosedSessions() utils.Observable[*types.Session] {
	return s.sessionTicker
}
