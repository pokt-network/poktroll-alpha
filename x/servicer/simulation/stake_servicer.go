package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"poktroll/x/servicer/keeper"
	"poktroll/x/servicer/types"
)

func SimulateMsgStakeServicer(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgStakeServicer{
			Address: simAccount.Address.String(),
		}

		// TODO: Handling the StakeServicer simulation

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "StakeServicer simulation not implemented"), nil, nil
	}
}
