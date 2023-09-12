package keeper

import (
	"poktroll/x/servicer/types"
)

var _ types.QueryServer = Keeper{}
