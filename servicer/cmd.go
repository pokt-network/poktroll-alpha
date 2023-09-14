package servicer

import (
	"context"
	"crypto/sha256"
	"fmt"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/pokt-network/smt"
	"github.com/spf13/cobra"

	"poktroll/client/pokt/cosmos"
	"poktroll/logger"
	"poktroll/miner"
	"poktroll/modules"
	"poktroll/relayer"
	"poktroll/runtime/di"
	"poktroll/servicer/config"
	sessionmanager "poktroll/session-manager"
	"poktroll/x/poktroll/types"
)

func GetServicerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "servicer",
		Short:                      fmt.Sprintf("Servicer commands for the %s module", types.ModuleName),
		SuggestionsMinimumDistance: 2,
		RunE:                       servicerCmd,
	}

	cmd.PersistentFlags().String(flags.FlagFrom, "", "Name or address of private key with which to sign")
	flags.AddKeyringFlags(cmd.PersistentFlags())

	return cmd
}

func servicerCmd(cmd *cobra.Command, args []string) error {
	injectorCtxValue := cmd.Context().Value(config.PoktrollDepInjectorContextKey)
	injector, ok := injectorCtxValue.(*di.Injector)
	if !ok {
		// TECHDEBT: return more useful errors.
		return fmt.Errorf("Invalid injcetor type %T", injectorCtxValue)
	}

	ctx := context.WithValue(cmd.Context(), config.PoktrollDepInjectorContextKey, injector)
	cmd.SetContext(ctx)

	clientCtx := client.GetClientContextFromCmd(cmd)
	factory, err := tx.NewFactoryCLI(clientCtx, cmd.Flags())
	if err != nil {
		// TECHDEBT: return more useful errors.
		return err
	}

	// TECHDEBT: we probably "should" be updating the client context instead of
	// reaching for the commands flag directly.
	fromFlagStr, err := cmd.Flags().GetString(flags.FlagFrom)
	if err != nil {
		// TECHDEBT: return more useful errors.
		return err
	}

	// NB: while we don't need to inject the key itself (just the name),
	// we should ensure that a key with the given name exists, otherwise
	// return the error.
	_, err = factory.Keybase().Key(fromFlagStr)
	if err != nil {
		// TECHDEBT: return more useful errors.
		return err
	}

	di.Provide(modules.ServicerToken, NewServicerModule(), injector)

	// PocketNetworkClient deps
	di.Provide(modules.ClientCtxInjectionToken, clientCtx, injector)
	di.Provide(modules.TxFactoryInjectionToken, factory, injector)
	di.Provide(modules.PocketNetworkClientToken, cosmos.NewLocalCosmosPocketClient(), injector)
	di.Provide(modules.KeyNameInjectionToken, fromFlagStr, injector)

	// Servicer's module deps
	di.Provide(modules.RelayerToken, relayer.NewRelayerModule(), injector)
	di.Provide(modules.MinerModuleToken, miner.NewMinerModule(), injector)
	di.Provide(modules.SessionManagerToken, sessionmanager.NewSessionManager(), injector)

	// Sparse Merkel Tree deps
	smtStore, err := getSMTStore(cmd)
	if err != nil {
		// TECHDEBT: return more useful errors.
		return err
	}
	di.Provide(miner.SMTStoreToken, smtStore, injector)
	di.Provide(miner.SMTHasherToken, sha256.New(), injector)

	// Global logger
	serverCtx := server.GetServerContextFromCmd(cmd)
	di.Provide(logger.CosmosLoggerToken, logger.NewLogger(serverCtx.Logger), injector)

	srvcr := di.HydrateMain(modules.ServicerToken, injector)
	if err := srvcr.Start(); err != nil {
		return err
	}

	return nil
}

func getSMTStore(cmd *cobra.Command) (smt.KVStore, error) {
	smtStorePath, err := getSMTStorePath(cmd)
	if err != nil {
		// TECHDEBT: return more useful errors.
		return nil, err
	}
	return smt.NewKVStore(smtStorePath)
}
func getSMTStorePath(cmd *cobra.Command) (string, error) {
	homeDirFlagStr, err := cmd.Flags().GetString(flags.FlagHome)
	if err != nil {
		// TECHDEBT: return more useful errors.
		return "", err
	}
	return filepath.Join(homeDirFlagStr, "smt.db"), nil
}