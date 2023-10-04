package cli

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

func CmdGetPortalAllowlist() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-portal-allowlist [portal-address]",
		Short: "Query get-portal-allowlist",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqPortalAddress := args[0]

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetPortalAllowlistRequest{
				PortalAddress: reqPortalAddress,
			}

			res, err := queryClient.GetPortalAllowlist(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdGetDelegatedPortals() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-delegated-portals [app-address]",
		Short: "Query get-delegated-portals",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqAppAddress := args[0]
			addr, err := sdk.AccAddressFromBech32(reqAppAddress)
			if err != nil {
				return fmt.Errorf("invalid app address [%s]: %w", reqAppAddress, err)
			}

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetDelegatedPortalsRequest{
				AppAddress: addr.String(),
			}

			res, err := queryClient.GetDelegatedPortals(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
