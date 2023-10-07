package cmd

import (
	"context"
	"fmt"
	ring_secp256k1 "github.com/athanorlabs/go-dleq/secp256k1"
	ring_types "github.com/athanorlabs/go-dleq/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"net/url"
	"os"
	"os/signal"
	"sync"

	cosmosclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/spf13/cobra"

	"poktroll/smartclient"
	client "poktroll/smartclient/client"
	"poktroll/smartclient/relayhandler"
	applicationTypes "poktroll/x/application/types"
	portalTypes "poktroll/x/portal/types"
	sessionTypes "poktroll/x/session/types"
)

const waitGroupContextKey = "smart_client_cmd_wait_group"

var (
	signingKeyName        string
	applicationListenHost string
	ringSinger            bool
)

func SmartClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "smartclient",
		Short: "A smart client command",
		RunE:  runSmartClient,
	}

	cmd.Flags().StringVar(&signingKeyName, "signing-key", "", "The signing key to use")
	cmd.Flags().StringVar(&applicationListenHost, "listen", "", "The application endpoint to use")
	cmd.Flags().BoolVar(&ringSinger, "ring-signer", false, "Use a ring signer instead of a simple signer")

	cmd.Flags().String(flags.FlagKeyringBackend, "", "Select keyring's backend (os|file|kwallet|pass|test)")
	cmd.Flags().String(flags.FlagNode, "tcp://localhost:36657", "tcp://<host>:<port> to tendermint rpc interface for this chain")

	return cmd
}

func runSmartClient(cmd *cobra.Command, args []string) error {
	// CONSIDERATION: there may be a more conventional, idiomatic, and/or convenient
	// way to track and cleanup goroutines. In the wait group solution, goroutines get a
	// reference to it via the context value and are expected to call `wg.Add(n)` and
	// `wg.Done()` appropriately.
	wg := new(sync.WaitGroup)
	ctx, cancelCtx := context.WithCancel(
		context.WithValue(
			cmd.Context(),
			waitGroupContextKey,
			wg,
		),
	)
	clientCtx := cosmosclient.GetClientContextFromCmd(cmd)
	clientFactory, err := tx.NewFactoryCLI(clientCtx, cmd.Flags())
	if err != nil {
		cancelCtx()
		panic(fmt.Errorf("failed to create tx factory: %w", err))
	}

	key, err := clientFactory.Keybase().Key(signingKeyName)
	if err != nil {
		cancelCtx()
		panic(fmt.Errorf("failed to get key: %w", err))
	}

	applicationAddress, err := key.GetAddress()
	if err != nil {
		cancelCtx()
		panic(fmt.Errorf("failed to get address: %w", err))
	}

	// use clientCtx.NodeURI to get the tendermint websocket endpoint
	blockQueryURL, err := url.Parse(clientCtx.NodeURI + "/websocket")
	if err != nil {
		cancelCtx()
		panic(fmt.Errorf("failed to parse block query URL: %w", err))
	}
	blockQueryURL.Scheme = "ws"

	// build the needed QueryClients (application, portal, session and account)
	applicationQueryClient := applicationTypes.NewQueryClient(clientCtx)
	portalQueryClient := portalTypes.NewQueryClient(clientCtx)
	sessionQueryClient := sessionTypes.NewQueryClient(clientCtx)
	accountQueryClient := authTypes.NewQueryClient(clientCtx)

	// build the new blocks subscription client
	blockQueryClient, err := client.NewBlocksQueryClient(ctx, blockQueryURL.String())
	if err != nil {
		cancelCtx()
		panic(fmt.Errorf("failed to create block query client: %w", err))
	}

	// use the ChooseFirstEndpoint strategy to select the first relayer endpoint
	endpointSelectionStrategy := &relayhandler.ChooseFirstEndpoint{}

	// create a signer from the keyring and signing key name
	// this should support a ring signature implementation
	// TODO: provide a flag to select the signer implementation
	var signer smartclient.Signer
	var signingKey ring_types.Scalar
	if !ringSinger {
		signer = smartclient.NewSimpleSigner(clientCtx.Keyring, signingKeyName)
	} else {
		signingKey, err = keyRecordToScalar(clientCtx.Keyring, signingKeyName)
		if err != nil {
			cancelCtx()
			panic(fmt.Errorf("failed to get signing key: %w", err))
		}
	}

	// ensure the protocol or any other part of the URL is not used in the listen address
	tcpListenAddr, err := url.Parse(applicationListenHost)
	if err != nil {
		cancelCtx()
		panic(fmt.Errorf("failed to parse application address: %w", err))
	}

	smartClient := relayhandler.NewRelayHandler(
		tcpListenAddr.Host,
		applicationQueryClient,
		portalQueryClient,
		sessionQueryClient,
		accountQueryClient,
		blockQueryClient,
		applicationAddress.String(),
		endpointSelectionStrategy,
		signer,
		signingKey,
	)
	if err != nil {
		cancelCtx()
		panic(fmt.Errorf("failed to create relay handler: %w", err))
	}

	if err := smartClient.Start(ctx); err != nil {
		cancelCtx()
		panic(fmt.Errorf("failed to start smart client: %w", err))
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

// keyRecordToScalar converts the private key obtained from a key record to a scalar
// point on the secp256k1 curve
func keyRecordToScalar(keyring keyring.Keyring, keyName string) (ring_types.Scalar, error) {
	keyRecord, err := keyring.Key(keyName)
	if err != nil {
		return nil, fmt.Errorf("key not found: %w", err)
	}
	local := keyRecord.GetLocal()
	if local == nil {
		return nil, fmt.Errorf("cannot extract private key from key record: nil")
	}
	priv, ok := local.PrivKey.GetCachedValue().(cryptotypes.PrivKey)
	if !ok {
		return nil, fmt.Errorf("cannot extract private key from key record: %T", local.PrivKey.GetCachedValue())
	}
	if _, ok := priv.(*secp256k1.PrivKey); !ok {
		return nil, fmt.Errorf("unexpected private key type: %T, want %T", priv, &secp256k1.PrivKey{})
	}
	crv := ring_secp256k1.NewCurve()
	privKey, err := crv.DecodeToScalar(priv.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %w", err)
	}
	return privKey, nil
}
