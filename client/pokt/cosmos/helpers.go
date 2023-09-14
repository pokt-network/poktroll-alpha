package cosmos

import (
	"context"
	types2 "poktroll/x/poktroll/types"

	"github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func txResultFromTxResponse(txResp *types.TxResponse) *types2.TxResult {
	return &types2.TxResult{
		Hash:   txResp.TxHash,
		Height: txResp.Height,
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