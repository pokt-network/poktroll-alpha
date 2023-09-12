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
)
