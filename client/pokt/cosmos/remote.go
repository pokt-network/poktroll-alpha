package cosmos

import (
	"context"
	"fmt"
	"poktroll/app"
	"reflect"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/pokt-network/smt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"poktroll/modules"
	"poktroll/runtime/di"
	"poktroll/x/poktroll/types"
)

var (
	_                      modules.PocketNetworkClient = &remoteCosmosPocketClient{}
	errInvalidActorTypeFmt                             = "invalid actor type: %s"
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

func (client *remoteCosmosPocketClient) Stake(
	ctx context.Context,
	actor *types.Actor,
	amount uint64,
) <-chan types.Maybe[*types.TxResult] {
	var (
		creator   string
		amount    uint64
		actorType string
		stakeInfo types.StakeInfo
		resultCh  = make(chan types.Maybe[*types.TxResult], 1)
	)

	// TODO_THIS_COMMIT: provide encoding config via DI (?)
	go func() {
		switch actor.GetActorType().(type) {
		case types.Actor_Servicer:
			stakeInfo = actor.GetServicer().GetStakeInfo()

			// TECHDEBT: align this once actorType is a go enum in the message handler
			servicerType := reflect.TypeOf(types.Actor_Servicer{})
			actorType = servicerType.Name()
		default:
			resultCh <- types.JustError[*types.TxResult](
				fmt.Errorf(errInvalidActorTypeFmt, actor.GetActorType()),
			)
		}
		msg := types.NewMsgStake(actor.GetServicer().GetSta)

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
	}()
	return resultCh
}

func (client *remoteCosmosPocketClient) Unstake(
	ctx context.Context,
	actor *types.Actor,
	amount uint64,
) <-chan types.Maybe[*types.TxResult] {
	panic("implement me")
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
