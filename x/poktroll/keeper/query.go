package keeper

import (
	"poktroll/x/poktroll/types"
)

var _ types.QueryServer = Keeper{}
