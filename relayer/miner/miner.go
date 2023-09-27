package miner

import (
	"context"
	"hash"
	"log"
	"sync"

	"poktroll/relayer/client"
	"poktroll/relayer/proxy"
	"poktroll/relayer/sessionmanager"
	"poktroll/utils"
)

// TODO: https://stackoverflow.com/questions/77190071/golang-best-practice-for-functions-intended-to-be-called-as-goroutines

// TODO_COMMENT: Explain what the responsibility of this structure is, how its used throughout
// and leave comments alongside each field.
type Miner struct {
	relays         utils.Observable[*proxy.RelayWithSession]
	sessionManager *sessionmanager.SessionManager
	client         client.ServicerClient
	hasher         hash.Hash
	sessionsMutex  sync.Mutex
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
		sessionsMutex:  sync.Mutex{},
	}

	return m
}

// MineRelays assigns the relays and sessions observables & starts their respective consumer goroutines.
func (m *Miner) MineRelays(ctx context.Context, relays utils.Observable[*proxy.RelayWithSession]) {
	m.relays = relays

	// these methods block, waiting for new sessions and relays respectively.
	go m.handleSessions(ctx)
	go m.handleRelays(ctx)
}

// handleSessions submits a claim for the ended session & starts a goroutine
// which will submit the corresponding proof when the respective proof window
// opens.
// IMPORTANT: This method is intended to be called as a new goroutine.
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

// TODO_REFACTOR: This function currently submits the claim & proof immediately. In reality, we're going
// to have to refactor the logic here such that:
// 1. Session ends & we submit the claim (TODO: account for session rollover that's configurable by the client)
// 2. Wait a few blocks after the claim is committed (TODO: Make this a governance parameter)
// 3. Build and submit the proof
// 4. Purune the session tree

// IMPORTANT: This method is intended to be called as a new goroutine.
func (m *Miner) handleSingleSession(ctx context.Context, session sessionmanager.SessionWithTree) {
	// this session should no longer be updated
	claimRoot, err := session.CloseTree()
	if err != nil {
		log.Printf("failed to close tree: %s", err)
		return
	}

	// SubmitClaim ensures on-chain claim inclusion
	if err := m.client.SubmitClaim(ctx, session.GetSessionId(), claimRoot); err != nil {
		log.Printf("failed to submit claim: %s", err)
		return
	}

	// TODO_REFACTOR: This should happen a few blocks later. WAIT should be an async process
	// based on block events. Prove should just be one atomic method.

	// generate and submit proof
	// past this point the proof is included on-chain and the local session can be cleared.
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
// IMPORTANT: This method is intended to be called as a new goroutine.
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

	m.sessionsMutex.Lock()
	defer m.sessionsMutex.Unlock()

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
	// TODO_REFACTOR: Make wait logic asynchronous and event based.

	// at this point the miner already waited for a number of blocks
	// use the latest block hash as the key to prove against.
	// TODO
	// 1. Create a function that converts the block hash to a path (to encapsulate the logic)
	// 2. Make sure the height at which we use the hash is deterministic (i.e. claimCommitHeight + govVariable)
	currentBlockHash := m.client.LatestBlock().Hash()
	path, valueHash, sum, proof, err := session.SessionTree().ProveClosest(currentBlockHash)
	if err != nil {
		return err
	}

	// SubmitProof ensures on-chain proof inclusion so we can safely prune the tree.
	return m.client.SubmitProof(ctx, claimRoot, path, valueHash, sum, proof)
}
