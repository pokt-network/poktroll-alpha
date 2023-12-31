package cli_test

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"cosmossdk.io/math"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"poktroll/testutil/network"
	"poktroll/testutil/nullify"
	"poktroll/x/servicer/client/cli"
	"poktroll/x/servicer/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func networkWithServicersObjects(t *testing.T, n int) (*network.Network, []types.Servicers) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	accounts := simtypes.RandomAccounts(r, n)

	for i := int64(0); i < int64(n); i++ {
		account := accounts[i].Address.String()
		coin := sdk.NewCoin("stake", math.NewInt(i))
		application := types.Servicers{
			Address: account,
			Stake:   &coin,
		}
		nullify.Fill(&application)
		state.ServicersList = append(state.ServicersList, application)
	}

	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), state.ServicersList
}

func TestShowServicers(t *testing.T) {
	net, objs := networkWithServicersObjects(t, 2)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	tests := []struct {
		desc    string
		address string

		args []string
		err  error
		obj  types.Servicers
	}{
		{
			desc:    "found",
			address: objs[0].Address,

			args: common,
			obj:  objs[0],
		},
		{
			desc:    "not found",
			address: strconv.Itoa(100000),

			args: common,
			err:  status.Error(codes.NotFound, "not found"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				tc.address,
			}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowServicers(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryGetServicersResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.Servicers)
				require.Equal(t,
					nullify.Fill(&tc.obj),
					nullify.Fill(&resp.Servicers),
				)
			}
		})
	}
}

func TestListServicers(t *testing.T) {
	net, objs := networkWithServicersObjects(t, 5)

	ctx := net.Validators[0].ClientCtx
	request := func(next []byte, offset, limit uint64, total bool) []string {
		args := []string{
			fmt.Sprintf("--%s=json", tmcli.OutputFlag),
		}
		if next == nil {
			args = append(args, fmt.Sprintf("--%s=%d", flags.FlagOffset, offset))
		} else {
			args = append(args, fmt.Sprintf("--%s=%s", flags.FlagPageKey, next))
		}
		args = append(args, fmt.Sprintf("--%s=%d", flags.FlagLimit, limit))
		if total {
			args = append(args, fmt.Sprintf("--%s", flags.FlagCountTotal))
		}
		return args
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(objs); i += step {
			args := request(nil, uint64(i), uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListServicers(), args)
			require.NoError(t, err)
			var resp types.QueryAllServicersResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.Servicers), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.Servicers),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListServicers(), args)
			require.NoError(t, err)
			var resp types.QueryAllServicersResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.Servicers), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.Servicers),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListServicers(), args)
		require.NoError(t, err)
		var resp types.QueryAllServicersResponse
		require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NoError(t, err)
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(objs),
			nullify.Fill(resp.Servicers),
		)
	})
}
