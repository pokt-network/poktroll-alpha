package types

import (
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
)

type ServicerI interface {
	proto.Message

	GetAddress() string
	GetStake() *types.Coin
}
	