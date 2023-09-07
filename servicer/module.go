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

	// Remove this when proof manager is updated
	tempHash []byte
}

func NewServicerModule() modules.ServicerModule {
	return &servicer{}
}

func (r *servicer) Module() modules.ServicerModule { return r }

func (r *servicer) Resolve(injector *di.Injector, path *[]string) {
	r.pocketNetworkClient = di.Resolve(modules.PocketNetworkClientToken, injector, path)
	r.relayer = di.Resolve(modules.RelayerToken, injector, path)
	r.sessionManager = di.Resolve(modules.SessionManagerToken, injector, path)
	r.miner = di.Resolve(modules.MinerModuleToken, injector, path)
	r.PrivateKey = di.Resolve(modules.PrivateKeyInjectionToken, injector, path)

	globalLogger := di.Resolve(modules.LoggerModuleToken, injector, path)
	r.logger = globalLogger.CreateLoggerForModule(modules.ServicerToken.Id())
}

func (r *servicer) CascadeStart() error {
	r.pocketNetworkClient.CascadeStart()
	r.relayer.CascadeStart()
	r.sessionManager.CascadeStart()
	return r.Start()
}

func (r *servicer) Start() error {
	go r.handleSessionEnd()
	go r.handleRelays()
	return nil
}

func (r *servicer) handleSessionEnd() {
	ch := r.sessionManager.OnSessionEnd()
	for session := range ch {
		if err := r.miner.SubmitClaim(); err != nil {
			continue
		}

		// Wait for some time
		r.miner.SubmitProof([]byte(session.BlockHash))
	}
}

func (r *servicer) handleRelays() {
	ch := r.relayer.Relays()
	for relay := range ch {
		serializedRelay := relay.Serialize()
		hash := crypto.SHA3Hash([]byte(serializedRelay))
		r.tempHash = hash
		r.miner.Update(hash, hash, 1)
	}
}
