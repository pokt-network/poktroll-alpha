package cmd

import (
	"context"
	"os"
	"os/signal"
	"sync"

	cosmosclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"poktroll/relayer"
	"poktroll/relayer/client"
)

var signingKeyName string
var wsURL string
var blocksPerSession uint32
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
	cmd.Flags().Uint32VarP(&blocksPerSession, "blocks-per-session", "b", 2, "Websocket URL to poktrolld node")
	cmd.Flags().StringVar(&smtStorePath, "smt-store", "", "Path to the SMT KV store")

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

	// CONSIDERATION: there may be a more conventional, idomatic, and/or convenient
	// way to track and cleanup goroutines. In the wait group solution, goroutines get a
	// reference to it via the context value and are expected to call `wg.Add(n)` and 
	// `wg.Done()` appropriately.
	wg := new(sync.WaitGroup)
	ctx, cancelCtx := context.WithCancel(
		context.WithValue(
			cmd.Context(),
			relayer.WaitGroupContextKey,
			wg,
		),
	)

	// IMPROVE: we tried this pattern because it seemed to be conventional across
	// some cosmos-sdk code. In our use case, it turned out to be problematic. In
	// the presence of shared and/or nested dependencies, call order starts to
	// matter.
	// CONSIDERATION: perhaps the `depinject` cosmos-sdk system or a builder
	// pattern would be more appropriate.
	// see: https://github.com/cosmos/cosmos-sdk/tree/main/depinject#depinject
	func (relayer *Relayer) WithKVStorePath(storePath string) *Relayer {
	c := client.NewServicerClient().
		WithTxFactory(clientFactory).
		WithSigningKeyUID(signingKeyName).
		WithClientCtx(clientCtx).
		WithWsURL(ctx, wsURL)

	// The order of the WithXXX methods matters for now.
	// TODO: Refactor this to a builder pattern.
	relayer := relayer.NewRelayer().
		WithKey(clientFactory.Keybase(), signingKeyName).
		WithServicerClient(c).
		WithBlocksPerSession(ctx, blocksPerSession).
		WithKVStorePath(smtStorePath)

	if err := relayer.Start(); err != nil {
		cancelCtx()
		return err
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, os.Kill)
	// Block until we receive an interrupt or kill signal (OS-agnostic)
	<-sigCh

	// Signal goroutines to stop
	cancelCtx()
	// Wait for all goroutines to finish
	wg.Wait()

	return nil
}
