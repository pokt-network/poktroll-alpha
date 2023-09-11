package cosmos

import (
	"context"
	txClient "github.com/cosmos/cosmos-sdk/client/tx"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/pokt-network/smt"
	"google.golang.org/grpc"

	"poktroll/app"
	"poktroll/modules"
	"poktroll/runtime/di"
	"poktroll/x/poktroll/types"
)

var (
	_ modules.PocketNetworkClient = &remoteCosmosPocketClient{}
)

type remoteCosmosPocketClient struct {
	//privateKey  *secp256k1.PrivKey
	txFactory   txClient.Factory
	grpcConn    *grpc.ClientConn
	txClient    tx.ServiceClient
	queryClient types.QueryClient
}

func NewRemoteCosmosPocketClient(
	ctx context.Context,
	grpcURI string,
) (modules.PocketNetworkClient, error) {
	grpcConn, err := newGRPCConn(ctx, grpcURI)
	if err != nil {
		return nil, err
	}

	return &remoteCosmosPocketClient{
		txClient:    tx.NewServiceClient(grpcConn),
		queryClient: types.NewQueryClient(grpcConn),
	}, nil
}

func (client *remoteCosmosPocketClient) Hydrate(injector *di.Injector, path *[]string) {
	//client.privateKey = di.Hydrate(modules.PrivateKeyInjectionToken, injector, path)
	client.txFactory = di.Hydrate(modules.TxFactoryInjectionToken, injector, path)
}

func (client *remoteCosmosPocketClient) CascadeStart() error {
	return nil
}

func (client *remoteCosmosPocketClient) Start() error {
	// CONSIDERATION: could move grpc dialing to here instead of the constructor.
	return nil
}

func (client *remoteCosmosPocketClient) Stop() {
	_ = client.grpcConn.Close()
}

func (client *remoteCosmosPocketClient) StakeServicer(
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

func (client *remoteCosmosPocketClient) StakeApplication(
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

func (client *remoteCosmosPocketClient) UnstakeServicer(
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
func (client *remoteCosmosPocketClient) UnstakeApplication(
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

func (client *remoteCosmosPocketClient) NewBlocks() <-chan *types.Block {
	panic("implement me")
}

func (client *remoteCosmosPocketClient) SubmitClaim(
	ctx context.Context,
	// TODO: what type should `claim` be?
	claim []byte,
) <-chan types.Maybe[*types.TxResult] {
	panic("implement me")
}

func (client *remoteCosmosPocketClient) SubmitProof(
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

func (client *remoteCosmosPocketClient) broadcastMessageTx(
	ctx context.Context,
	resultCh chan<- types.Maybe[*types.TxResult],
	msg cosmosTypes.Msg,
) {

	// CONSIDERATION: provide encoding config via DI instead of importing the
	// cosmos app here (?)
	txConfig := app.MakeEncodingConfig().TxConfig
	txBuilder := txConfig.NewTxBuilder()
	if err := txBuilder.SetMsgs(msg); err != nil {
		resultCh <- types.JustError[*types.TxResult](err)
	}

	txBz, err := txConfig.TxEncoder()(txBuilder.GetTx())
	bcastTxResp, err := client.txClient.BroadcastTx(ctx, &tx.BroadcastTxRequest{
		TxBytes: txBz,
		Mode:    tx.BroadcastMode_BROADCAST_MODE_SYNC,
	})

	if err != nil {
		resultCh <- types.JustError[*types.TxResult](err)
	}

	resultCh <- types.Just(txResultFromTxResponse(bcastTxResp.GetTxResponse()))
}
