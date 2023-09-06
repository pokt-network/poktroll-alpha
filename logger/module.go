package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"

	"poktroll/modules"
	"poktroll/runtime/configs"
	"poktroll/runtime/di"
)

type loggerModule struct {
	zerolog.Logger
	config *configs.LoggerConfig
}

var logLevelToZeroLog = map[configs.LogLevel]zerolog.Level{
	configs.LogLevel_LOG_LEVEL_UNSPECIFIED: zerolog.NoLevel,
	configs.LogLevel_LOG_LEVEL_DEBUG:       zerolog.DebugLevel,
	configs.LogLevel_LOG_LEVEL_INFO:        zerolog.InfoLevel,
	configs.LogLevel_LOG_LEVEL_WARN:        zerolog.WarnLevel,
	configs.LogLevel_LOG_LEVEL_ERROR:       zerolog.ErrorLevel,
	configs.LogLevel_LOG_LEVEL_FATAL:       zerolog.FatalLevel,
	configs.LogLevel_LOG_LEVEL_PANIC:       zerolog.PanicLevel,
}

var logFormatToEnum = map[string]configs.LogFormat{
	"json":   configs.LogFormat_LOG_FORMAT_JSON,
	"pretty": configs.LogFormat_LOG_FORMAT_PRETTY,
}

func NewGlobalLogger(cfg *configs.LoggerConfig) *loggerModule {
	return &loggerModule{
		Logger: zerolog.New(os.Stdout).With().Timestamp().Logger(),
		config: cfg,
	}
}

func (m *loggerModule) Module() modules.LoggerModule { return m }

// CreateLoggerForModule implements the respective `modules.Logger` interface member.
func (m *loggerModule) CreateLoggerForModule(moduleName string) *modules.Logger {
	logger := m.Logger.With().Str("module", moduleName).Logger()
	return &logger
}

func (m *loggerModule) Resolve(injector *di.Injector, path *[]string) {
	m.Create()
}

func (m *loggerModule) Create() {
	m.CreateLoggerForModule("global")
	if pocketLogLevel, ok := configs.LogLevel_value[`LOG_LEVEL_`+strings.ToUpper(m.config.GetLevel())]; ok {
		zerolog.SetGlobalLevel(logLevelToZeroLog[configs.LogLevel(pocketLogLevel)])
	} else {
		zerolog.SetGlobalLevel(zerolog.NoLevel)
	}

	if logFormatToEnum[m.config.GetFormat()] == configs.LogFormat_LOG_FORMAT_PRETTY {
		logStructure := zerolog.ConsoleWriter{Out: os.Stdout}
		logStructure.FormatLevel = func(i interface{}) string {
			return fmt.Sprintf("level=%s", strings.ToUpper(i.(string)))
		}

		m.Logger = m.Logger.Output(logStructure)
		m.Logger.Info().Msg("using pretty log format")
	}
}

func (m *loggerModule) CascadeStart() error {
	return m.Start()
}

func (m *loggerModule) Start() error {
	m.Logger = *m.CreateLoggerForModule("global")
	return nil
}

func (m *loggerModule) GetLogger() modules.Logger {
	return m.Logger
}

// SetFields sets the fields for the global logger
func (m *loggerModule) SetFields(fields map[string]any) {
	m.Logger = m.Logger.With().Fields(fields).Logger()
}

// UpdateFields updates the fields for the global logger
func (m *loggerModule) UpdateFields(fields map[string]any) {
	m.Logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
		for k, v := range fields {
			c = c.Interface(k, v)
		}
		return c
	})
}
