package configs

import (
	"encoding/json"
	"log"
	"os"
	"poktroll/runtime/defaults"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// IMPROVE: add a SaveConfig() function to save the config to a file and
// generate a default config file for the user. Add it as a new command
// into the CLI.

type Config struct {
	RootDirectory   string `json:"root_directory"`
	ClientDebugMode bool   `json:"client_debug_mode"`

	Persistence      *PersistenceConfig `json:"persistence"`
	Logger           *LoggerConfig      `json:"logger"`
	BlocksPerSession int64              `json:"blocks_per_session"`
}

// ParseConfig parses the config file and returns a Config struct
func ParseConfig(cfgFile string) *Config {
	config := NewDefaultConfig()

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("/etc/pocket/")  // path to look for the config file in
		viper.AddConfigPath("$HOME/.pocket") // call multiple times to add many search paths
		viper.AddConfigPath(".")             // optionally look for config in the working directory
		viper.SetConfigName("config")        // name of config file (without extension)
		viper.SetConfigType("json")          // REQUIRED if the config file does not have the extension in the name
	}

	// The lines below allow for environment variables configuration (12 factor app)
	// Eg: POCKET_CONSENSUS_PRIVATE_KEY=somekey would override `consensus.private_key` in config.json
	// If the key is not set in the config, the env var will not be used.
	viper.SetEnvPrefix("POCKET")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	verbose := viper.GetBool("verbose")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok && cfgFile == "" {
			if verbose {
				log.Default().Printf("No config provided, using defaults")
			}
		} else {
			// TODO: This is a log call to avoid import cycles. Refactor logger_config.proto to avoid this.
			log.Fatalf("[ERROR] fatal error reading config file %s", err.Error())
		}
	} else {
		// TODO: This is a log call to avoid import cycles. Refactor logger_config.proto to avoid this.
		if verbose {
			log.Default().Printf("Using config file: %s", viper.ConfigFileUsed())
		}
	}

	decoderConfig := func(dc *mapstructure.DecoderConfig) {
		// This is to leverage the `json` struct tags without having to add `mapstructure` ones.
		// Until we have complex use cases, this should work just fine.
		dc.TagName = "json"
	}
	// Detect if we need to use json.Unmarshal instead of viper.Unmarshal
	if err := viper.Unmarshal(&config, decoderConfig); err != nil {
		cfgData := viper.AllSettings()
		cfgJSON, _ := json.Marshal(cfgData)

		// last ditch effort to unmarshal the config
		if err := json.Unmarshal(cfgJSON, &config); err != nil {
			log.Fatalf("[ERROR] failed to unmarshal config %s", err.Error())
		}
	}

	return config
}

// setViperDefaults this is a hacky way to set the default values for Viper so env var overrides work.
// DISCUSS: is there a better way to do this?
func setViperDefaults(cfg *Config) {
	// convert the config struct to a map with the json tags as keys
	cfgData, err := json.Marshal(cfg)
	if err != nil {
		log.Fatalf("[ERROR] failed to marshal config %s", err.Error())
	}
	var cfgMap map[string]any
	if err := json.Unmarshal(cfgData, &cfgMap); err != nil {
		log.Fatalf("[ERROR] failed to unmarshal config %s", err.Error())
	}

	for k, v := range cfgMap {
		viper.SetDefault(k, v)
	}
}

func NewDefaultConfig(options ...func(*Config)) *Config {
	cfg := &Config{
		RootDirectory: defaults.DefaultRootDirectory,
		Persistence:   &PersistenceConfig{
			// PostgresUrl:    defaults.DefaultPersistencePostgresURL,
			// BlockStorePath: defaults.DefaultPersistenceBlockStorePath,
		},
		Logger: &LoggerConfig{
			Level:  defaults.DefaultLoggerLevel,
			Format: defaults.DefaultLoggerFormat,
		},
		BlocksPerSession: defaults.DefaultBlocksPerSession,
	}

	for _, option := range options {
		option(cfg)
	}

	// set Viper defaults so POCKET_ env vars work without having to set in config file
	setViperDefaults(cfg)

	return cfg
}

// WithNodeSchema is an option to configure the schema for the node's database.
func WithNodeSchema(schema string) func(*Config) {
	return func(cfg *Config) {
		cfg.Persistence.NodeSchema = schema
	}
}

// CreateTempConfig creates a temporary config for testing purposes only
func CreateTempConfig(cfg *Config) (*Config, error) {
	tmpfile, err := os.CreateTemp("", "test_config_*.json")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpfile.Name())

	content, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}

	if _, err := tmpfile.Write(content); err != nil {
		return nil, err
	}

	if err := tmpfile.Close(); err != nil {
		return nil, err
	}

	return ParseConfig(tmpfile.Name()), nil
}