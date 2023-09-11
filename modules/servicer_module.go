package modules

import (
	cosmosClient "github.com/cosmos/cosmos-sdk/client"
	txClient "github.com/cosmos/cosmos-sdk/client/tx"
	"poktroll/runtime/di"
)

var (
	ServicerToken = di.NewInjectionToken[ServicerModule]("servicer")
	// CONSIDERATION: use "pocket" key interface types instead of cosmos-sdk's.
	//PrivateKeyInjectionToken = di.NewInjectionToken[*secp256k1.PrivKey]("servicerPrivateKey")
	KeyNameInjectionToken   = di.NewInjectionToken[string]("servicerKeyName")
	TxFactoryInjectionToken = di.NewInjectionToken[txClient.Factory]("servicerTxFactory")
	ClientCtxInjectionToken = di.NewInjectionToken[cosmosClient.Context]("servicerClientCtx")
)

type ServicerModule interface {
	di.Module
}
