package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"poktroll/x/portal/types"
	"strconv"
)

var _ = strconv.Itoa(0)

func CmdWhitelistApplication() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whitelist-application",
		Short: "Broadcast message whitelist_application",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			appAddress := args[0]
			appAddr, err := sdk.AccAddressFromBech32(appAddress)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgWhitelistApplication(
				clientCtx.GetFromAddress().String(),
				appAddr.String(),
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
