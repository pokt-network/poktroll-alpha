package prooflifecycle

import (
	"context"
	"hash"
	"poktroll/modules"
	"poktroll/runtime/di"

	"github.com/pokt-network/smt"
)

var SMTHasherToken = di.NewInjectionToken[hash.Hash]("SMTHasher")
var SMTStoreToken = di.NewInjectionToken[smt.KVStore]("SMTStore")

type miner struct {
	smst   smt.SMST
	client modules.PocketNetworkClient
	logger modules.Logger
}

func CreateProofLifecycle() modules.MinerModule {
	return &miner{}
}

func (m *miner) SubmitClaim() error {
	claim := m.smst.Root()
	result := <-m.client.SubmitClaim(context.TODO(), claim)
	return result.Error()
}

func (m *miner) SubmitProof(hash []byte) error {
	key, valueHash, sum, proof, err := m.smst.ProveClosest(hash)
	if err != nil {
		return err
	}

	result := <-m.client.SubmitProof(context.TODO(), key, valueHash, sum, proof, err)
	return result.Error()
}

func (m *miner) Update(key []byte, value []byte, weight uint64) error {
	return m.smst.Update(key, value, weight)
}

func (m *miner) Resolve(injector *di.Injector, path *[]string) {
	hasher := di.Resolve(SMTHasherToken, injector, path)
	store := di.Resolve(SMTStoreToken, injector, path)
	m.smst = *smt.NewSparseMerkleSumTree(store, hasher)

	m.logger = *di.Resolve(modules.LoggerModuleToken, injector, path).
		CreateLoggerForModule(modules.MinerModuleToken.Id())

	m.client = di.Resolve(modules.PocketNetworkClientToken, injector, path)
}

func (t *miner) Start() error {
	return nil
}

func (m *miner) CascadeStart() error {
	if err := m.client.CascadeStart(); err != nil {
		return err
	}

	return m.Start()
}
