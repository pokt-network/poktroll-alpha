package smst

import (
	"hash"

	"github.com/pokt-network/cmt-pokt/modules"
	"github.com/pokt-network/cmt-pokt/runtime/di"
	"github.com/pokt-network/smt"
)

var SMTHasherToken = di.NewInjectionToken[hash.Hash]("SMTHasher")
var SMTStoreToken = di.NewInjectionToken[smt.MapStore]("SMTStore")

type smst struct {
	smt.SMST
}

func NewSMST() modules.ProofManager {
	return &smst{}
}

func NewMapStore() smt.MapStore {
	return smt.NewSimpleMap()
}

func (t *smst) Resolve(injector *di.Injector, path *[]string) {
	hasher := di.Resolve(SMTHasherToken, injector, path)
	store := di.Resolve(SMTStoreToken, injector, path)

	t.SMST = *smt.NewSparseMerkleSumTree(store, hasher)
}

func (t *smst) CascadeStart() error {
	return t.Start()
}

func (t *smst) Start() error {
	return nil
}
