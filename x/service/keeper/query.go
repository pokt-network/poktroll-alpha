package keeper

import (
	"poktroll/x/service/types"
)

var _ types.QueryServer = Keeper{}
