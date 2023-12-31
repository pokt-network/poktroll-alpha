package client

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/polydawn/refmt/json"

	"poktroll/utils"
	"poktroll/x/application/types"
)

var ErrNotDelegateMsg = "expected delegate websocket msg"

func NewDelegateMsg(delegateMsgBz []byte) (*types.EventDelegate, error) {
	dMsg := new(types.EventDelegate)
	if err := json.Unmarshal(delegateMsgBz, dMsg); err != nil {
		return nil, err
	}

	log.Print("delegated")

	// If msg does not match the expected format then the address will be empty
	if dMsg.Address == "" {
		return nil, fmt.Errorf(ErrNotDelegateMsg, string(delegateMsgBz))
	}

	return dMsg, nil
}

type DelegateQueryClient struct {
	conn                *websocket.Conn
	delegateNotifee     utils.Observable[*types.EventDelegate]
	latestDelegateMutex *sync.RWMutex
	latestDelegate      *types.EventDelegate
}

func NewDelegateQueryClient(ctx context.Context, endpoint string) (*DelegateQueryClient, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse endpoint: %w", err)
	}

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to dial endpoint: %w", err)
	}

	client := &DelegateQueryClient{
		conn:                conn,
		latestDelegateMutex: &sync.RWMutex{},
	}

	if err := client.subscribeToDelegate(ctx); err != nil {
		return nil, err
	}

	return client, nil
}

func (qc *DelegateQueryClient) DelegateNotifee() utils.Observable[*types.EventDelegate] {
	return qc.delegateNotifee
}

func (qc *DelegateQueryClient) LatestDelegation(ctx context.Context) *types.EventDelegate {
	qc.latestDelegateMutex.RLock()
	defer qc.latestDelegateMutex.RUnlock()
	// block until we have a delegation to return
	if qc.latestDelegate == nil {
		subscription := qc.delegateNotifee.Subscribe(ctx)
		<-subscription.Ch()
		subscription.Unsubscribe()
	}
	return qc.latestDelegate
}

func (qc *DelegateQueryClient) listen(ctx context.Context, blocksNotifier chan *types.EventDelegate) {
	for {
		_, msg, err := qc.conn.ReadMessage()
		if err != nil {
			return
		}

		if strings.TrimSpace(string(msg)) == `{"jsonrpc":"2.0","result":{},"id":1}` {
			continue
		}

		dMsg, err := NewDelegateMsg(msg)
		expectedErr := fmt.Errorf(ErrNotDelegateMsg, string(msg))
		switch {
		case err == nil:
		case err.Error() == expectedErr.Error():
		case err != nil:
			log.Printf("failed to handle websocket msg: %s\n", err)
		}

		qc.latestDelegate = dMsg
		blocksNotifier <- dMsg
	}
}

func (qc *DelegateQueryClient) subscribeToDelegate(ctx context.Context) error {
	delegateNotifee, delegateNotifier := utils.NewControlledObservable[*types.EventDelegate](nil)

	subscribeMsg := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "subscribe",
		"id":      1,
		"params": map[string]interface{}{
			// subscribe to all events from the application module
			"query": "tm.event='Tx' AND message.module='application'",
		},
	}

	subscribeMsgBz, err := json.Marshal(subscribeMsg)
	if err != nil {
		return err
	}

	if err := qc.conn.WriteMessage(websocket.TextMessage, subscribeMsgBz); err != nil {
		return err
	}

	go qc.listen(ctx, delegateNotifier)
	go func() {
		<-ctx.Done()
		log.Println("closing websocket")
		_ = qc.conn.Close()
	}()
	qc.delegateNotifee = delegateNotifee
	return nil
}
