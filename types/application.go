package types

import (
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
)

type ApplicationI interface {
	proto.Message

	GetAddress() string
	GetStake() *types.Coin
}
