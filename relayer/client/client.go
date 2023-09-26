package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	cosmosClient "github.com/cosmos/cosmos-sdk/client"
	txClient "github.com/cosmos/cosmos-sdk/client/tx"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	authClient "github.com/cosmos/cosmos-sdk/x/auth/client"

	"poktroll/utils"
	"poktroll/x/servicer/types"
)

var (
	_ ServicerClient = &servicerClient{}
	// TECHDEBT: consolidate into a(n) errors file(s)/package
	// errEmptyAddress is used when address hasn't been configured but is required.
	errEmptyAddress = fmt.Errorf("client address is empty")
)

type servicerClient struct {
	// nextRequestId is a *unique* ID intended to be monotonically incremented
	// and used to uniquely identify distinct RPC requests.
	nextRequestId uint64
	// address is the on-chain account address of this client (relayer / servicer).
	address string
	// txFactory is a cosmos-sdk tx factory which encapsulates everything
	// necessary to sign transactions given a client context.
	txFactory txClient.Factory
	// clientCtx is a cosmos-sdk client context which encapsulates everything
	// necessary to construct, encode, and broadcast transactions.
	clientCtx cosmosClient.Context

	blocksNotifee utils.Observable[types.Block]
	// TODO_CONSIDERATION: using an observable for received tx messages & a filter
	// for `#signAndBroadcastTx()` callers to react to the specific tx in question
	// instead of using shared memory across goroutines (`txByHash`) would likely
	// improve readability and maintainability. This would likely require a new
	// "buffered controllable observable"; i.e. a controlled observable which uses
	// buffered channels to avoid blocking channel sender.
	//
	//txsNotifee    utils.Observable[*cosmosTypes.TxResponse]

	// txsMutex protectx txsByHash and txsByHashByTimeout maps
	txsMutex sync.Mutex
	// txsByHash maps tx hash to a channel which will receive an error or nil,
	// and close, when the tx with the given hash is committed.
	txsByHash map[string]chan error
	// txsByHashByTimeout maps timeout (block) height to a map of txsByHash. It
	// is used to ensure that tx error channels receive and close in the event
	// that they have not already by the given timeout height.
	txsByHashByTimeout map[uint64]map[string]chan error

	// latestBlockMutex protext latestBlock.
	latestBlockMutex sync.RWMutex
	// latestBlock is the latest block that has been committed.
	latestBlock types.Block

	// Configuration
	// =============
	// keyName is the name of the key as per the CLI keyring/keybase.
	// See: `poktrolld keys list --help`.
	keyName string
	// wsURL is the URL of the websocket endpoint to connect to for RPC
	// service over websocket transport (with /subscribe support).
	wsURL string
	// INCOMPLETE: this should be configurable & integrated w/ viper, flags, etc.
	// txTimeoutHeightOffset is the number of blocks after the latest block
	// that a tx should be considered invalid if it has not been committed.
	txTimeoutHeightOffset uint32
}

func NewServicerClient() *servicerClient {
	return &servicerClient{
		txsByHash:          make(map[string]chan error),
		txsByHashByTimeout: make(map[uint64]map[string]chan error),
	}
}

func (client *servicerClient) signAndBroadcastMessageTx(
	ctx context.Context,
	msg cosmosTypes.Msg,
) (txErrCh chan error, err error) {
	// construct tx
	txConfig := client.clientCtx.TxConfig
	txBuilder := txConfig.NewTxBuilder()
	if err = txBuilder.SetMsgs(msg); err != nil {
		return nil, err
	}

	// calculate timeout height
	timeoutHeight := client.LatestBlock().Height() +
		uint64(client.txTimeoutHeightOffset)

	// TECHDEBT: this should be configurable
	txBuilder.SetGasLimit(200000)
	txBuilder.SetTimeoutHeight(timeoutHeight)

	// sign tx
	if err := authClient.SignTx(
		client.txFactory,
		client.clientCtx,
		client.keyName,
		txBuilder,
		false,
		false,
	); err != nil {
		return nil, err
	}

	// ensure tx is valid
	// NOTE: this makes the tx valid; i.e. it is *REQUIRED*
	if err := txBuilder.GetTx().ValidateBasic(); err != nil {
		return nil, err
	}

	// serialize tx
	txBz, err := client.encodeTx(txBuilder)
	if err != nil {
		return nil, err
	}

	txResponse, err := client.clientCtx.BroadcastTxSync(txBz)
	if err != nil {
		return nil, err
	}

	txResponseJSON, err := json.MarshalIndent(txResponse, "", "  ")
	if err != nil {
		panic(err)
	}

	// NB: the hex representation of the tx hash can has no canonical case but
	// must be consistent.
	txHash := strings.ToLower(txResponse.TxHash)
	log.Printf("txResponse: %s\n", txResponseJSON)

	return client.updateTxs(ctx, txHash, timeoutHeight)
}

func (client *servicerClient) updateTxs(
	ctx context.Context,
	txHash string,
	timeoutHeight uint64,
) (txErrCh chan error, err error) {
	client.txsMutex.Lock()
	defer client.txsMutex.Unlock()

	// Initialize txsByHashByTimeout map if necessary.
	txsByHash, ok := client.txsByHashByTimeout[timeoutHeight]
	if !ok {
		txsByHash = make(map[string]chan error)
		client.txsByHashByTimeout[timeoutHeight] = txsByHash
	}

	// Initialize txsByHash map in txsByHashByTimeout map if necessary.
	txErrCh, ok = txsByHash[txHash]
	if !ok {
		// NB: intentionally buffered to avoid blocking on send. Only intended
		// to send/receive a single error.
		txErrCh = make(chan error, 1)
		txsByHash[txHash] = txErrCh
	}
	// Initialize txsByHash map if necessary.
	if _, ok := client.txsByHash[txHash]; !ok {
		// NB: both maps hold a reference to the same channel so that we can check
		// if the channel has already been closed when timing out.
		client.txsByHash[txHash] = txErrCh
	}

	// TODO_THIS_COMMIT: check txResponse for error in logs, parse & send on
	// txErrCh if tx failed!!!
	return txErrCh, nil
}

func (client *servicerClient) encodeTx(txBuilder cosmosClient.TxBuilder) ([]byte, error) {
	return client.clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
}

// TODO_IMPROVE: Implement proper options for `servicerClient`
func (client *servicerClient) WithSigningKey(keyName string, address string) *servicerClient {
	client.keyName = keyName
	client.address = address

	return client
}

func (client *servicerClient) WithWsURL(ctx context.Context, wsURL string) *servicerClient {
	client.wsURL = wsURL
	client.blocksNotifee = client.subscribeToBlocks(ctx)
	client.subscribeToOwnTxs(ctx, client.blocksNotifee)
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

func (client *servicerClient) WithTxTimeoutHeightOffset(offset uint32) *servicerClient {
	client.txTimeoutHeightOffset = offset
	return client
}
