package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"poktroll/x/servicer/types"
)

var _ = strconv.Itoa(0)

func CmdProof() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proof [root] [path] [value-hash] [sum] [proof-bz]",
		Short: "Broadcast message proof",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argRoot := args[0]
			argPath := args[1]
			argValueHash := args[2]
			argSum, err := cast.ToInt32E(args[3])
			if err != nil {
				return err
			}
			argProofBz := args[4]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgProof(
				clientCtx.GetFromAddress().String(),
				argRoot,
				argPath,
				argValueHash,
				argSum,
				argProofBz,
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
