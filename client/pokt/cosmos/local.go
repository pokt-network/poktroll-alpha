package cosmos

import (
	"context"

	txClient "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/pokt-network/smt"

	"poktroll/app"
	"poktroll/modules"
	"poktroll/runtime/di"
	"poktroll/x/poktroll/types"
)

var (
	_ modules.PocketNetworkClient = &localCosmosPocketClient{}
)

type localCosmosPocketClient struct {
	privateKey *secp256k1.PrivKey
	//txFactory txClient.Factory
	//clientCtx cosmosClient.Context
}

func NewLocalCosmosPocketClient(ctx context.Context) (modules.PocketNetworkClient, error) {
	return &localCosmosPocketClient{}, nil
}

func (client *localCosmosPocketClient) Hydrate(injector *di.Injector, path *[]string) {
	client.privateKey = di.Hydrate(modules.PrivateKeyInjectionToken, injector, path)
	//client.txFactory = di.Hydrate(modules.TxFactoryInjectionToken, injector, path)
	//client.clientCtx = di.Hydrate(modules.ClientCtxInjectionToken, injector, path)
}

func (client *localCosmosPocketClient) CascadeStart() error {
	return nil
}

func (client *localCosmosPocketClient) Start() error {
	// CONSIDERATION: could move grpc dialing to here instead of the constructor.
	return nil
}

func (client *localCosmosPocketClient) Stop() {
	//_ = client.grpcConn.Close()
}

func (client *localCosmosPocketClient) StakeServicer(
	ctx context.Context,
	servicer *types.Servicer,
	amount string,
) <-chan types.Maybe[*types.TxResult] {
	var (
		resultCh = make(chan types.Maybe[*types.TxResult], 1)
	)

	msg := types.NewMsgStake(
		servicer.StakeInfo.GetAddress(),
		amount,
		// TECHDEBT: update once `poktroll.Keeper#StakeActor()` is refactored.
		types.ServicerPrefix,
	)

	go client.broadcastMessageTx(ctx, resultCh, msg)
	return resultCh
}

func (client *localCosmosPocketClient) StakeApplication(
	ctx context.Context,
	application *types.Application,
	amount string,
) <-chan types.Maybe[*types.TxResult] {
	resultCh := make(chan types.Maybe[*types.TxResult], 1)

	// TODO_THIS_COMMIT: provide encoding config via DI (?)
	msg := types.NewMsgStake(
		application.StakeInfo.GetAddress(),
		amount,
		// TECHDEBT: update once `poktroll.Keeper#StakeActor()` is refactored.
		types.ServicerPrefix,
	)

	go client.broadcastMessageTx(ctx, resultCh, msg)
	return resultCh
}

func (client *localCosmosPocketClient) UnstakeServicer(
	ctx context.Context,
	servicer *types.Servicer,
	amount string,
) <-chan types.Maybe[*types.TxResult] {
	resultCh := make(chan types.Maybe[*types.TxResult], 1)

	msg := types.NewMsgUnstake(
		servicer.StakeInfo.GetAddress(),
		amount,
		// TECHDEBT: update once `poktroll.Keeper#StakeActor()` is refactored.
		types.ServicerPrefix,
	)

	go client.broadcastMessageTx(ctx, resultCh, msg)
	return resultCh
}
func (client *localCosmosPocketClient) UnstakeApplication(
	ctx context.Context,
	application *types.Application,
	amount string,
) <-chan types.Maybe[*types.TxResult] {
	resultCh := make(chan types.Maybe[*types.TxResult], 1)

	msg := types.NewMsgUnstake(
		application.StakeInfo.GetAddress(),
		amount,
		// TECHDEBT: update once `poktroll.Keeper#StakeActor()` is refactored.
		types.ServicerPrefix,
	)

	go client.broadcastMessageTx(ctx, resultCh, msg)
	return resultCh
}

func (client *localCosmosPocketClient) NewBlocks() <-chan *types.Block {
	panic("implement me")
}

func (client *localCosmosPocketClient) SubmitClaim(
	ctx context.Context,
	// TODO: what type should `claim` be?
	claim []byte,
) <-chan types.Maybe[*types.TxResult] {
	panic("implement me")
}

func (client *localCosmosPocketClient) SubmitProof(
	ctx context.Context,
	closestKey []byte,
	closestValueHash []byte,
	closestSum uint64,
	// TODO: what type should `claim` be?
	proof *smt.SparseMerkleProof,
	err error,
) <-chan types.Maybe[*types.TxResult] {
	panic("implement me")
}

func (client *localCosmosPocketClient) broadcastMessageTx(
	ctx context.Context,
	resultCh chan<- types.Maybe[*types.TxResult],
	msg cosmosTypes.Msg,
) {
	//client.clientCtx.BroadcastMode
	//if err := txClient.GenerateOrBroadcastTxWithFactory(
	//	client.clientCtx,
	//	client.txFactory,
	//	msg,
	//); err != nil {
	//	resultCh <- types.JustError[*types.TxResult](err)
	//}

	// TODO_THIS_COMMIT: use DI to get updated client context!
	txConfig := app.MakeEncodingConfig().TxConfig
	txBuilder := txConfig.NewTxBuilder()
	if err := txBuilder.SetMsgs(msg); err != nil {
		resultCh <- types.JustError[*types.TxResult](err)
	}

	txBuilder.SetSignatures(signing.SignatureV2{
		PubKey: client.privateKey.PubKey(),
		Data: &signing.SingleSignatureData{
			// TODO: what is this / how does it work?
			SignMode:  txConfig.SignModeHandler().DefaultMode(),
			Signature: nil,
		},
		Sequence:
	})

	txClient.SignWithPrivKey()

	resultCh <- types.Just(txResultFromTxResponse(bcastTxResp.GetTxResponse()))
}
