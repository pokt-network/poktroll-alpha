package client

import (
	"context"

	cosmosClient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	txClient "github.com/cosmos/cosmos-sdk/client/tx"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	authClient "github.com/cosmos/cosmos-sdk/x/auth/client"
	"github.com/pokt-network/smt"

	"poktroll/relayminer/types"
	"poktroll/x/servicer/client"
)

var (
	_ client.ServicerClient = &servicerClient{}
)

type servicerClient struct {
	keyName   string
	txFactory txClient.Factory
	clientCtx cosmosClient.Context
}

func NewServicerClient(
	keyName string,
	txFactory tx.Factory,
	clientCtx cosmosClient.Context,
) client.ServicerClient {
	return &servicerClient{
		keyName:   keyName,
		txFactory: txFactory,
		clientCtx: clientCtx,
	}
}

func (client *servicerClient) NewBlocks() <-chan types.Block {
	panic("implement me")
}

func (client *servicerClient) SubmitClaim(
	ctx context.Context,
	// TODO: what type should `claim` be?
	claim []byte,
) error {
	panic("implement me")
}

func (client *servicerClient) SubmitProof(
	ctx context.Context,
	closestKey []byte,
	closestValueHash []byte,
	closestSum uint64,
	// TODO: what type should `claim` be?
	proof *smt.SparseMerkleProof,
) error {
	//
	client.broadcastMessageTx(ctx, msg)
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
