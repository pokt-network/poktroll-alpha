package cli

import (
	"strconv"

	"poktroll/x/poktroll/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdServicers() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "servicers",
		Short: "Query Servicers",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			request := &types.QueryServicersRequest{}

			res, err := queryClient.Servicers(cmd.Context(), request)
			if err != nil {
				clientCtx.PrintString("could not get servicers")
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
