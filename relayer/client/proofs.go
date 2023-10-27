package client

import (
	"context"

	"github.com/pokt-network/smt"

	"poktroll/x/servicer/types"
	sharedtypes "poktroll/x/shared/types"
)

func (client *servicerClient) SubmitProof(
	ctx context.Context,
	session *sharedtypes.Session,
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
		Session:         session,
		ServicerAddress: client.address,
		SmstRootHash:    smtRootHash,
		Path:            closestKey,
		ValueHash:       closestValueHash,
		SmstSum:         closestSum,
		Proof:           proofBz,
	}
	if _, err = client.signAndBroadcastMessageTx(ctx, msg); err != nil {
		return err
	}
	return nil
}
