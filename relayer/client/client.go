package client

import (
	"context"
	"encoding/json"
	"fmt"
	authClient "github.com/cosmos/cosmos-sdk/x/auth/client"
	"log"
	"sync"

	cosmosClient "github.com/cosmos/cosmos-sdk/client"
	txClient "github.com/cosmos/cosmos-sdk/client/tx"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
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
	log.Printf("signing key name: %s", client.keyName)
	_, err := client.txFactory.Keybase().Key(client.keyName)
	if err != nil {
		return fmt.Errorf("failed to get key with UID %q: %w", client.keyName, err)
	}

	if err := authClient.SignTx(
		client.txFactory,
		client.clientCtx,
		client.keyName,
		txBuilder,
		false,
		true,
	); err != nil {
		return err
	}

	txBuilder.SetGasLimit(200000)
	tx := txBuilder.GetTx()
	sigs, err := tx.GetSignaturesV2()
	for _, sig := range sigs {
		sigJSON, err := json.MarshalIndent(sig, "", "  ")
		if err != nil {
			return err
		}
		log.Printf(string(sigJSON))
	}

	// serialize tx
	txBz, err := txConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return err
	}

	//var txBz []byte
	txRes, err := client.clientCtx.BroadcastTxSync(txBz)
	if err != nil {
		return err
	}

	txResJSON, err := json.MarshalIndent(txRes, "", "  ")
	if err != nil {
		return err
	}
	log.Printf(string(txResJSON))

	log.Printf("broadcast tx w/ hash: %q", txRes.TxHash)

	return nil
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
