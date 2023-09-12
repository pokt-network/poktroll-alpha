package logger

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"

	"poktroll/runtime/di"
)

var (
	// TECHDEBT: LoggerModule's `Logger` should be an interface.
	//_ modules.LoggerModule = &Logger{}
	CosmosLoggerToken = di.NewInjectionToken[*CosmosLogger]("logger")
)

type CosmosLogger struct {
	log.Logger
}

func NewLogger(logger log.Logger) *CosmosLogger {
	return &CosmosLogger{Logger: logger}
}

func (l *CosmosLogger) CreateLoggerForModule(moduleId string) *CosmosLogger {
	return &CosmosLogger{
		Logger: l.Logger.With("module", moduleId),
	}
}

func (l *CosmosLogger) Hydrate(_ *di.Injector, _ *[]string) {}
func (l *CosmosLogger) CascadeStart() error {
	return nil
}
func (l *CosmosLogger) Start() error {
	return nil
}
func (l *CosmosLogger) Stop() error {
	return fmt.Errorf("cannot stop cosmos logger implementation")
}
