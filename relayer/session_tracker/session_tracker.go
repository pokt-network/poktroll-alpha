package sessiontracker

import (
	"context"
	"fmt"
	"poktroll/utils"
	"poktroll/x/servicer/types"
)

type SessionTracker struct {
	blocksPerSession int64
	session          *types.Session
	sessionTicker    utils.Observable[*types.Session]
	latestSecret     []byte

	newSessions chan *types.Session
	newBlocks   chan types.Block
}

func NewSessionTracker(ctx context.Context, newBlocks chan types.Block) *SessionTracker {
	sm := &SessionTracker{newBlocks: newBlocks}
	sm.sessionTicker, sm.newSessions = utils.NewControlledObservable[*types.Session](nil)

	go sm.handleBlocks(ctx)

	return sm
}

func (sm *SessionTracker) ClosedSessions() utils.Observable[*types.Session] {
	return sm.sessionTicker
}

func (sm *SessionTracker) handleBlocks(ctx context.Context) {
	// tick sessions along as new blocks are received
	for block := range sm.newBlocks {
		select {
		case <-ctx.Done():
			return
		default:
		}
		// discover a new session every `blocksPerSession` blocks
		if int64(block.Height())%sm.blocksPerSession == 0 {
			sessionNumber := int64(block.Height()) / sm.blocksPerSession

			sm.session = &types.Session{
				SessionNumber: sessionNumber,
				SessionHeight: sessionNumber * sm.blocksPerSession,
				BlockHash:     block.Hash(),
			}

			// set the latest secret for claim and proof use
			fmt.Printf("block hash: %s\n", block.Hash())
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
