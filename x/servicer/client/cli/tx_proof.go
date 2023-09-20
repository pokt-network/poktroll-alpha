package cli

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"poktroll/x/servicer/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdProof() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proof [root hex] [path hex] [value-hash hex] [sum] [proof-bz hex]",
		Short: "Broadcast message proof",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argRoot, err := hex.DecodeString(args[0])
			if err != nil {
				return fmt.Errorf("unable to hex decode root hash argument: %w", err)
			}

			argPath, err := hex.DecodeString(args[1])
			if err != nil {
				return fmt.Errorf("unable to hex decode path argument: %w", err)
			}

			argValueHash, err := hex.DecodeString(args[2])
			if err != nil {
				return fmt.Errorf("unable to hex decode value hash argument: %w", err)
			}

			argSum, err := cast.ToUint64E(args[3])
			if err != nil {
				return err
			}

			argProofBz, err := hex.DecodeString(args[4])
			if err != nil {
				return fmt.Errorf("unable to hex decode proof argument: %w", err)
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg, err := types.NewMsgProof(
				clientCtx.GetFromAddress().String(),
				argRoot,
				argPath,
				argValueHash,
				argSum,
				argProofBz,
			)
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
