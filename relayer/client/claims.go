package client

import (
	"context"
	"encoding/json"
	"fmt"
	"poktroll/x/servicer/types"
)

func (client *servicerClient) SubmitClaim(
	ctx context.Context,
	smtRootHash []byte,
) error {
	if client.address == "" {
		return errEmptyAddress
	}

	if _, ok := client.committedClaims[string(smtRootHash)]; ok {
		<-client.committedClaims[string(smtRootHash)]
		return nil
	}

	client.committedClaims[string(smtRootHash)] = make(chan struct{})

	msg := &types.MsgClaim{
		Creator:     client.address,
		SmtRootHash: smtRootHash,
	}
	if err := client.broadcastMessageTx(ctx, msg); err != nil {
		return err
	}

	<-client.committedClaims[string(smtRootHash)]
	return nil
}

func (client *servicerClient) subscribeToClaims(ctx context.Context) {
	query := fmt.Sprintf("message.module='servicer' AND message.action='claim' AND message.sender='%s'", client.address)

	msgHandler := func(ctx context.Context, msg []byte) error {
		var claim types.EventClaimed
		if err := json.Unmarshal(msg, &claim); err != nil {
			return err
		}
		if claimCommittedCh, ok := client.committedClaims[string(claim.SmtRootHash)]; ok {
			claimCommittedCh <- struct{}{}
		}
		return nil
	}
	client.subscribeWithQuery(ctx, query, msgHandler)
}
