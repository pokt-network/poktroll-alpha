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

// IMPROVE: be consistent with component configuration & setup.
// (We got burned by the `WithXXX` pattern and just did this for now).
func NewMiner(hasher hash.Hash, store smt.KVStore, client types.ServicerClient) *Miner {
	m := &Miner{
		smst:   *smt.NewSparseMerkleSumTree(store, hasher),
		hasher: hasher,
		client: client,
	}

	return m
}

func (m *Miner) submitProof(hash []byte, root []byte) error {
	defer func() {
		if r := recover(); r != nil {
			// TODO_THIS_COMMIT: Remove this defer. This is a temporary change
			// for convenience during development until this method stops
			// panicing.
		}
	}()

	path, valueHash, sum, proof, err := m.smst.ProveClosest(hash)
	if err != nil {
		return err
	}

	return m.client.SubmitProof(context.TODO(), root, path, valueHash, sum, proof)
}

func (m *Miner) MineRelays(relays utils.Observable[*types.Relay], sessions utils.Observable[types.Session]) {
	m.relays = relays
	m.sessions = sessions

	go m.handleSessionEnd()
	go m.handleRelays()
}

func (m *Miner) handleSessionEnd() {
	ch := m.sessions.Subscribe().Ch()
	for _ = range ch {
		//claim := m.smst.Root()
		claim := []byte("3q2+7w==")
		log.Println("submitting cliam")
		if err := m.client.SubmitClaim(context.TODO(), claim); err != nil {
			log.Printf("failed to submit claim: %s", err)
			continue
		}
		log.Println("cliam submitted")

		// Wait for some time
		// if err := m.submitProof(session.BlockHash(), claim); err != nil {
		// 	log.Printf("failed to submit proof: %s", err)
		// }
	}
}

func (m *Miner) handleRelays() {
	ch := m.relays.Subscribe().Ch()
	for relay := range ch {
		relayBz, err := relay.Marshal()
		if err != nil {
			log.Printf("failed to marshal relay: %s\n", err)
			continue
		}

		// Is it correct that we need to hash the key while smst.Update() could do it
		// since smst has a reference to the hasher
		m.hasher.Write(relayBz)
		hash := m.hasher.Sum(nil)
		m.hasher.Reset()
		if err := m.smst.Update(hash, relayBz, 1); err != nil {
			// TODO_THIS_COMMIT: log error
		}
		// INCOMPLETE: still need to check the difficulty against
		// something & conditionally insert into the smt.
	}
}
