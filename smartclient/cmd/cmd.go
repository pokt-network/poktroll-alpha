package cmd

import (
	"fmt"

	cosmosclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"poktroll/smartclient"
	client "poktroll/smartclient/client"
	"poktroll/smartclient/relayhandler"
	applicationTypes "poktroll/x/application/types"
	sessionTypes "poktroll/x/session/types"
)

var (
	signingKeyName  string
	relayerEndpoint string
)

var SmartClientCmd = &cobra.Command{
	Use:   "smartclient",
	Short: "A smart client command",
	RunE:  runSmartClient,
}

func init() {
	SmartClientCmd.Flags().StringVar(&signingKeyName, "signing-key", "", "The signing key to use")
	SmartClientCmd.Flags().StringVar(&relayerEndpoint, "relayer-endpoint", "", "The relayer endpoint to use")
	SmartClientCmd.MarkFlagRequired("signing-key")
	SmartClientCmd.MarkFlagRequired("relayer-endpoint")
}

func runSmartClient(cmd *cobra.Command, args []string) error {
	clientCtx := cosmosclient.GetClientContextFromCmd(cmd)
	//clientFactory, err := tx.NewFactoryCLI(clientCtx, cmd.Flags())
	//if err != nil {
	//	return err
	//}

	ctx := cmd.Context()

	key, err := clientCtx.Keyring.Key(signingKeyName)
	if err != nil {
		panic(fmt.Errorf("failed to get key: %w", err))
	}

	applicationAddress, err := key.GetAddress()
	if err != nil {
		panic(fmt.Errorf("failed to get address: %w", err))
	}

	applicationQueryClient := applicationTypes.NewQueryClient(clientCtx)
	sessionQueryClient := sessionTypes.NewQueryClient(clientCtx)
	blockQueryClient, err := client.NewBlocksQueryClient(ctx, relayerEndpoint)
	endpointSelectionStrategy := &relayhandler.ChooseFirstEndpoint{}
	signer := smartclient.NewSimpleSigner(clientCtx.Keyring, signingKeyName)

	smartClient := relayhandler.NewRelayHandler(
		relayerEndpoint,
		applicationQueryClient,
		sessionQueryClient,
		blockQueryClient,
		applicationAddress.String(),
		endpointSelectionStrategy,
		signer,
	)
	if err != nil {
		return err
	}

	smartClient.Start(ctx)

	return nil
}
