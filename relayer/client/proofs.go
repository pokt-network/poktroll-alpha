package client

import (
	"context"
	"poktroll/x/servicer/types"

	"github.com/pokt-network/smt"
)

func (client *servicerClient) SubmitProof(
	ctx context.Context,
	smtRootHash []byte,
	closestKey []byte,
	closestValueHash []byte,
	closestSum uint64,
	smtProof *smt.SparseMerkleProof,
) error {
	if client.address == "" {
		return errEmptyAddress
	}

	proofBz, err := smtProof.Marshal()
	if err != nil {
		return err
	}

	msg := &types.MsgProof{
		Servicer:     client.address,
		SmstRootHash: smtRootHash,
		Path:         closestKey,
		ValueHash:    closestValueHash,
		Sum:          closestSum,
		Proof:        proofBz,
	}
	_, err = client.signAndBroadcastMessageTx(ctx, msg)
	if err != nil {
		return err
	}
	return nil
}
