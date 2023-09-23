package client

import (
	"context"
	"fmt"
	"sync"

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
	keyName          string
	address          string
	txFactory        txClient.Factory
	clientCtx        cosmosClient.Context
	wsURL            string
	nextRequestId    uint64
	blocksNotifee    utils.Observable[types.Block]
	latestBlock      types.Block
	commitedClaimsMu sync.Mutex
	committedClaims  map[string]chan struct{}
}

func NewServicerClient() *servicerClient {
	return &servicerClient{
		committedClaims: make(map[string]chan struct{}),
	}
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

// TODO_IMPROVE: Implement proper options for `servicerClient`
func (client *servicerClient) WithSigningKey(keyName string, address string) *servicerClient {
	client.keyName = keyName
	client.address = address

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
