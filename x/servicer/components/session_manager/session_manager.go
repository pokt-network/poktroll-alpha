package sessionmanager

import (
	"fmt"

	"poktroll/utils"
	"poktroll/x/servicer/types"
)

type SessionManager struct {
	blocksPerSession int64
	session          *types.Session
	sessionTicker    utils.Observable[*types.Session]
	latestSecret     []byte

	newSessions chan *types.Session
	newBlocks   chan *types.Block
}

func NewSessionManager(newBlocks chan *types.Block) *SessionManager {
	sm := &SessionManager{
		session:   &types.Session{},
		newBlocks: newBlocks,
	}
	sm.sessionTicker, sm.newSessions = utils.NewControlledObservable[*types.Session](nil)

	// TODO_THIS_COMMIT: make blocksPerSession a module param.
	sm.blocksPerSession = 2

	go sm.handleBlocks()

	return sm
}

func (sm *SessionManager) ClosedSessions() utils.Observable[*types.Session] {
	return sm.sessionTicker
}

func (sm *SessionManager) handleBlocks() {
	// tick sessions along as new blocks are received
	for block := range sm.newBlocks {
		// discover a new session every `blocksPerSession` blocks
		if block.Height%sm.blocksPerSession == 0 {
			sm.session = &types.Session{
				SessionNumber: block.Height / sm.blocksPerSession,
				SessionHeight: sm.session.SessionNumber * sm.blocksPerSession,
				BlockHash:     block.Hash,
			}

			// set the latest secret for claim and proof use
			sm.latestSecret = block.Hash
			go func() {
				fmt.Println("NEW SESSION")
				sm.newSessions <- sm.session
			}()
		}
	}
}
