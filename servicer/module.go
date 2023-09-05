package servicer

import (
	"context"

	"github.com/pokt-network/smt"

	"poktroll/modules"
	"poktroll/runtime/di"
	"poktroll/shared/crypto"
)

type servicer struct {
	pocketNetworkClient modules.PocketNetworkClient
	relayer             modules.RelayerModule
	sessionManager      modules.SessionManager
	proofManager        modules.ProofManager
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
	r.proofManager = di.Resolve(modules.ProofManagerToken, injector, path)
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
	for range ch {
		claim := r.proofManager.Root()
		if claim == nil {
			continue
		}
		<-r.pocketNetworkClient.SubmitClaim(context.TODO(), claim)

		// Wait for some time
		//key := session.BlockHash
		key := r.tempHash
		value, _, err := r.proofManager.Get([]byte(key))
		if err != nil {
			r.logger.Error().AnErr("key", err).Msg("getting key")
		}

		proof, err := r.proofManager.Prove(r.tempHash)
		if err != nil {
			r.logger.Error().AnErr("proof", err).Msg("getting proof")
		}

		result := smt.VerifySumProof(proof, claim, key, value, r.proofManager.Sum(), r.proofManager.Spec())
		r.logger.Info().Bool("result", result).Msg("Sumproof")

		<-r.pocketNetworkClient.SubmitProof(context.TODO(), proof)
	}
}

func (r *servicer) handleRelays() {
	ch := r.relayer.Relays()
	for relay := range ch {
		serializedRelay := relay.Serialize()
		hash := crypto.SHA3Hash([]byte(serializedRelay))
		r.tempHash = hash
		r.proofManager.Update(hash, hash, 1)
	}
}
