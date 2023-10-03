package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/application module sentinel errors
var (
	ErrNilStakeAmount          = sdkerrors.Register(ModuleName, 1, "Stake amount is nil")
	ErrEmptyStakeAmount        = sdkerrors.Register(ModuleName, 2, "Stake amount is empty")
	ErrStakeAmountMustBeHigher = sdkerrors.Register(ModuleName, 3, "The stake amount for a previously staked application must be explicitly higher than the prior amount")
	ErrUnstakingNonExistentApp = sdkerrors.Register(ModuleName, 4, "Could not unstake non-existent application")
	ErrNoServicesToStake       = sdkerrors.Register(ModuleName, 5, "Must stake for at least one service")
	ErrApplicationNotFound     = sdkerrors.Register(ModuleName, 6, "Application not found")
	ErrMaxDelegatedReached     = sdkerrors.Register(ModuleName, 7, "Application has reached the maximum number of delegated portals")
	ErrPortalAlreadyDelegated  = sdkerrors.Register(ModuleName, 8, "Application has already delegated to this portal")
	ErrPortalNotDelegated      = sdkerrors.Register(ModuleName, 9, "Application has not delegated to this portal")
	ErrPortalNotFound          = sdkerrors.Register(ModuleName, 10, "Portal not found")
	ErrAppNotWhitelisted       = sdkerrors.Register(ModuleName, 11, "Application is not whitelisted for this portal")
	ErrCannotUndelegateSelf    = sdkerrors.Register(ModuleName, 12, "Cannot undelegate self")
	ErrNotSelfDelegated        = sdkerrors.Register(ModuleName, 13, "Application is not self-delegated")
)
