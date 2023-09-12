package servicer

import (
	"poktroll/logger"
	"poktroll/modules"
	"poktroll/runtime/di"
)

type servicer struct {
	pocketNetworkClient modules.PocketNetworkClient
	relayer             modules.RelayerModule
	sessionManager      modules.SessionManager
	miner               modules.MinerModule
	logger              logger.CosmosLogger
}

func NewServicerModule() modules.ServicerModule {
	return &servicer{}
}

func (r *servicer) Hydrate(injector *di.Injector, path *[]string) {
	r.pocketNetworkClient = di.Hydrate(modules.PocketNetworkClientToken, injector, path)
	r.relayer = di.Hydrate(modules.RelayerToken, injector, path)
	r.miner = di.Hydrate(modules.MinerModuleToken, injector, path)
	r.sessionManager = di.Hydrate(modules.SessionManagerToken, injector, path)
	globalLogger := di.Hydrate(logger.CosmosLoggerToken, injector, path)
	r.logger = *globalLogger.CreateLoggerForModule(modules.ServicerToken.Id())
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
