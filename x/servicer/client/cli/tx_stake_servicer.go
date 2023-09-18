package cli

import (
	"encoding/json"
	"os"
	"strconv"

	"poktroll/x/servicer/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdStakeServicer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stake-servicer",
		Short: "Broadcast message stake-servicer",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			pathToAppConfig := args[0]
			contents, err := os.ReadFile(pathToAppConfig)
			if err != nil {
				return err
			}

			var stakeMsg types.MsgStakeServicer
			if err := json.Unmarshal(contents, &stakeMsg); err != nil {
				return err
			}
			stakeMsg.Address = clientCtx.GetFromAddress().String()
			if err := stakeMsg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &stakeMsg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
