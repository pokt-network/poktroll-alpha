package servicer

import (
	"poktroll/modules"
	"poktroll/runtime/di"
	"poktroll/shared/crypto"
)

type servicer struct {
	pocketNetworkClient modules.PocketNetworkClient
	relayer             modules.RelayerModule
	sessionManager      modules.SessionManager
	miner               modules.MinerModule
	logger              *modules.Logger
	PrivateKey          crypto.PrivateKey
}

func NewServicerModule() modules.ServicerModule {
	return &servicer{}
}

func (r *servicer) Hydrate(injector *di.Injector, path *[]string) {
	r.pocketNetworkClient = di.Hydrate(modules.PocketNetworkClientToken, injector, path)
	r.relayer = di.Hydrate(modules.RelayerToken, injector, path)
	r.miner = di.Hydrate(modules.MinerModuleToken, injector, path)
	r.sessionManager = di.Hydrate(modules.SessionManagerToken, injector, path)
	r.PrivateKey = di.Hydrate(modules.PrivateKeyInjectionToken, injector, path)
	globalLogger := di.Hydrate(modules.LoggerModuleToken, injector, path)
	r.logger = globalLogger.CreateLoggerForModule(modules.ServicerToken.Id())
}

func (r *servicer) CascadeStart() error {
	r.pocketNetworkClient.CascadeStart()
	r.relayer.CascadeStart()
	r.miner.CascadeStart()
	r.sessionManager.CascadeStart()
	return r.Start()
}

func (r *servicer) Start() error {
	r.miner.MineRelays(r.relayer.Relays(), r.sessionManager.ClosedSessions())
	return nil
}
