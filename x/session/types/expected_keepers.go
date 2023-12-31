package types

import (
	// poktroll "poktroll/types"

	apptypes "poktroll/x/application/types"
	svctypes "poktroll/x/servicer/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}

type ApplicationKeeper interface {
	// GetAllApplication(ctx sdk.Context) ([]apptypes.Application, error)
	GetApplication(ctx sdk.Context, address string) (app apptypes.Application, found bool)
}

type ServicerKeeper interface {
	GetAllServicers(ctx sdk.Context) []svctypes.Servicers
}
