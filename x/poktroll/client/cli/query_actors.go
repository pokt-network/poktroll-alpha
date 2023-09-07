package cli

import (
	"strconv"

	"poktroll/x/poktroll/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdActors() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "actors",
		Short: "Query actors",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryActorsRequest{}

			res, err := queryClient.Actors(cmd.Context(), params)
			if err != nil {
				clientCtx.PrintString("could not get actors")
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
