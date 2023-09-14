package cli

import (
	"strconv"

	"poktroll/x/servicer/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdStakeServicer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stake-servicer",
		Short: "Broadcast message stake-servicer",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			stakeAmountString := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			stakeAmount, err := sdk.ParseCoinNormalized(stakeAmountString)
			if err != nil {
				return err
			}

			msg := types.NewMsgStakeServicer(
				clientCtx.GetFromAddress().String(),
				stakeAmount,
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
