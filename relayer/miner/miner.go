package miner

import (
	"context"
	"hash"
	"log"

	"github.com/pokt-network/smt"

	"poktroll/utils"
	"poktroll/x/servicer/types"
)

type Miner struct {
	smst     smt.SMST
	relays   utils.Observable[*types.Relay]
	sessions utils.Observable[types.Session]
	client   types.ServicerClient
	hasher   hash.Hash
}

func NewMiner(hasher hash.Hash, store smt.KVStore, client types.ServicerClient) *Miner {
	m := &Miner{
		smst:   *smt.NewSparseMerkleSumTree(store, hasher),
		hasher: hasher,
		client: client,
	}

	go m.handleSessionEnd()
	go m.handleRelays()

	return m
}

func (m *Miner) submitProof(hash []byte, root []byte) error {
	path, valueHash, sum, proof, err := m.smst.ProveClosest(hash)
	if err != nil {
		return err
	}

	return m.client.SubmitProof(context.TODO(), root, path, valueHash, sum, proof)
}

func (m *Miner) MineRelays(relays utils.Observable[*types.Relay], sessions utils.Observable[types.Session]) {
	m.relays = relays
	m.sessions = sessions
}

func (m *Miner) handleSessionEnd() {
	ch := m.sessions.Subscribe().Ch()
	for session := range ch {
		claim := m.smst.Root()
		if err := m.client.SubmitClaim(context.TODO(), claim); err != nil {
			log.Printf("failed to submit claim: %s", err)
			continue
		}

		// Wait for some time
		if err := m.submitProof(session.BlockHash(), claim); err != nil {
			log.Printf("failed to submit proof: %s", err)
		}
	}
}

func (m *Miner) handleRelays() {
	ch := m.relays.Subscribe().Ch()
	for relay := range ch {
		// TODO get the serialized byte representation of the relay
		relayBz, err := relay.Marshal()
		if err != nil {
			//m.logger.Error("failed to marshal relay: %s", err)
			continue
		}

		// Is it correct that we need to hash the key while smst.Update() could do it
		// since smst has a reference to the hasher
		hash := m.hasher.Sum([]byte(relayBz))
		if err := m.smst.Update(hash, relayBz, 1); err != nil {
			// TODO_THIS_COMMIT: log error
		}
	}
}
