package client

import (
	"context"
	"encoding/json"
	"log"
	"poktroll/x/servicer/types"
)

func (client *servicerClient) SubmitClaim(
	ctx context.Context,
	smtRootHash []byte,
) error {
	smtRootHashStr := string(smtRootHash)
	log.Println("SANITY CHECK!")
	if client.address == "" {
		log.Println("client.address == ''")
		return errEmptyAddress
	}

	client.commitedClaimsMu.Lock()
	defer client.commitedClaimsMu.Unlock()
	if _, ok := client.committedClaims[smtRootHashStr]; ok {
		log.Println("existing commited claim channel")
		<-client.committedClaims[string(smtRootHash)]
		return nil
	} else {
		log.Println("no commited claim channel")
	}

	client.committedClaims[smtRootHashStr] = make(chan struct{})

	log.Println("broadcasting...")

	//queryClient := types.NewQueryClient(client.clientCtx)
	//_, err := queryClient.ServicersAll(ctx, &types.QueryAllServicersRequest{})
	//if err != nil {
	//	return err
	//}
	//log.Println("!!!QUERIED SERVICERS!!!")

	msg := &types.MsgClaim{
		Creator:     client.address,
		SmtRootHash: smtRootHash,
	}
	if err := client.broadcastMessageTx(ctx, msg); err != nil {
		return err
	}

	log.Println("..done broadcasting!")

	log.Println("reading from commited claim channel...")
	<-client.committedClaims[smtRootHashStr]
	log.Println(".. done reading from commited claim channel!")
	return nil
}

func (client *servicerClient) subscribeToClaims(ctx context.Context) {
	// TODO_THIS_COMMIT: query should be speific to THIS servicer.
	//query := fmt.Sprintf("tm.event='EventClaimed' ??.servicer_address='%s'", client.address)
	//query := "tm.event='EventClaimed'"
	//query := "message.module='servicer' AND message.action='Claim'"
	//query := "message.module='servicer'"
	//query := "tm.event='Tx' AND message.module='servicer'"
	query := "tm.event='Tx'"

	msgHandler := func(ctx context.Context, msg []byte) error {
		log.Println("event claimed!")
		log.Printf("%s", string(msg))
		var claim types.EventClaimed
		if err := json.Unmarshal(msg, &claim); err != nil {
			return err
		}
		log.Printf("claim: %+v", claim)

		client.commitedClaimsMu.Lock()
		defer client.commitedClaimsMu.Unlock()
		if claimCommittedCh, ok := client.committedClaims[string(claim.SmtRootHash)]; ok {
			log.Println("sending to claim commited channel")
			claimCommittedCh <- struct{}{}
		} else {
			log.Println("didn't find claim commited channel!")
		}
		return nil
	}
	client.subscribeWithQuery(ctx, query, msgHandler)
}
