package sessionmanager

import (
	"crypto/sha256"
	"log"
	"os"

	"github.com/pokt-network/smt"

	"poktroll/x/session/types"
)

type SessionWithTree interface {
	SessionTree() *smt.SMST
	CloseTree() ([]byte, error)
	DeleteTree() error
}

var _ SessionWithTree = &sessionWithTree{}

type sessionWithTree struct {
	sessionInfo *types.Session
	tree        *smt.SMST
	treeStore   smt.KVStore
	claimedRoot []byte
	closed      bool
	storePath   string
	onDelete    func()
}

func (s *sessionWithTree) SessionTree() *smt.SMST {
	// if the tree is closed, we need to re-open it from disk
	if s.closed {
		store, err := smt.NewKVStore(s.storePath)
		if err != nil {
			log.Println("error creating store for session", err)
			return nil
		}

		s.treeStore = store
		s.tree = smt.ImportSparseMerkleSumTree(s.treeStore, sha256.New(), s.claimedRoot)
	}

	return s.tree
}

// get the root of the no longer updatable tree
func (s *sessionWithTree) CloseTree() (root []byte, err error) {
	claimedRoot := s.tree.Root()

	// we need the claimed root so we can re-open the tree from disk for proof submission
	s.claimedRoot = claimedRoot

	if err := s.tree.Commit(); err != nil {
		return nil, err
	}

	if err := s.treeStore.Stop(); err != nil {
		return nil, err
	}

	// mark tree/kvstore as closed
	s.closed = true
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
