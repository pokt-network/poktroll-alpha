package keeper

import (
	"poktroll/x/application/types"
)

var _ types.QueryServer = Keeper{}
