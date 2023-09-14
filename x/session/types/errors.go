package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/session module sentinel errors
var (
	ErrRetrieveApps     = sdkerrors.Register(ModuleName, 1, "Could not retrieve apps")
	ErrNoServicersFound = sdkerrors.Register(ModuleName, 2, "No servicers found")
	ErrFindApp          = sdkerrors.Register(ModuleName, 3, "Could not find app")
)
