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

func CmdDelegateToPortal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegate-to-portal [portal-pub-key]",
		Short: "Broadcast message DelegateToPortal",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argPortalPubKey := args[0]
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := types.NewMsgDelegateToPortal(
				clientCtx.GetFromAddress().String(),
				argPortalPubKey,
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
