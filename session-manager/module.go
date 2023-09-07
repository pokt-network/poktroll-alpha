package sessionmanager

import (
	"poktroll/modules"
	"poktroll/runtime/di"
	"poktroll/types"
	"sync"
)

var _ modules.SessionManager = &sessionManager{}

type sessionManager struct {
	pocketNetworkClient modules.PocketNetworkClient
	blocksPerSession    uint64
	session             *types.Session
	sessionTicker       chan *types.Session
	latestSecret        string
	started             bool
	logger              modules.Logger
	mu                  sync.RWMutex
	listeners           [](chan<- *types.Session)
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
		for session := range s.sessionTicker {
			s.mu.RLock()
			for _, ch := range s.listeners {
				ch <- session
			}
			s.mu.RUnlock()
		}
		s.mu.Lock()
		for _, ch := range s.listeners {
			close(ch)
		}
		s.listeners = [](chan<- *types.Session){}
		s.mu.Unlock()
	}()

	go func() {
		// tick sessions along as new blocks are received
		for block := range s.pocketNetworkClient.OnNewBlock() {
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
					s.sessionTicker <- s.session
				}()
				break // do we not want to continue here?
				// I think break terminates the for loop that's reading off the new sessions and that feels weird?
			}
		}
	}()

	s.started = true
	return nil
}

func (s *sessionManager) GetSession() *types.Session {
	return s.session
}

func (s *sessionManager) ClosedSessions() (sessions <-chan *types.Session, closeChan func()) {
	ch := make(chan *types.Session)
	len := len(s.listeners)
	s.listeners = append(s.listeners, ch)
	closeChan = func() {
		channelToClose := len
		// remove the channel from the list of listeners
		s.listeners = append(s.listeners[:channelToClose], s.listeners[channelToClose+1:]...)
		close(ch)
	}
	return ch, closeChan
}

// Make this module uninjectable due to single consumer return channels
func (s *sessionManager) Uninjectable() {}
