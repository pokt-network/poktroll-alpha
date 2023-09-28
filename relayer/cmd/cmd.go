package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"poktroll/relayer/client"
	"sync"

	cosmosclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"poktroll/relayer"
)

var signingKeyName string
var wsURL string
var smtStorePath string

func RelayerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "relayer",
		Short: "Run a relayer",
		Long:  `Run a relayer`,
		RunE:  runRelayer,
	}

	// TECHDEBT: integrate these flags with the client context (i.e. flags, config, viper, etc.)
	// This is simpler to do with server-side configs (see rootCmd#PersistentPreRunE).
	// Will require more effort than currently justifiable.
	cmd.Flags().StringVar(&signingKeyName, "signing-key", "", "Name of the key to sign transactions")
	cmd.Flags().StringVar(&wsURL, "ws-url", "ws://localhost:36657/websocket", "Websocket URL to poktrolld node; formatted as ws://<host>:<port>[/path]")
	cmd.Flags().StringVar(&smtStorePath, "smt-store", "smt", "Path to the SMT KV store")

	cmd.Flags().String(flags.FlagKeyringBackend, "", "Select keyring's backend (os|file|kwallet|pass|test)")
	cmd.Flags().String(flags.FlagNode, "tcp://localhost:36657", "tcp://<host>:<port> to tendermint rpc interface for this chain")

	return cmd
}

func runRelayer(cmd *cobra.Command, _ []string) error {
	clientCtx := cosmosclient.GetClientContextFromCmd(cmd)
	clientFactory, err := tx.NewFactoryCLI(clientCtx, cmd.Flags())
	if err != nil {
		return err
	}

	// CONSIDERATION: there may be a more conventional, idiomatic, and/or convenient
	// way to track and cleanup goroutines. In the wait group solution, goroutines get a
	// reference to it via the context value and are expected to call `wg.Add(n)` and
	// `wg.Done()` appropriately.
	wg := new(sync.WaitGroup)
	ctx, cancelCtx := context.WithCancel(
		context.WithValue(
			cmd.Context(),
			client.WaitGroupContextKey,
			wg,
		),
	)

	// TODO_REFACTOR: Factor out the key retrieval and address extraction.
	key, err := clientFactory.Keybase().Key(signingKeyName)
	if err != nil {
		panic(fmt.Errorf("failed to get key with UID %q: %w", signingKeyName, err))
	}
	address, err := key.GetAddress()
	if err != nil {
		panic(fmt.Errorf("failed to get address for key with UID %q: %w", signingKeyName, err))
	}

	servicerClient := client.NewServicerClient().
		WithTxFactory(clientFactory).
		WithSigningKey(signingKeyName, address.String()).
		WithClientCtx(clientCtx).
		WithWsURL(ctx, wsURL).
		// TECHDEBT: this should be a config field.
		WithTxTimeoutHeightOffset(5)

	// INCOMPLETE: this should be populated from some relayer config.
	serviceEndpoints := map[string][]string{
		"svc1": {"ws://anvil:8547/"},
		"svc2": {"http://anvil:8547"},
	}

	// The order of the WithXXX methods matters for now.
	// TODO: Refactor this to a builder pattern.
	relayer := relayer.NewRelayer().
		WithKey(ctx, clientFactory.Keybase(), signingKeyName, address.String(), clientCtx, servicerClient, serviceEndpoints).
		WithServicerClient(servicerClient).
		WithKVStorePath(ctx, filepath.Join(clientCtx.HomeDir, smtStorePath))

	if err := relayer.Start(); err != nil {
		cancelCtx()
		return err
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	// Block until we receive an interrupt or kill signal (OS-agnostic)
	<-sigCh

	// Signal goroutines to stop
	cancelCtx()
	// Wait for all goroutines to finish
	wg.Wait()

	return nil
}
