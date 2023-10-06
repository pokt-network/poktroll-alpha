package client

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"poktroll/utils"
	"poktroll/x/servicer/types"
	"strings"

	errorsmod "cosmossdk.io/errors"
	abciTypes "github.com/cometbft/cometbft/abci/types"
	cometTypes "github.com/cometbft/cometbft/types"
)

var (
	errNotTxMsg                = errorsmod.Register("relayer/client", 1, "Expected tx websocket msg")
	errInvalidTimedOutTxHash   = errorsmod.Register("relayer/client", 2, "Invalid time out tx hash")
	errFailedToFetchTimedOutTx = errorsmod.Register("relayer/client", 3, "Failed to fetch error for timed out tx")
	errTxTimeOut               = errorsmod.Register("relayer/client", 4, "Tx timed out")
)

// cometTxResponseWebsocketMsg is used to deserialize incoming websocket messages from
// the tx subscription.
type cometTxResponseWebsocketMsg struct {
	Tx     []byte            `json:"tx"`
	Events []abciTypes.Event `json:"events"`
}

// subscribeToOwnTxs subscribes to txs which were signed/sent by this client's
// address using a single websocket connection.
func (client *servicerClient) subscribeToOwnTxs(
	ctx context.Context,
	blocksNotifee utils.Observable[types.Block],
) {
	// NB: cometbft event subscription query
	// (see: https://docs.cosmos.network/v0.47/core/events#subscribing-to-events)
	query := fmt.Sprintf("tm.event='Tx' AND message.sender='%s'", client.address)

	// TODO_CONSIDERATION: using an observable for received tx messages & a filter
	// for `#signAndBroadcastTx()` callers to react to the specific tx in question
	// instead of using shared memory across goroutines (`txByHash`) would likely
	// improve readability and maintainability. This would likely require a new
	// "buffered controllable observable"; i.e. a controlled observable which uses
	// buffered channels to avoid blocking channel sender.
	//
	//txsNotifee, txsNotifier := utils.NewControlledObservable[*cosmosTypes.TxResponse](nil)
	msgHandler := client.txsFactoryHandler()
	client.subscribeWithQuery(ctx, query, msgHandler)

	//return txsNotifee
	go client.timeoutTxs(ctx, blocksNotifee)
}

// Closes the error channels for expect transactions from the latest block when it times out.
// IMPORTANT: THis is intended to be run as a goroutine!
func (client *servicerClient) timeoutTxs(
	ctx context.Context,
	blocksNotifee utils.Observable[types.Block],
) {
	ch := blocksNotifee.Subscribe().Ch()
	// TODO: Add a comment
	for block := range ch {
		select {
		case <-ctx.Done():
			return
		default:
		}

		// TODO_TECHDEBT: move latest block assignment to a dedicated subscription / goroutine
		// Update latest block
		client.latestBlockMutex.Lock()
		client.latestBlock = block
		client.latestBlockMutex.Unlock()

		client.txsMutex.Lock()
		txsByHash, ok := client.txsByHashByTimeout[block.Height()]
		if !ok {
			// No txs to time out this block height.
			client.txsMutex.Unlock()
			continue
		}

		for txHash, txErrCh := range txsByHash {
			select {
			// If the tx has been seen by the subscription, then the txErrCh
			// will have been closed by the websocket message handler after
			// parsing and sending the error.
			case err, ok := <-txErrCh:
				if ok {
					// TODO_THIS_COMMIT: finish thinking this through.
					panic(fmt.Errorf("txErrCh should be closed; got err: %w", err))
				}
				delete(txsByHash, txHash)
				client.txsMutex.Unlock()
				continue
			default:
			}

			// Otherwise, send a timeout error on, close, and delete txErrCh.
			txErrCh <- fmt.Errorf("tx timed out: %s", txHash)
			close(txErrCh)
			delete(txsByHash, txHash)

			go client.getTxTimeoutError(ctx, txHash)
		}

		delete(client.txsByHashByTimeout, block.Height())
		client.txsMutex.Unlock()
	}
}

// txsFactoryHandler returns a websocket message handler function which attempts
// to deserialize a tx event message, find its corresponding txErrCh, send an
// error if present, & close it.
func (client *servicerClient) txsFactoryHandler() messageHandler {
	return func(ctx context.Context, msg []byte) error {
		txMsg, err := client.newCometTxResponseMsg(msg)
		switch {
		case err == nil:
		case errNotTxMsg.Is(err):
			return nil
		case err != nil:
			return fmt.Errorf("failed to parse new tx message: %w", err)
		}

		fmt.Printf("TX MSG:\n%s\n", string(msg))
		txHash := strings.ToLower(
			fmt.Sprintf("%x", string(
				cometTypes.Tx(txMsg.Tx).Hash(),
			)),
		)

		client.txsMutex.Lock()
		defer client.txsMutex.Unlock()

		txErrCh, ok := client.txsByHash[txHash]
		if !ok {
			log.Println("error: received tx for which no txErrCh exists")
			return nil
		}
		// TODO_THIS_COMMIT: check tx for errors, parse & send if present!!!
		txErrCh <- nil
		close(txErrCh)
		delete(client.txsByHash, txHash)

		// TODO_CONSIDERATION: do we really need both of these maps?
		for timeoutHeight, txsByHash := range client.txsByHashByTimeout {
			for txHashAsKey, _ := range txsByHash {
				if txHash == txHashAsKey {
					delete(txsByHash, txHash)
				}
			}
			if len(txsByHash) == 0 {
				delete(client.txsByHashByTimeout, timeoutHeight)
			}
		}
		return nil
	}
}

// newCometTxResponseMsg attempts to deserialize the given bytes into a comet tx event byte slice.
// if the resulting block has a height of zero, assume the message was not a
// block message and return an errNotBlockMsg error.
func (client *servicerClient) newCometTxResponseMsg(txMsgBz []byte) (*cometTxResponseWebsocketMsg, error) {
	txResponseMsg := new(cometTxResponseWebsocketMsg)
	if err := json.Unmarshal(txMsgBz, txResponseMsg); err != nil {
		return nil, err
	}

	// If msg does not match the expected format then block will be its zero value.
	if bytes.Equal(txResponseMsg.Tx, []byte{}) {
		return nil, errNotTxMsg.Wrapf("got: %s", string(txMsgBz))
	}

	return txResponseMsg, nil
}

// This function is intended to be called as a goroutine
// TODO_DISCUSS: Should it be prefixed with `go`?
func (client *servicerClient) getTxTimeoutError(ctx context.Context, txHashHex string) error {
	txHash, err := hex.DecodeString(txHashHex)
	if err != nil {
		return errInvalidTimedOutTxHash.Wrapf("got: %s", txHashHex)
	}

	txResponse, err := client.clientCtx.Client.Tx(ctx, txHash, false)
	if err != nil {
		return errFailedToFetchTimedOutTx.Wrapf("got tx: %s: %s", txHashHex, err.Error())
	}
	return errTxTimeOut.Wrapf("got: %x: %s", txHashHex, txResponse.TxResult.Log)
}
