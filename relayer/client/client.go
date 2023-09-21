package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	cosmosClient "github.com/cosmos/cosmos-sdk/client"
	txClient "github.com/cosmos/cosmos-sdk/client/tx"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	authClient "github.com/cosmos/cosmos-sdk/x/auth/client"
	"github.com/gorilla/websocket"
	"github.com/pokt-network/smt"

	"poktroll/relayer"
	"poktroll/utils"
	"poktroll/x/servicer/types"
)

var (
	_               types.ServicerClient = &servicerClient{}
	errEmptyAddress                      = fmt.Errorf("client address is empty")
)

type Block struct {
	height uint64
	hash   []byte
}

func (b Block) Height() uint64 {
	return b.height
}

func (b Block) Hash() []byte {
	return b.hash
}

type servicerClient struct {
	keyName         string
	address         string
	txFactory       txClient.Factory
	clientCtx       cosmosClient.Context
	wsURL           string
	committedClaims map[string]chan struct{}
	nextRequestId   uint64
	blocksNotifee   utils.Observable[types.Block]
}

func NewServicerClient() *servicerClient {
	return &servicerClient{
		committedClaims: make(map[string]chan struct{}),
	}
}

func (client *servicerClient) Blocks() utils.Observable[types.Block] {
	return client.blocksNotifee
}

func (client *servicerClient) SubmitClaim(
	ctx context.Context,
	smtRootHash []byte,
) error {
	if client.address == "" {
		return errEmptyAddress
	}

	if _, ok := client.committedClaims[string(smtRootHash)]; ok {
		<-client.committedClaims[string(smtRootHash)]
		return nil
	}

	client.committedClaims[string(smtRootHash)] = make(chan struct{})

	msg := &types.MsgClaim{
		Creator:     client.address,
		SmtRootHash: smtRootHash,
	}
	if err := client.broadcastMessageTx(ctx, msg); err != nil {
		return err
	}

	<-client.committedClaims[string(smtRootHash)]
	return nil
}

func (client *servicerClient) SubmitProof(
	ctx context.Context,
	smtRootHash []byte,
	closestKey []byte,
	closestValueHash []byte,
	closestSum uint64,
	smtProof *smt.SparseMerkleProof,
) error {
	if client.address == "" {
		return errEmptyAddress
	}

	proofBz, err := smtProof.Marshal()
	if err != nil {
		return err
	}

	msg := &types.MsgProof{
		Creator:   client.address,
		Root:      smtRootHash,
		Path:      closestKey,
		ValueHash: closestValueHash,
		Sum:       closestSum,
		Proof:     proofBz,
	}
	if err := client.broadcastMessageTx(ctx, msg); err != nil {
		return err
	}
	return nil
}

func (client *servicerClient) broadcastMessageTx(
	ctx context.Context,
	msg cosmosTypes.Msg,
) error {
	// construct tx
	txConfig := client.clientCtx.TxConfig
	txBuilder := txConfig.NewTxBuilder()
	if err := txBuilder.SetMsgs(msg); err != nil {
		return err
	}

	// sign tx
	if err := authClient.SignTx(
		client.txFactory,
		client.clientCtx,
		client.keyName,
		txBuilder,
		false,
		false,
	); err != nil {
		return err
	}

	// serialize tx
	txBz, err := txConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return err
	}

	if _, err := client.clientCtx.BroadcastTxSync(txBz); err != nil {
		return err
	}

	return nil
}

// listen blocks on reading messages from a websocket connection, it is intended
// to be called from within a go routine.
func (client *servicerClient) listen(ctx context.Context, conn *websocket.Conn, msgHandler messageHandler) {
	wg, haveWaitGroup := ctx.Value(relayer.WaitGroupContextKey).(*sync.WaitGroup)
	if haveWaitGroup {
		// Increment the relayer wait group to track this goroutine
		wg.Add(1)
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("closing websocket")
			_ = conn.Close()
			if haveWaitGroup {
				// Decrement the wait group as this goroutine stops
				wg.Done()
			}
			return
		default:
		}

		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				// NB: stop this goroutine if the websocket connection is closed
				return
			}
			log.Printf("skipping due to websocket error: %s\n", err)
			// TODO: handle other errors (?)
			continue
		}

		if err := msgHandler(ctx, msg); err != nil {
			log.Printf("skipping due to message handler error: %s\n", err)
			continue
		}
	}
}

func (client *servicerClient) WithSigningKeyUID(uid string) *servicerClient {
	key, err := client.txFactory.Keybase().Key(uid)

	if err != nil {
		panic(fmt.Errorf("failed to get key with UID %q: %w", uid, err))
	}

	address, err := key.GetAddress()
	if err != nil {
		panic(fmt.Errorf("failed to get address for key with UID %q: %w", uid, err))
	}

	client.keyName = uid
	client.address = address.String()

	return client
}

func (client *servicerClient) WithWsURL(ctx context.Context, wsURL string) *servicerClient {
	client.wsURL = wsURL
	client.blocksNotifee = client.subscribeToBlocks(ctx)
	client.subscribeToClaims(ctx)
	return client
}

func (client *servicerClient) WithTxFactory(txFactory txClient.Factory) *servicerClient {
	client.txFactory = txFactory
	return client
}

func (client *servicerClient) WithClientCtx(clientCtx cosmosClient.Context) *servicerClient {
	client.clientCtx = clientCtx
	return client
}

type messageHandler func(ctx context.Context, msg []byte) error

func (client *servicerClient) subscribeWithQuery(ctx context.Context, query string, msgHandler messageHandler) {
	conn, _, err := websocket.DefaultDialer.Dial(client.wsURL, nil)
	if err != nil {
		panic(fmt.Errorf("failed to connect to websocket: %w", err))
	}

	requestId := client.getNextRequestId()
	conn.WriteJSON(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "subscribe",
		"id":      requestId,
		"params": map[string]interface{}{
			"query": query,
		},
	})

	go client.listen(ctx, conn, msgHandler)
}

func (client *servicerClient) subscribeToBlocks(ctx context.Context) utils.Observable[types.Block] {
	query := "tm.event='NewBlock'"

	blocksNotifee, blocksNotifier := utils.NewControlledObservable[types.Block](nil)
	msgHandler := handleBlocksFactory(blocksNotifier)
	client.subscribeWithQuery(ctx, query, msgHandler)

	return blocksNotifee
}

func (client *servicerClient) subscribeToClaims(ctx context.Context) {
	query := fmt.Sprintf("message.module='servicer' AND message.action='claim' AND message.sender='%s'", client.address)

	msgHandler := func(ctx context.Context, msg []byte) error {
		var claim types.EventClaimed
		if err := json.Unmarshal(msg, &claim); err != nil {
			return err
		}
		if claimCommittedCh, ok := client.committedClaims[string(claim.Root)]; ok {
			claimCommittedCh <- struct{}{}
		}
		return nil
	}
	client.subscribeWithQuery(ctx, query, msgHandler)
}

func (client *servicerClient) getNextRequestId() uint64 {
	client.nextRequestId++
	return client.nextRequestId
}

func handleBlocksFactory(blocksNotifier chan types.Block) messageHandler {
	return func(ctx context.Context, msg []byte) error {
		block, err := NewTendermintBlockEvent(msg)
		if err != nil {
			return fmt.Errorf("skipping due to new block event error: %w", err)
		}

		// If msg does not contain data then block is nil, we can ignore it
		if block == nil {
			return fmt.Errorf("skipping because block is nil")
		}

		log.Printf("new block; height: %d, hash: %x\n", block.Height(), block.Hash())
		blocksNotifier <- block

		return nil
	}
}
