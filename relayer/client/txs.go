package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	tmtypes "github.com/cometbft/cometbft/types"
	"poktroll/utils"
	"poktroll/x/servicer/types"
)

var (
	_           types.Block = &cometBlockWebsocketMsg{}
	errNotTxMsg             = "expected tx websocket msg; got: %s"
)

// cometBlockWebsocketMsg is used to deserialize incoming websocket messages from
// the block subscription. It implements the types.Block interface by loosely
// wrapping cometbft's block type, into which messages are deserialized.
type cometTxResponseWebsocketMsg struct {
	Tx []byte `json:"tx"`
}

func (client *servicerClient) subscribeToOwnTxs(
	ctx context.Context,
	blocksNotifee utils.Observable[types.Block],
) {
	query := fmt.Sprintf("tm.event='Tx' AND message.sender='%s'", client.address)

	// TODO_CONSIDERATION: using an observable for received tx messages & a filter
	// for `#signAndBroadcastTx()` callers to react to the specific tx in question
	// instead of using shared memory across goroutines (`txByHash`) would likely
	// improve readability and maintainability. This would likely require a new
	// "buffered controllable observable"; i.e. a controlled observable which uses
	// buffered channels to avoid blocking channel sender.
	//
	//txsNotifee, txsNotifier := utils.NewControlledObservable[*cosmosTypes.TxResponse](nil)
	msgHandler := client.handleTxsFactory()
	client.subscribeWithQuery(ctx, query, msgHandler)

	//return txsNotifee
	go client.timeoutTxs(ctx, blocksNotifee)

	return
}

func (client *servicerClient) timeoutTxs(
	ctx context.Context,
	blocksNotifee utils.Observable[types.Block],
) {
	ch := blocksNotifee.Subscribe().Ch()
	for block := range ch {
		select {
		case <-ctx.Done():
			return
		default:
		}

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
		}

		delete(client.txsByHashByTimeout, block.Height())
		client.txsMutex.Unlock()
	}
}

// func handleTxsFactory(txsNotifier chan *cosmosTypes.TxResponse) messageHandler {
func (client *servicerClient) handleTxsFactory() messageHandler {
	return func(ctx context.Context, msg []byte) error {
		txMsg, err := client.newCometTxResponseMsg(msg)
		expectedErr := fmt.Errorf(errNotTxMsg, string(msg))
		switch {
		case err == nil:
		case err.Error() == expectedErr.Error():
			return nil
		case err != nil:
			return fmt.Errorf("failed to parse new tx message: %w", err)
		}

		fmt.Printf("TX MSG:\n%s\n", string(msg))
		txHash := fmt.Sprintf("%x", string(tmtypes.Tx(txMsg.Tx).Hash()))

		client.txsMutex.Lock()
		defer client.txsMutex.Unlock()

		txErrCh, ok := client.txsByHash[txHash]
		if !ok {
			panic("txErrCh not found")
		}
		if txErrCh == nil {
			// INCOMPLETE: handle and/or invalidate this case.
			fmt.Println("txErrCh is nil")
			return nil
		}
		// TODO_THIS_COMMIT: check tx for errors, parse & send if present!!!
		txErrCh <- nil
		close(txErrCh)
		delete(client.txsByHash, txHash)

		// TODO_CONSIDERATION: do we really need both of these maps?
		for timeoutHeight, txsByHash := range client.txsByHashByTimeout {
			for txHash, _ := range txsByHash {
				if txHash == txHash {
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

func (client *servicerClient) newCometTxResponseMsg(txMsgBz []byte) (*cometTxResponseWebsocketMsg, error) {
	txResponseMsg := new(cometTxResponseWebsocketMsg)
	if err := json.Unmarshal(txMsgBz, txResponseMsg); err != nil {
		return nil, err
	}

	// If msg does not match the expected format then block will be its zero value.
	if bytes.Equal(txResponseMsg.Tx, []byte{}) {
		return nil, fmt.Errorf(errNotTxMsg, string(txMsgBz))
	}

	return txResponseMsg, nil
}