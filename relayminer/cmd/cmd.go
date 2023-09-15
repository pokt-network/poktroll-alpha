package cmd

import (
	cosmosclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"poktroll/relayminer"
	"poktroll/relayminer/client"
)

var signingKeyName string

func RelayMinerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "relay-miner",
		// Collides with poktrolld subcommand namespace of same name
		// Aliases: []string{"sevicer"}
		Short: "Run a relay miner",
		Long:  `Run a relay miner`,
		RunE:  runRelayMiner,
	}

	cmd.Flags().StringVar(&signingKeyName, "signing-key", "", "Name of the key to sign transactions")

	return cmd
}

func runRelayMiner(cmd *cobra.Command, args []string) error {
	// construct client
	clientCtx := cosmosclient.GetClientContextFromCmd(cmd)
	clientFactory, err := tx.NewFactoryCLI(clientCtx, cmd.Flags())
	if err != nil {
		return err
	}

	c := client.NewServicerClient(signingKeyName, clientFactory, clientCtx)

	relayMiner := relayminer.NewRelayMiner(c)

	return relayMiner.Start()
}
