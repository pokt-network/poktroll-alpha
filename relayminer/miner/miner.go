package miner

import (
	"hash"

	"github.com/pokt-network/smt"

	"poktroll/utils"
	"poktroll/x/servicer/client"
	"poktroll/x/servicer/types"
)

type Miner struct {
	smst smt.SMST
	// TECHDEBT: update after switching to logger module (i.e. once
	// servicer is external to poktrolld)
	//logger   log.Logger
	relays   utils.Observable[*types.Relay]
	sessions utils.Observable[*types.Session]
	client   client.ServicerClient

	hasher hash.Hash
}

func NewMiner(hasher hash.Hash, store smt.KVStore, client client.ServicerClient) *Miner {
	m := &Miner{
		smst:   *smt.NewSparseMerkleSumTree(store, hasher),
		hasher: hasher,
		client: client,
	}

	go m.handleSessionEnd()
	go m.handleRelays()

	return m
}

func (m *Miner) submitClaim() error {
	//claim := m.smst.Root()
	//result := m.client.SubmitClaim(context.TODO(), claim)
	//return result.Error()
	return nil
}

func (m *Miner) submitProof(hash []byte) error {
	//key, valueHash, sum, proof, err := m.smst.ProveClosest(hash)
	//if err != nil {
	//	return err
	//}

	//result := m.client.SubmitProof(context.TODO(), key, valueHash, sum, proof, err)
	//return result.Error()
	return nil
}

func (m *Miner) MineRelays(relays utils.Observable[*types.Relay], sessions utils.Observable[*types.Session]) {
	m.relays = relays
	m.sessions = sessions
}

func (m *Miner) handleSessionEnd() {
	ch := m.sessions.Subscribe().Ch()
	for session := range ch {
		if err := m.submitClaim(); err != nil {
			continue
		}

		// Wait for some time
		m.submitProof([]byte(session.BlockHash))
	}
}

func (m *Miner) handleRelays() {
	ch := m.relays.Subscribe().Ch()
	for relay := range ch {
		//m.logger.Info("TODO handle relay 🔂 %+v", relay)

		// TODO get the serialized byte representation of the relay
		relayBz, err := relay.Marshal()
		if err != nil {
			//m.logger.Error("failed to marshal relay: %s", err)
			continue
		}

		// Is it correct that we need to hash the key while smst.Update() could do it
		// since smst has a reference to the hasher
		hash := m.hasher.Sum([]byte(relayBz))
		m.update(hash, relayBz, 1)
	}
}

func (m *Miner) update(key []byte, value []byte, weight uint64) error {
	return m.smst.Update(key, value, weight)
}
