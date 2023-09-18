package cli

import (
	"strconv"

	"poktroll/x/session/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdGetSession() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-session",
		Short: "Query get-session",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			appAddress := args[0]
			serviceId := args[1]
			blockHeight, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetSessionRequest{
				BlockHeight: blockHeight,
				AppAddress:  appAddress,
				ServiceId:   serviceId,
			}

			res, err := queryClient.GetSession(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
