package client

import (
	"context"
	"github.com/pokt-network/smt"
	"poktroll/x/servicer/types"
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
		Creator:   client.address,
		Root:      smtRootHash,
		Path:      closestKey,
		ValueHash: closestValueHash,
		Sum:       closestSum,
		Proof:     proofBz,
	}
	_, _, err = client.signAndBroadcastMessageTx(ctx, msg)
	if err != nil {
		return err
	}
	return nil
}
