package servicer

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"poktroll/client/pokt/cosmos"
	"poktroll/modules"
	"poktroll/runtime/di"
	"poktroll/x/poktroll/types"
)

func GetServicerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Servicer commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       servicerCmd,
	}

	return cmd
}

func servicerCmd(cmd *cobra.Command, args []string) error {
	injector := di.NewInjector()
	ctx := context.WithValue(cmd.Context(), config.PoktrollDepInjectorContextKey, injector)
	cmd.SetContext(ctx)

	clientCtx := client.GetClientContextFromCmd(cmd)
	factory, err := tx.NewFactoryCLI(clientCtx, cmd.Flags())
	if err != nil {
		return err
	}

	// NB: while we don't need to inject the key itself (just the name),
	// we should ensure that a key with the given name exists, otherwise
	// return the error.
	// QUESTION: does `clientCtx.GetFromName()` get a default value?
	fmt.Printf("key: %s\n", clientCtx.GetFromName())
	key, err := factory.Keybase().Key(clientCtx.GetFromName())
	if err != nil {
		// TECHDEBT:
		return err
	}

	di.Provide(modules.ClientCtxInjectionToken, clientCtx, injector)
	di.Provide(modules.TxFactoryInjectionToken, factory, injector)
	di.Provide(modules.PocketNetworkClientToken, cosmos.NewLocalCosmosPocketClient(), injector)
	di.Provide(modules.KeyNameInjectionToken, key.Name, injector)
	di.Provide(modules.ServicerToken, NewServicerModule(), injector)

	srvcr := di.HydrateMain(modules.ServicerToken, injector)
	if err := srvcr.Start(); err != nil {
		return err
	}

	return nil
}
