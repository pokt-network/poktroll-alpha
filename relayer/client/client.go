package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	cometTypes "github.com/cometbft/cometbft/types"
	cosmosClient "github.com/cosmos/cosmos-sdk/client"
	txClient "github.com/cosmos/cosmos-sdk/client/tx"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	authClient "github.com/cosmos/cosmos-sdk/x/auth/client"

	"poktroll/utils"
	"poktroll/x/servicer/types"
)

var (
	_               types.ServicerClient = &servicerClient{}
	errEmptyAddress                      = fmt.Errorf("client address is empty")
)

type servicerClient struct {
	nextRequestId uint64
	address       string
	txFactory     txClient.Factory
	clientCtx     cosmosClient.Context

	blocksNotifee utils.Observable[types.Block]
	txsNotifee    utils.Observable[*cosmosTypes.TxResponse]

	txsMutex           sync.Mutex
	txsByHash          map[string]chan error
	txsByHashByTimeout map[uint64]map[string]chan error

	latestBlockMutex sync.RWMutex
	latestBlock      types.Block

	// Configuration
	keyName               string
	wsURL                 string
	txTimeoutHeightOffset uint32
}

func NewServicerClient() *servicerClient {
	return &servicerClient{
		latestBlock: &cometBlockWebsocketMsg{Block: cometTypes.Block{
			Header: cometTypes.Header{
				Height: 5000000,
			},
		}},
		txsByHash:          make(map[string]chan error),
		txsByHashByTimeout: make(map[uint64]map[string]chan error),
	}
}

func (client *servicerClient) signAndBroadcastMessageTx(
	ctx context.Context,
	msg cosmosTypes.Msg,
) (txHash string, timeoutHeight uint64, err error) {
	// construct tx
	txConfig := client.clientCtx.TxConfig
	txBuilder := txConfig.NewTxBuilder()
	if err = txBuilder.SetMsgs(msg); err != nil {
		return "", 0, err
	}

	// calculate timeout height
	timeoutHeight = client.LatestBlock().Height() +
		uint64(client.txTimeoutHeightOffset)

	client.txsMutex.Lock()
	defer client.txsMutex.Unlock()

	txsByHash, ok := client.txsByHashByTimeout[timeoutHeight]
	if !ok {
		txsByHash = make(map[string]chan error)
		client.txsByHashByTimeout[timeoutHeight] = txsByHash
	}

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
		return "", 0, err
	}

	// ensure tx is valid
	// NOTE: this makes the tx valid; i.e. it is *REQUIRED*
	if err := txBuilder.GetTx().ValidateBasic(); err != nil {
		return "", 0, err
	}

	// serialize tx
	txBz, err := client.encodeTx(txBuilder)
	if err != nil {
		return "", 0, err
	}

	txResponse, err := client.clientCtx.BroadcastTxSync(txBz)
	if err != nil {
		return "", 0, err
	}

	txResponseJSON, err := json.MarshalIndent(txResponse, "", "  ")
	if err != nil {
		panic(err)
	}

	txHash = strings.ToLower(txResponse.TxHash)
	newTxErrCh := make(chan error, 1)
	txErrCh, ok := txsByHash[txHash]
	if !ok {
		txErrCh = newTxErrCh
		txsByHash[txHash] = txErrCh
	}
	if _, ok := client.txsByHash[txHash]; !ok {
		client.txsByHash[txHash] = txErrCh
	}

	// TODO_THIS_COMMIT: check txResponse for error in logs, parse & send on
	// txErrCh if tx failed!!!
	log.Printf("txResponse: %s\n", txResponseJSON)

	return txHash, timeoutHeight, nil
}

func (client *servicerClient) encodeTx(txBuilder cosmosClient.TxBuilder) ([]byte, error) {
	return client.clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
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
