package servicer

import (
	"fmt"

	"github.com/spf13/cobra"
)

// define dependencies' injector tokens
// - rpc port

func GetServicerCmd() *cobra.Command {
	// inject servicer module dependencies (or refactor somewhere) (e.g. configs)

	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Servicer commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       servicerCmd,
	}

	//cmd.AddCommand(...)

	return cmd
}

func servicerCmd(cmd *cobra.Command, args []string) error {
	// resolve servicer module
	// start servicer
	return nil
}
