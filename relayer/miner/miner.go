package miner

import (
	"context"
	"hash"
	"log"

	"poktroll/relayer/client"
	"poktroll/relayer/proxy"
	"poktroll/relayer/sessionmanager"
	"poktroll/utils"
)

type Miner struct {
	relays         utils.Observable[*proxy.RelayWithSession]
	sessionManager *sessionmanager.SessionManager
	client         client.ServicerClient
	hasher         hash.Hash
}

// IMPROVE: be consistent with component configuration & setup.
// (We got burned by the `WithXXX` pattern and just did this for now).
func NewMiner(
	hasher hash.Hash,
	client client.ServicerClient,
	sessionManager *sessionmanager.SessionManager,
) *Miner {
	m := &Miner{
		hasher:         hasher,
		client:         client,
		sessionManager: sessionManager,
	}

	return m
}

// MineRelays assigns the relays and sessions observables & starts their
// respective consumer goroutines.
func (m *Miner) MineRelays(ctx context.Context, relays utils.Observable[*proxy.RelayWithSession]) {
	m.relays = relays

	// these methods block, waiting for new sessions and relays respectively.
	go m.handleSessions(ctx)
	go m.handleRelays(ctx)
}

// handleSessionEnd submits a claim for the ended session & starts a goroutine
// which will submit the corresponding proof when the respective proof window
// opens.
func (m *Miner) handleSessions(ctx context.Context) {
	subscription := m.sessionManager.Sessions().Subscribe()
	go func() {
		<-ctx.Done()
		subscription.Unsubscribe()
	}()

	ch := subscription.Ch()
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

	// generate and submit proof
	// past this point the proof is included on-chain and the session can be pruned.
	if err := m.waitAndProve(ctx, session, claimRoot); err != nil {
		log.Printf("failed to submit proof: %s", err)
		return
	}

	// prune tree now that proof is submitted
	if err := session.DeleteTree(); err != nil {
		log.Printf("failed to prune tree: %s", err)
		return
	}
}

// handleRelays blocks until a relay is received, then handles it in a new
// goroutine.
func (m *Miner) handleRelays(ctx context.Context) {
	subscription := m.relays.Subscribe()
	go func() {
		<-ctx.Done()
		subscription.Unsubscribe()
	}()

	ch := subscription.Ch()
	// process each relay in parallel
	for relay := range ch {
		go m.handleSingleRelay(ctx, relay)
	}
}

// handleSingleRelay validates, executes, & hashes the relay. If the relay's difficulty
// is above the mining difficulty, it's inserted into SMST.
func (m *Miner) handleSingleRelay(
	ctx context.Context,
	relayWithSession *proxy.RelayWithSession,
) {
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

	// INCOMPLETE: still need to check the difficulty against
	// something & conditionally insert into the smt.
	if err := smst.Update(hash, relayBz, 1); err != nil {
		log.Printf("failed to update smt: %s\n", err)
		return
	}
}

func (m *Miner) waitAndProve(ctx context.Context, session sessionmanager.SessionWithTree, claimRoot []byte) error {
	// TODO: implement wait logic here

	// at this point the miner already waited for a number of blocks
	// use the latest block hash as the key to prove against.
	currentBlockHash := m.client.LatestBlock().Hash()
	path, valueHash, sum, proof, err := session.SessionTree().ProveClosest(currentBlockHash)
	if err != nil {
		return err
	}

	// SubmitProof ensures on-chain proof inclusion so we can safely prune the tree.
	return m.client.SubmitProof(ctx, claimRoot, path, valueHash, sum, proof)
}
