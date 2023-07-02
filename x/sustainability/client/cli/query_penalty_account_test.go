package cli_test

import (
	"fmt"
	"testing"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"google.golang.org/grpc/status"

	"github.com/cascadiafoundation/github.com/cascadiafoundation/cascadia/testutil/network"
	"github.com/cascadiafoundation/github.com/cascadiafoundation/cascadia/testutil/nullify"

	"github.com/cascadiafoundation/cascadia/x/sustainability/client/cli"
	"github.com/cascadiafoundation/cascadia/x/sustainability/types"
)

func networkWithPenaltyAccountObjects(t *testing.T) (*network.Network, types.PenaltyAccount) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	penaltyAccount := &types.PenaltyAccount{}
	nullify.Fill(&penaltyAccount)
	state.PenaltyAccount = penaltyAccount
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), *state.PenaltyAccount
}

func TestShowPenaltyAccount(t *testing.T) {
	net, obj := networkWithPenaltyAccountObjects(t)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	for _, tc := range []struct {
		desc string
		args []string
		err  error
		obj  types.PenaltyAccount
	}{
		{
			desc: "get",
			args: common,
			obj:  obj,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			var args []string
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowPenaltyAccount(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryGetPenaltyAccountResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.PenaltyAccount)
				require.Equal(t,
					nullify.Fill(&tc.obj),
					nullify.Fill(&resp.PenaltyAccount),
				)
			}
		})
	}
}