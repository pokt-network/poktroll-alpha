package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"poktroll/x/application/types"
)

var _ = strconv.Itoa(0)

func CmdUndelegateFromPortal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "undelegate-from-portal [portal_pub_key]",
		Short: "Broadcast message undelegate_from_portal",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			portalPubKey := args[0]
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := types.NewMsgUndelegateFromPortal(
				clientCtx.GetFromAddress().String(),
				portalPubKey,
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
