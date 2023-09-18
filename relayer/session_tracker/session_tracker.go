package sessiontracker

import (
	"context"
	"poktroll/utils"
	"poktroll/x/servicer/types"
)

var _ types.Session = &session{}

type SessionTracker struct {
	blocksPerSession uint32
	session          types.Session
	sessionTicker    utils.Observable[types.Session]
	latestSecret     []byte

	newSessions chan types.Session
	blockTicker utils.Observable[types.Block]
}

type session struct {
	sessionNumber uint64
	sessionHeight uint64
	blockHash     []byte
}

func NewSessionTracker(ctx context.Context, blocksPerSession uint32, blockTicker utils.Observable[types.Block]) *SessionTracker {
	sm := &SessionTracker{blockTicker: blockTicker, blocksPerSession: blocksPerSession}
	sm.sessionTicker, sm.newSessions = utils.NewControlledObservable[types.Session](nil)

	go sm.handleBlocks(ctx)

	return sm
}

func (sm *SessionTracker) ClosedSessions() utils.Observable[types.Session] {
	return sm.sessionTicker
}

func (sm *SessionTracker) handleBlocks(ctx context.Context) {
	// tick sessions along as new blocks are received
	ch := sm.blockTicker.Subscribe().Ch()
	for block := range ch {
		select {
		case <-ctx.Done():
			return
		default:
		}
		// discover a new session every `blocksPerSession` blocks
		if block.Height()%uint64(sm.blocksPerSession) == 0 {
			sessionNumber := block.Height() / uint64(sm.blocksPerSession)

			sm.session = &session{
				sessionNumber: sessionNumber,
				sessionHeight: sessionNumber * uint64(sm.blocksPerSession),
				blockHash:     block.Hash(),
			}

			// set the latest secret for claim and proof use
			sm.latestSecret = block.Hash()
			go func() {
				select {
				case <-ctx.Done():
					return
				default:
					sm.newSessions <- sm.session
				}
			}()
		}
	}
}

func (s *session) SessionNumber() uint64 {
	return s.sessionNumber
}

func (s *session) SessionHeight() uint64 {
	return s.sessionHeight
}

func (s *session) BlockHash() []byte {
	return s.blockHash
}
