package sessionmanager

import (
	"context"
	"crypto/sha256"
	"log"
	"path/filepath"

	"github.com/pokt-network/smt"

	"poktroll/utils"
	"poktroll/x/servicer/types"
	sessionTypes "poktroll/x/session/types"
)

type SessionManager struct {
	// map[sessionEndHeight]map[sessionId]SessionWithTree
	// sessionEndHeight groups sessions that end at the same height
	// supports the case where ALL sessions end at the same height
	// supports different sessions ending (e.g. per service)
	sessions         map[uint64]map[string]SessionWithTree
	sessionsNotifier chan map[string]SessionWithTree // channel emitting map[sessionId]SessionWithTree
	sessionsNotifee  utils.Observable[map[string]SessionWithTree]
	client           types.ServicerClient
	storeDirectory   string // directory that will contain session tree stores
}

func NewSessionManager(ctx context.Context, storeDirectory string, client types.ServicerClient) *SessionManager {
	sessions := make(map[uint64]map[string]SessionWithTree)
	sm := &SessionManager{client: client, storeDirectory: storeDirectory, sessions: sessions}
	sm.sessionsNotifee, sm.sessionsNotifier = utils.NewControlledObservable[map[string]SessionWithTree](nil)

	go sm.handleBlocks(ctx)

	return sm
}

// emits all sessions that have ended
func (sm *SessionManager) Sessions() utils.Observable[map[string]SessionWithTree] {
	return sm.sessionsNotifee
}

// returns the tree for a given session, creating it if it doesn't exist
func (sm *SessionManager) EnsureSessionTree(sessionInfo *sessionTypes.Session) *smt.SMST {
	// get session end so we can group sessions that end at the same height
	// make sure we do not off by one
	sessionEnd := sessionInfo.SessionBlockStartHeight + sessionInfo.NumBlocksPerSession
	sessionId := sessionInfo.SessionId

	// make sure to have a container for sessions that end at this height
	if _, ok := sm.sessions[sessionEnd]; !ok {
		sm.sessions[sessionEnd] = make(map[string]SessionWithTree)
	}

	// create session tree if it doesn't exist (first relay for this session)
	// we need to get its store so we can close it later since we can't access it from the tree
	if _, ok := sm.sessions[sessionEnd][sessionId]; !ok {
		storePath := filepath.Join(sm.storeDirectory, sessionId)
		tree, store, err := sm.createTreeForSession(storePath)
		if err != nil {
			log.Println("error creating tree for session", err)
			return nil
		}

		removeFromSessionsMap := func() {
			delete(sm.sessions[sessionEnd], sessionId)

			// delete sessionEnd map if it's empty
			itemsCount := 0
			for _, _ = range sm.sessions[sessionEnd] {
				itemsCount++
				break
			}

			if itemsCount == 0 {
				delete(sm.sessions, sessionEnd)
			}
		}

		sm.sessions[sessionEnd][sessionId] = &sessionWithTree{
			sessionInfo:           sessionInfo,
			tree:                  tree,
			treeStore:             store,
			storePath:             storePath,
			removeFromSessionsMap: removeFromSessionsMap,
		}
	}

	return sm.sessions[sessionEnd][sessionId].SessionTree()
}

func (sm *SessionManager) handleBlocks(ctx context.Context) {
	// tick sessions along as new blocks are received
	ch := sm.client.Blocks().Subscribe().Ch()
	for block := range ch {
		select {
		case <-ctx.Done():
			return
		default:
		}
		height := block.Height()
		// if some sessions end by this block, process them
		if sessions, ok := sm.sessions[height]; !ok {
			sm.sessionsNotifier <- sessions
		}
	}
}

func (sm *SessionManager) createTreeForSession(storePath string) (*smt.SMST, smt.KVStore, error) {
	treeStore, err := smt.NewKVStore(storePath)
	if err != nil {
		return nil, nil, err
	}

	tree := smt.NewSparseMerkleSumTree(treeStore, sha256.New())
	return tree, treeStore, nil
}
