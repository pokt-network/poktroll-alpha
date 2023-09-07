package prooflifecycle

import (
	"context"
	"hash"

	"github.com/pokt-network/smt"

	"poktroll/modules"
	"poktroll/runtime/di"
	"poktroll/types"
)

var SMTHasherToken = di.NewInjectionToken[hash.Hash]("SMTHasher")
var SMTStoreToken = di.NewInjectionToken[smt.KVStore]("SMTStore")

type miner struct {
	smst     smt.SMST
	client   modules.PocketNetworkClient
	logger   modules.Logger
	relays   <-chan *types.Relay
	sessions <-chan *types.Session
}

func CreateMiner() modules.MinerModule {
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

func (m *miner) MineRelays(relays <-chan *types.Relay, sessions <-chan *types.Session) {
	m.relays = relays
	m.sessions = sessions
}

func (m *miner) handleSessionEnd() {
	for session := range m.sessions {
		if err := m.SubmitClaim(); err != nil {
			continue
		}

		// Wait for some time
		m.SubmitProof([]byte(session.BlockHash))
	}
}

func (m *miner) handleRelays() {
	for relay := range m.relays {
		m.logger.Debug().Msgf("TODO handle relay 🔂 %+v", relay)

		// TODO get access to a relayer module

		// TODO get the serialized byte representation of the relay
		// serializedRelay := relay.Serialize()

		// TODO update the claim and proof tree
		// key=SHA3HASH(serializedRelay)
		// value=serializedRelay
		// hash := crypto.SHA3Hash([]byte(serializedRelay))
		// m.Update(hash, hash, 1)
	}
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
	go t.handleSessionEnd()
	go t.handleRelays()
	return nil
}

func (m *miner) CascadeStart() error {
	if err := m.client.CascadeStart(); err != nil {
		return err
	}

	return m.Start()
}
