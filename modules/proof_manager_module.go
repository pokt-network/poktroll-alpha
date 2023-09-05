package modules

import (
	"github.com/pokt-network/smt"

	"poktroll/runtime/di"
)

var ProofManagerToken = di.NewInjectionToken[ProofManager]("proofManager")

type ProofManager interface {
	di.Module
	Spec() *smt.TreeSpec
	Get(key []byte) ([]byte, uint64, error)
	Update(key, value []byte, weight uint64) error
	Delete(key []byte) error
	Prove(key []byte) (*smt.SparseMerkleProof, error)
	Commit() error
	Root() []byte
	Sum() uint64
}
