package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"poktroll/x/servicer/types"
)

var _ = strconv.Itoa(0)

func CmdProofs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proofs [servicer-address]",
		Short: "Query proofs",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqServicerAddress := args[0]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryProofsRequest{

				ServicerAddress: reqServicerAddress,
			}

			res, err := queryClient.Proofs(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
