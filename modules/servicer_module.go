package modules

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"

	"poktroll/runtime/di"
)

var (
	ServicerToken = di.NewInjectionToken[ServicerModule]("servicer")
	// CONSIDERATION: use "pocket" key interface types instead of cosmos-sdk's.
	PrivateKeyInjectionToken = di.NewInjectionToken[*secp256k1.PrivKey]("servicerPrivateKey")
	//TxFactoryInjectionToken = di.NewInjectionToken[txClient.Factory]("servicerTxFactory")
	//ClientCtxInjectionToken = di.NewInjectionToken[cosmosClient.Context]("servicerClientCtx")
)

type ServicerModule interface {
	di.Module
}
