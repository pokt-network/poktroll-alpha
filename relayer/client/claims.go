package client

import (
	"context"
	"poktroll/x/servicer/types"
)

// SubmitClaim implements the respective method on the ServicerClient interface.
func (client *servicerClient) SubmitClaim(
	ctx context.Context,
	smtRootHash []byte,
) error {
	if client.address == "" {
		return errEmptyAddress
	}

	msg := &types.MsgClaim{
		Creator:     client.address,
		SmtRootHash: smtRootHash,
	}
	txErrCh, err := client.signAndBroadcastMessageTx(ctx, msg)
	if err != nil {
		return err
	}

	return <-txErrCh
}
