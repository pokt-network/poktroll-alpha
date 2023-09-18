package client

import (
	"context"
	"fmt"

	cosmosClient "github.com/cosmos/cosmos-sdk/client"
	txClient "github.com/cosmos/cosmos-sdk/client/tx"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	authClient "github.com/cosmos/cosmos-sdk/x/auth/client"
	"github.com/gorilla/websocket"
	"github.com/pokt-network/smt"

	"poktroll/utils"
	"poktroll/x/servicer/types"
)

var (
	_ types.ServicerClient = &servicerClient{}
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
	keyName   string
	txFactory txClient.Factory
	clientCtx cosmosClient.Context
	wsClient  *websocket.Conn
	newBlocks utils.Observable[types.Block]
}

func NewServicerClient() *servicerClient {
	return &servicerClient{}
}

func (client *servicerClient) NewBlocks() utils.Observable[types.Block] {
	return client.newBlocks
}

func (client *servicerClient) SubmitClaim(
	ctx context.Context,
	smtRootHash []byte,
) error {
	msg := &types.MsgClaim{
		Creator:     client.clientCtx.FromAddress.String(),
		SmtRootHash: smtRootHash,
	}
	if err := client.broadcastMessageTx(ctx, msg); err != nil {
		return err
	}
	return nil
}

func (client *servicerClient) SubmitProof(
	ctx context.Context,
	smtRootHash []byte,
	closestKey []byte,
	closestValueHash []byte,
	closestSum uint64,
	// TODO: what type should `claim` be?
	proof *smt.SparseMerkleProof,
) error {
	proofBz, err := proof.Marshal()
	if err != nil {
		return err
	}

	msg := &types.MsgProof{
		Creator:   client.clientCtx.FromAddress.String(),
		Root:      smtRootHash,
		Path:      closestKey,
		ValueHash: closestValueHash,
		// CONSIDERATION: should we change this type in the protobuf?
		Sum:     int32(closestSum),
		ProofBz: proofBz,
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

func (client *servicerClient) listen(ctx context.Context, newBlocks chan types.Block) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		_, _, err := client.wsClient.ReadMessage()
		if err != nil {
			continue
		}

		newBlocks <- Block{
			height: 1,
			hash:   []byte(""),
		}
	}
}

func (client *servicerClient) WithSigningKeyUID(uid string) *servicerClient {
	client.keyName = uid
	return client
}

func (client *servicerClient) WithWsURL(ctx context.Context, wsURL string) *servicerClient {
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		panic(fmt.Errorf("failed to connect to websocket: %w", err))
	}

	conn.WriteJSON(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "subscribe",
		"id":      0,
		"params": map[string]interface{}{
			"query": "tm.event='NewBlock'",
		},
	})

	newBlocks, controller := utils.NewControlledObservable[types.Block](nil)

	client.wsClient = conn
	client.newBlocks = newBlocks

	go client.listen(ctx, controller)

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
