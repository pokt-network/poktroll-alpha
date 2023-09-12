package logger

import (
	"fmt"
	"github.com/cometbft/cometbft/libs/log"

	"poktroll/runtime/di"
)

var (
	// TECHDEBT: LoggerModule's `Logger` should be an interface.
	//_ modules.LoggerModule = &Logger{}
	CosmosLoggerToken = di.NewInjectionToken[*Logger]("logger")
)

type Logger struct {
	log.Logger
}

func NewLogger(logger log.Logger) *Logger {
	return &Logger{Logger: logger}
}

func (l *Logger) CreateLoggerForModule(moduleId string) *Logger {
	return &Logger{
		Logger: l.Logger.With("module", moduleId),
	}
}

func (l *Logger) Hydrate(_ *di.Injector, _ *[]string) {}
func (l *Logger) CascadeStart() error {
	return nil
}
func (l *Logger) Start() error {
	return nil
}
func (l *Logger) Stop() error {
	return fmt.Errorf("cannot stop cosmos logger implementation")
}
