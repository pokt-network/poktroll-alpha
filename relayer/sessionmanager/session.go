package sessionmanager

import (
	"crypto/sha256"
	"log"
	"os"
	"sync"

	"github.com/pokt-network/smt"

	"poktroll/x/session/types"
)

var _ SessionWithTree = &sessionWithTree{}

type SessionWithTree interface {
	GetSessionId() string
	SessionTree() *smt.SMST
	CloseTree() ([]byte, error)
	DeleteTree() error
}

type sessionWithTree struct {
	sessionInfo    *types.Session
	tree           *smt.SMST
	treeStore      smt.KVStore
	claimedSMTRoot []byte
	isClosed       bool   // TODO_COMMENT: What does this mean? E.g. can no more relays be added to it?
	storePath      string // TODO_CONSIDERATION: Can this not be part of treeStore?
	onDelete       func() // TODO_CONSIDERATION: onDeleteFn?
	sessionMutex   *sync.Mutex
}

func NewSessionWithTree(
	sessionInfo *types.Session,
	tree *smt.SMST,
	treeStore smt.KVStore,
	storePath string,
	onDelete func(),
) *sessionWithTree {
	return &sessionWithTree{
		sessionInfo: sessionInfo,
		tree:        tree,
		treeStore:   treeStore,
		storePath:   storePath,
		onDelete:    onDelete,
		isClosed:    false,
	}
}

func (s *sessionWithTree) SessionTree() *smt.SMST {
	// if the tree is closed, we need to re-open it from disk
	if s.isClosed {
		store, err := smt.NewKVStore(s.storePath)
		if err != nil {
			log.Println("error creating store for session", err)
			return nil
		}

		s.treeStore = store
		s.tree = smt.ImportSparseMerkleSumTree(s.treeStore, sha256.New(), s.claimedSMTRoot)
	}

	return s.tree
}

func (s *sessionWithTree) GetSessionId() string {
	return s.sessionInfo.SessionId
}

// get the root of the no longer updatable tree
func (s *sessionWithTree) CloseTree() (root []byte, err error) {
	claimedRoot := s.tree.Root()

	// we need the claimed root so we can re-open the tree from disk for proof submission
	s.claimedSMTRoot = claimedRoot

	if err := s.tree.Commit(); err != nil {
		return nil, err
	}

	if err := s.treeStore.Stop(); err != nil {
		return nil, err
	}

	// mark tree/kvstore as closed
	s.isClosed = true
	return claimedRoot, nil
}

func (s *sessionWithTree) DeleteTree() error {
	// DISCUSS: why use treeStore.ClearAll() instead of os.RemoveAll?
	// see: https://pkg.go.dev/os#RemoveAll
	if err := s.treeStore.ClearAll(); err != nil {
		return err
	}

	if err := s.treeStore.Stop(); err != nil {
		return err
	}

	if err := os.RemoveAll(s.storePath); err != nil {
		return err
	}

	s.onDelete()

	return nil
}
