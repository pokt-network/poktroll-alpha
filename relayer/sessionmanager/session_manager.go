package sessionmanager

import (
	"context"
	"crypto/sha256"
	"log"
	"path/filepath"

	"github.com/pokt-network/smt"

	"poktroll/relayer/client"
	"poktroll/utils"
	sessionTypes "poktroll/x/session/types"
)

type sessionTreeMap map[string]SessionWithTree

type SessionManager struct {
	// map[sessionEndHeight]map[sessionId]SessionWithTree
	// sessionEndHeight groups sessions that end at the same height
	// supports the case where ALL sessions end at the same height
	// supports different sessions ending (e.g. per service)
	sessions         map[uint64]sessionTreeMap
	sessionsNotifier chan sessionTreeMap // channel emitting map[sessionId]SessionWithTree
	sessionsNotifee  utils.Observable[sessionTreeMap]
	client           client.ServicerClient
	storeDirectory   string // directory that will contain session tree stores
}

func NewSessionManager(
	ctx context.Context,
	storeDirectory string,
	client client.ServicerClient,
) *SessionManager {
	sessions := make(map[uint64]sessionTreeMap)
	sm := &SessionManager{client: client, storeDirectory: storeDirectory, sessions: sessions}
	sm.sessionsNotifee, sm.sessionsNotifier = utils.NewControlledObservable[sessionTreeMap](nil)

	go sm.handleBlocks(ctx)

	return sm
}

// emits all sessions that have ended
func (sm *SessionManager) Sessions() utils.Observable[sessionTreeMap] {
	return sm.sessionsNotifee
}

// returns the tree for a given session, creating it if it doesn't exist
func (sm *SessionManager) EnsureSessionTree(sessionInfo *sessionTypes.Session) *smt.SMST {
	// get session end so we can group sessions that end at the same height
	// make sure we do not off by one
	sessionId := sessionInfo.SessionId

	// make sure to have a container for sessions that end at this height
	if _, ok := sm.sessions[sessionInfo.GetSessionEndHeight()]; !ok {
		sm.sessions[sessionInfo.GetSessionEndHeight()] = make(sessionTreeMap)
	}

	// create session tree if it doesn't exist (first relay for this session)
	// we need to get its store so we can close it later since we can't access it from the tree
	if _, ok := sm.sessions[sessionInfo.GetSessionEndHeight()][sessionId]; !ok {
		storePath := filepath.Join(sm.storeDirectory, sessionId)
		tree, store, err := sm.createTreeForSession(storePath)
		if err != nil {
			log.Println("error creating tree for session", err)
			return nil
		}

		sm.sessions[sessionInfo.GetSessionEndHeight()][sessionId] = NewSessionWithTree(
			sessionInfo,
			tree,
			store,
			storePath,
			sm.sessionCleanupFactory(sessionInfo),
		)
	}

	return sm.sessions[sessionInfo.GetSessionEndHeight()][sessionId].SessionTree()
}

func (sm *SessionManager) handleBlocks(ctx context.Context) {
	ch := sm.client.BlocksNotifee().Subscribe(ctx).Ch()
	for block := range ch {
		// if some sessions end by this block, process them
		if sessions, ok := sm.sessions[block.Height()]; ok {
			sm.sessionsNotifier <- sessions
		}
	}
}

func (sm *SessionManager) createTreeForSession(storePath string) (*smt.SMST, smt.KVStore, error) {
	treeStore, err := smt.NewKVStore(storePath)
	if err != nil {
		return nil, nil, err
	}

	tree := smt.NewSparseMerkleSumTree(treeStore, sha256.New(), smt.WithValueHasher(nil))
	return tree, treeStore, nil
}

func (sm *SessionManager) sessionCleanupFactory(sessionInfo *sessionTypes.Session) func() {
	return func() {
		delete(sm.sessions[sessionInfo.GetSessionEndHeight()], sessionInfo.SessionId)

		// delete sessionEnd map if it's empty
		if len(sm.sessions[sessionInfo.GetSessionEndHeight()]) == 0 {
			delete(sm.sessions, sessionInfo.GetSessionEndHeight())
		}
	}
}
