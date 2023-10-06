package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"poktroll/x/portal/keeper"
	"poktroll/x/portal/types"
)

func SimulateMsgUnstakePortal(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgUnstakePortal{
			Address: simAccount.Address.String(),
		}

		// TODO: Handling the UnstakePortal simulation

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "UnstakePortal simulation not implemented"), nil, nil
	}
}
