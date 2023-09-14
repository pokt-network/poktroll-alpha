package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
)

type SessionI interface {
	proto.Message

	GetAllApplication(ctx sdk.Context) (list []ApplicationI)
	GetServicers(ctx sdk.Context) (list []ServicerI)
}
