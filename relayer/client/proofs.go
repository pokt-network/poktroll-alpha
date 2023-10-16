package client

import (
	"context"
	"poktroll/x/servicer/types"

	"github.com/pokt-network/smt"
)

func (client *servicerClient) SubmitProof(
	ctx context.Context,
	sessionId string,
	smtRootHash []byte,
	smtProof *smt.SparseMerkleClosestProof,
) error {
	if client.address == "" {
		return errEmptyAddress
	}

	proofBz, err := smtProof.Marshal()
	if err != nil {
		return err
	}

	msg := &types.MsgProof{
		SessionId:       sessionId,
		ServicerAddress: client.address,
		SmstRootHash:    smtRootHash,
		Proof:           proofBz,
	}
	txErrCh, err := client.signAndBroadcastMessageTx(ctx, msg)
	if err != nil {
		return err
	}

	return <-txErrCh
}
