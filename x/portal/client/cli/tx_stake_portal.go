package cli

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"
	"strings"

	"poktroll/x/portal/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

var allowlistedApps = make([]string, 0)

func CmdStakePortal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stake-portal",
		Short: "Broadcast message stake-portal",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			stakeAmountString := args[0]
			serviceIdsCommaSeparated := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			stakeAmount, err := sdk.ParseCoinNormalized(stakeAmountString)
			if err != nil {
				return err
			}
			msg := types.NewMsgStakePortal(
				clientCtx.GetFromAddress().String(),
				stakeAmount,
				strings.Split(serviceIdsCommaSeparated, ","),
				allowlistedApps,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	f := cmd.Flags()
	f.StringSliceVar(&allowlistedApps, "allowlisted-apps", []string{}, "comma separated list of allowlisted application addresses")

	return cmd
}
