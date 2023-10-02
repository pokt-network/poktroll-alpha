package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"poktroll/x/portal/types"
)

var _ = strconv.Itoa(0)

func CmdListPortals() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-portals",
		Short: "Query list-portals",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			params := &types.QueryAllPortalsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.Portals(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowPortal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-portal [index]",
		Short: "shows a portal",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			argAddress := args[0]
			params := &types.QueryGetPortalRequest{
				Address: argAddress,
			}

			res, err := queryClient.Portal(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdGetPortalWhitelist() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-portal-whitelist [portal-address]",
		Short: "Query get_portal_whitelist",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqPortalAddress := args[0]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetPortalWhitelistRequest{
				PortalAddress: reqPortalAddress,
			}

			res, err := queryClient.GetPortalWhitelist(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
