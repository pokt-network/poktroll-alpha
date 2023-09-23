package client

import (
	"context"
	"fmt"
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
	txHash, timeoutHeight, err := client.signAndBroadcastMessageTx(ctx, msg)
	if err != nil {
		return err
	}

	// TODO_THIS_COMMIT: factor out to a new method.
	client.txsMutex.Lock()
	if _, ok := client.txsByHashByTimeout[timeoutHeight]; !ok {
		// INCOMPLETE: handle and/or invalidate this case.
		panic(fmt.Errorf("txsByHash not found"))
	}

	txErrCh, ok := client.txsByHash[txHash]
	if !ok {
		// INCOMPLETE: handle and/or invalidate this case.
		panic(fmt.Errorf("txErrCh not found"))
	}
	client.txsMutex.Unlock()

	return <-txErrCh
}
