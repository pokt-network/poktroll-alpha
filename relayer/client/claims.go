package client

import (
	"context"
	"poktroll/x/servicer/types"
)

// SubmitClaim implements the respective method on the ServicerClient interface.
func (client *servicerClient) SubmitClaim(
	ctx context.Context,
	// TODO_REFACTOR: we should be passing sessionHeader everywhere instead of sessionId
	sessionId string,
	smtRootHash []byte,
) error {
	if client.address == "" {
		return errEmptyAddress
	}

	msg := &types.MsgClaim{
		ServicerAddress: client.address,
		SmstRootHash:    smtRootHash,
		SessionId:       sessionId,
	}
	txErrCh, err := client.signAndBroadcastMessageTx(ctx, msg)
	if err != nil {
		return err
	}

	return <-txErrCh
}
