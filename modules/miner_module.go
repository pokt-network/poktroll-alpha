package modules

import (
	"poktroll/runtime/di"
)

var MinerModuleToken = di.NewInjectionToken[MinerModule]("miner")

type MinerModule interface {
	di.Module
	Update(key []byte, value []byte, weight uint64) error
	SubmitClaim() error
	SubmitProof(key []byte) error
}
