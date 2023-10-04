package types

import (
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	GetPubKey(ctx sdk.Context, addr sdk.AccAddress) (pk cryptotypes.PubKey, err error)
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.

type BankKeeper interface {
	// SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	// SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

type PortalKeeper interface {
	SetDelegator(ctx sdk.Context, appAddress string, delegatedPortals Delegatees)
	GetAllowlist(ctx sdk.Context, portalAddress string) (val []string, found bool)
}
