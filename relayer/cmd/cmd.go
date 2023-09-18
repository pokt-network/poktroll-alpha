package cmd

import (
	cosmosclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"poktroll/relayer"
	"poktroll/relayer/client"
)

var signingKeyName string
var wsURL string

func RelayerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "relayer",
		Short: "Run a relayer",
		Long:  `Run a relayer`,
		RunE:  runRelayer,
	}

	cmd.Flags().StringVar(&signingKeyName, "signing-key", "", "Name of the key to sign transactions")
	cmd.Flags().StringVar(&wsURL, "ws-url", "ws://localhost:26657/websocket", "Websocket URL to poktrolld node")

	return cmd
}

func runRelayer(cmd *cobra.Command, args []string) error {
	// construct client
	clientCtx := cosmosclient.GetClientContextFromCmd(cmd)
	clientFactory, err := tx.NewFactoryCLI(clientCtx, cmd.Flags())
	if err != nil {
		return err
	}

	ctx := cmd.Context()

	c := client.NewServicerClient().
		WithSigningKeyUID(signingKeyName).
		WithTxFactory(clientFactory).
		WithClientCtx(clientCtx).
		WithWsURL(ctx, wsURL)

	relayer := relayer.NewRelayer(ctx, c)

	return relayer.Start()
}
