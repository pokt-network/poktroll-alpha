package client

import (
	"context"

	cosmosClient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	txClient "github.com/cosmos/cosmos-sdk/client/tx"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	authClient "github.com/cosmos/cosmos-sdk/x/auth/client"
	"github.com/pokt-network/smt"

	relayminer "poktroll/relayminer/types"
	"poktroll/x/servicer/client"
	"poktroll/x/servicer/types"
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

func (client *servicerClient) NewBlocks() <-chan relayminer.Block {
	panic("implement me")
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
