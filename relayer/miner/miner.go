package miner

import (
	"context"
	"hash"
	"log"

	"poktroll/relayer/proxy"
	"poktroll/relayer/sessionmanager"
	"poktroll/utils"
	"poktroll/x/servicer/types"
)

type Miner struct {
	relays         utils.Observable[*proxy.RelayWithSession]
	sessionManager *sessionmanager.SessionManager
	client         types.ServicerClient
	hasher         hash.Hash
}

// IMPROVE: be consistent with component configuration & setup.
// (We got burned by the `WithXXX` pattern and just did this for now).
func NewMiner(
	hasher hash.Hash,
	client types.ServicerClient,
	sessionManager *sessionmanager.SessionManager,
) *Miner {
	m := &Miner{
		hasher:         hasher,
		client:         client,
		sessionManager: sessionManager,
	}

	return m
}

func (m *Miner) MineRelays(ctx context.Context, relays utils.Observable[*proxy.RelayWithSession]) {
	m.relays = relays

	// these methods block, waiting for new sessions and relays respectively.
	go m.handleSessions(ctx)
	go m.handleRelays(ctx)
}

func (m *Miner) handleSessions(ctx context.Context) {
	ch := m.sessionManager.Sessions().Subscribe().Ch()
	// this emits each time a batch of sessions is ready to be processed.
	for closedSessions := range ch {
		// process sessions in parallel.
		for _, session := range closedSessions {
			go m.handleSingleSession(ctx, session)
		}
	}
}

func (m *Miner) handleSingleSession(ctx context.Context, session sessionmanager.SessionWithTree) {
	// this session should no longer be updated
	claimRoot, err := session.CloseTree()
	if err != nil {
		log.Printf("failed to close tree: %s", err)
		return
	}

	// SubmitClaim ensures on-chain claim inclusion
	if err := m.client.SubmitClaim(ctx, claimRoot); err != nil {
		log.Printf("failed to submit claim: %s", err)
		return
	}

	// TODO: implement wait logic here

	// generate and submit proof
	// past this point the proof is included on-chain and the session can be pruned.
	if err := m.submitProof(ctx, session, claimRoot); err != nil {
		log.Printf("failed to submit proof: %s", err)
		return
	}

	// prune tree now that proof is submitted
	if err := session.PruneTree(); err != nil {
		log.Printf("failed to prune tree: %s", err)
		return
	}
}

func (m *Miner) handleRelays(ctx context.Context) {
	ch := m.relays.Subscribe().Ch()
	// process each relay in parallel
	for relay := range ch {
		go m.handleSingleRelay(relay)
	}
}

func (m *Miner) handleSingleRelay(relayWithSession *proxy.RelayWithSession) {
	relayBz, err := relayWithSession.Relay.Marshal()
	if err != nil {
		log.Printf("failed to marshal relay: %s\n", err)
		return
	}

	// Is it correct that we need to hash the key while smst.Update() could do it
	// since smst has a reference to the hasher
	m.hasher.Write(relayBz)
	hash := m.hasher.Sum(nil)
	m.hasher.Reset()

	// ensure the session tree exists for this relay
	smst := m.sessionManager.EnsureSessionTree(relayWithSession.Session)

	if err := smst.Update(hash, relayBz, 1); err != nil {
		log.Printf("failed to update smt: %s\n", err)
		return
	}
	// INCOMPLETE: still need to check the difficulty against
	// something & conditionally insert into the smt.
}

func (m *Miner) submitProof(ctx context.Context, session sessionmanager.SessionWithTree, claimRoot []byte) error {
	defer func() {
		if r := recover(); r != nil {
			// TODO_THIS_COMMIT: Remove this defer. This is a temporary change
			// for convenience during development until this method stops
			// panicing.
		}
	}()

	// at this point the miner already waited for a number of blocks
	// use the latest block hash as the key to prove against.
	currentBlockHash := m.client.GetLatestBlock().Hash()
	path, valueHash, sum, proof, err := session.SessionTree().ProveClosest(currentBlockHash)
	if err != nil {
		return err
	}

	// SubmitProof ensures on-chain proof inclusion so we can safely prune the tree.
	err = m.client.SubmitProof(ctx, claimRoot, path, valueHash, sum, proof)
	if err != nil {
		return err
	}

	return nil
}
