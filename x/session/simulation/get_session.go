package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"poktroll/x/session/keeper"
	"poktroll/x/session/types"
)

func SimulateMsgGetSession(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgGetSession{
			Address: simAccount.Address.String(),
		}

		// TODO: Handling the GetSession simulation

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "GetSession simulation not implemented"), nil, nil
	}
}
