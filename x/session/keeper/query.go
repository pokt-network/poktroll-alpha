package keeper

import (
	"poktroll/x/session/types"
)

var _ types.QueryServer = Keeper{}
