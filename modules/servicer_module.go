package modules

import (
	"poktroll/runtime/di"

	cosmosClient "github.com/cosmos/cosmos-sdk/client"
	txClient "github.com/cosmos/cosmos-sdk/client/tx"
)

var (
	ServicerToken           = di.NewInjectionToken[ServicerModule]("servicer")
	KeyNameInjectionToken   = di.NewInjectionToken[string]("servicerKeyName")
	TxFactoryInjectionToken = di.NewInjectionToken[txClient.Factory]("servicerTxFactory")
	ClientCtxInjectionToken = di.NewInjectionToken[cosmosClient.Context]("servicerClientCtx")
)

type ServicerModule interface {
	di.Module
}
