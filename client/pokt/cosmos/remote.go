package cosmos

import (
	"context"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/pokt-network/smt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"poktroll/app"
	"poktroll/modules"
	"poktroll/runtime/di"
	"poktroll/x/poktroll/types"
)

var (
	_ modules.PocketNetworkClient = &remoteCosmosPocketClient{}
)

type remoteCosmosPocketClient struct {
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

func (client *remoteCosmosPocketClient) Hydrate(injector *di.Injector, path *[]string) {}

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

func txResultFromTxResponse(txResp *cosmosTypes.TxResponse) *types.TxResult {
	return &types.TxResult{
		Hash:   txResp.TxHash,
		Height: uint64(txResp.Height),
	}
}

func newGRPCConn(
	ctx context.Context,
	target string,
) (*grpc.ClientConn, error) {
	return grpc.DialContext(
		ctx, target,
		// NB: this connection is necessary for the client to work, there's no
		// benefit to handling a dial error asynchronously.
		grpc.WithBlock(),
		// TODO_THIS_COMMIT: don't use insecure transport credentials
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}
