package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"poktroll/x/servicer/types"
)

var _ = strconv.Itoa(0)

func CmdClaims() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claims [servicer-address]",
		Short: "Query claims",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqServicerAddress := args[0]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryClaimsRequest{
				ServicerAddress: reqServicerAddress,
			}

			res, err := queryClient.Claims(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
