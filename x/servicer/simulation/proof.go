package simulation

import (
	"math/rand"

	"poktroll/x/servicer/keeper"
	"poktroll/x/servicer/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgProof(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgProof{
			Servicer: simAccount.Address.String(),
		}

		// TODO: Handling the Proof simulation

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "Proof simulation not implemented"), nil, nil
	}
}
