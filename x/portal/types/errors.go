package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/portal module sentinel errors
var (
	ErrNilStakeAmount             = sdkerrors.Register(ModuleName, 1, "Stake amount is nil")
	ErrEmptyStakeAmount           = sdkerrors.Register(ModuleName, 2, "Stake amount is empty")
	ErrStakeAmountMustBeHigher    = sdkerrors.Register(ModuleName, 3, "The stake amount for a previously staked portal must be explicitly higher than the prior amount")
	ErrUnstakingNonExistentPortal = sdkerrors.Register(ModuleName, 4, "Could not unstake non-existent portal")
	ErrNoServicesToStake          = sdkerrors.Register(ModuleName, 5, "Must stake for at least one service")
	ErrAppAlreadyWhitelisted      = sdkerrors.Register(ModuleName, 6, "Application already whitelisted by the portal")
	ErrPortalNotFound             = sdkerrors.Register(ModuleName, 7, "Portal not found in state")
	ErrAppNotWhitelisted          = sdkerrors.Register(ModuleName, 8, "Application not whitelisted by the portal")
)
