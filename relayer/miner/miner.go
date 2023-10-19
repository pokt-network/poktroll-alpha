package miner

import (
	"context"
	"encoding/binary"
	"fmt"
	"hash"
	"log"
	"math/rand"
	"sync"

	"poktroll/relayer/client"
	"poktroll/relayer/proxy"
	"poktroll/relayer/sessionmanager"
	"poktroll/utils"
	claimproofparams "poktroll/x/servicer/keeper"
	"poktroll/x/servicer/types"
	sessionkeeper "poktroll/x/session/keeper"
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
	// TODO: Fetch and process incomplete sessions before listening to new ones
	// Add persistence to pending claims and proofs of a given session
	// Use m.handleSingleSession for each of them
	ch := m.sessionManager.Sessions().Subscribe(ctx).Ch()
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

	// TODO: Query if a claim has been created for this session and skip waitAndClaim if so.

	// SubmitClaim ensures on-chain claim inclusion
	if err := m.waitAndClaim(ctx, session, claimRoot); err != nil {
		log.Printf("failed to submit claim: %s", err)
		return
	}

	// HACK: this is a hack to get the claim submission block height assuming that getting a block
	// right after the claim is submitted will return the block at which the claim was submitted.
	// seems that m.client.LatestBlock() is not returning the latest block.
	claimSubmissionBlockHeight := int64(m.client.LatestBlock(ctx).Height() + 1)
	log.Printf("claim submitted at block height: %d", claimSubmissionBlockHeight)

	// TODO_REFACTOR: This should happen a few blocks later. WAIT should be an async process
	// based on block events. Prove should just be one atomic method.

	// generate and submit proof
	// past this point the proof is included on-chain and the local session can be cleared.
	if err := m.waitAndProve(ctx, session, claimRoot, claimSubmissionBlockHeight); err != nil {
		log.Printf("failed to submit proof: %s", err)
		return
	}

	log.Printf("proof submitted at block height: %d", m.client.LatestBlock(ctx).Height())

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
	ch := m.relays.Subscribe(ctx).Ch()
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

// getClaimSubmissionWindow returns the earliest and latest block heights at which
// a claim can be submitted for the current session.
// explanation of how the earliest and latest submission block height is determined is available in the
// poktroll/x/servicer/keeper/msg_server_claim.go file
func (m *Miner) getClaimSubmissionWindow(ctx context.Context, sessionNumber int64) (earliest int64, latest int64, err error) {
	numSessionBlocks := int64(sessionkeeper.NumSessionBlocks)

	sessionStartHeight := sessionNumber * numSessionBlocks
	earliestClaimSubmissionBlockHeight := sessionStartHeight + numSessionBlocks + claimproofparams.GovEarliestClaimSubmissionBlocksOffset

	// we wait for earliestClaimSubmissionBlockHeight to be received before proceeding since we need its hash
	// to know where this servicer's claim submission window starts.
	log.Printf("waiting for global earliest claim submission block height: %d", earliestClaimSubmissionBlockHeight)
	block, err := m.waitForBlock(ctx, uint64(earliestClaimSubmissionBlockHeight))
	if err != nil {
		return 0, 0, err
	}

	log.Printf("received earliest claim submission block height: %d, use its hash to have a random submission for the servicer", block.Height())

	earliestClaimSubmissionBlockHash := block.Hash()
	log.Printf("using block %d's hash %x as randomness", block.Height(), earliestClaimSubmissionBlockHash)
	rngSeed, _ := binary.Varint(earliestClaimSubmissionBlockHash)
	randomNumber := rand.NewSource(rngSeed).Int63()
	randClaimSubmissionBlockHeightOffset := randomNumber % (claimproofparams.GovLatestClaimSubmissionBlocksInterval - claimproofparams.GovClaimSubmissionBlocksWindow - 1)
	earliestServicerClaimSubmissionBlockHeight := earliestClaimSubmissionBlockHeight + randClaimSubmissionBlockHeightOffset + 1
	latestServicerClaimSubmissionBlockHeight := earliestServicerClaimSubmissionBlockHeight + claimproofparams.GovClaimSubmissionBlocksWindow

	return earliestServicerClaimSubmissionBlockHeight, latestServicerClaimSubmissionBlockHeight, nil
}

func (m *Miner) waitAndClaim(
	ctx context.Context,
	session sessionmanager.SessionWithTree,
	claimRoot []byte,
) error {
	earliestBlockHeight, _, err := m.getClaimSubmissionWindow(ctx, int64(session.GetSessionInfo().SessionNumber))
	if err != nil {
		return err
	}

	log.Printf("earliest claim submission block height for this servicer: %d", earliestBlockHeight)
	block, err := m.waitForBlock(ctx, uint64(earliestBlockHeight-1))
	if err != nil {
		return err
	}
	log.Printf("currentBlock: %d, submitting claim", block.Height())
	return m.client.SubmitClaim(ctx, session.GetSessionInfo(), claimRoot)
}

// getProofSubmissionWindow returns the earliest and latest block heights at which
// a proof can be submitted.
// explanation of how the earliest and latest submission block height is determined is available in the
// poktroll/x/servicer/keeper/msg_server_claim.go file
func (m *Miner) getProofSubmissionWindow(ctx context.Context, claimSubmissionHeight int64) (earliest int64, latest int64, err error) {
	earliestProofSubmissionBlockHeight := claimSubmissionHeight + claimproofparams.GovEarliestProofSubmissionBlocksOffset

	// we wait for earliestProofSubmissionBlockHeight to be received before proceeding since we need its hash
	log.Printf("waiting for global earliest proof submission block height: %d", earliestProofSubmissionBlockHeight)
	block, err := m.waitForBlock(ctx, uint64(earliestProofSubmissionBlockHeight))
	if err != nil {
		return 0, 0, err
	}
	earliestProofSubmissionBlockHash := block.Hash()
	log.Printf("using block %d's hash %x as randomness", block.Height(), earliestProofSubmissionBlockHash)

	rngSeed, _ := binary.Varint(earliestProofSubmissionBlockHash)
	randomNumber := rand.NewSource(rngSeed).Int63()
	randProofSubmissionBlockHeightOffset := randomNumber % (claimproofparams.GovLatestProofSubmissionBlocksInterval - claimproofparams.GovProofSubmissionBlocksWindow - 1)
	earliestServicerProofSubmissionBlockHeight := earliestProofSubmissionBlockHeight + randProofSubmissionBlockHeightOffset + 1

	latestServicerClaimSubmissionBlockHeight := earliestProofSubmissionBlockHeight + claimproofparams.GovProofSubmissionBlocksWindow

	return earliestServicerProofSubmissionBlockHeight, latestServicerClaimSubmissionBlockHeight, nil
}

func (m *Miner) waitAndProve(ctx context.Context, session sessionmanager.SessionWithTree, claimRoot []byte, claimSubmissionHeight int64) error {
	// TODO_REFACTOR: Make wait logic asynchronous and event based.

	// use the servicer's earliest proof submission block as the key to prove against.
	// this blocks until the earliest submission block height is received and its hash known.
	earliestBlockHeight, _, err := m.getProofSubmissionWindow(ctx, claimSubmissionHeight)
	if err != nil {
		return err
	}

	log.Printf("earliest proof submission block height for this servicer: %d", earliestBlockHeight)
	block, err := m.waitForBlock(ctx, uint64(earliestBlockHeight-1))
	if err != nil {
		return err
	}

	proofBlockHash := block.Hash()
	log.Printf("used session id %s", session.GetSessionId())
	sessionTree := session.SessionTree()
	proof, err := sessionTree.ProveClosest(proofBlockHash)
	if err != nil {
		return err
	}

	log.Printf("currentBlock: %d, submitting proof", block.Height()+1)
	// SubmitProof ensures on-chain proof inclusion so we can safely prune the tree.
	return m.client.SubmitProof(
		ctx,
		session.GetSessionId(),
		claimRoot,
		proof,
	)
}

// waitForBlock blocks until the block at the given height is received.
func (m *Miner) waitForBlock(ctx context.Context, height uint64) (types.Block, error) {
	currentBlock := m.client.LatestBlock(ctx)
	if currentBlock.Height() == height {
		return currentBlock, nil
	}

	subscription := m.client.BlocksNotifee().Subscribe(ctx)
	defer subscription.Unsubscribe()
	ch := subscription.Ch()
	for block := range ch {
		if block.Height() == height {
			return block, nil
		} else if block.Height() > height {
			return nil, fmt.Errorf("too late for block %d; current: %d", height, block.Height())
		} else {
			log.Printf("waiting for block %d; current: %d", height, block.Height())
		}
	}

	return nil, fmt.Errorf("block not received")
}
