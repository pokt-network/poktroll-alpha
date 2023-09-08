package runtime

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"poktroll/runtime/configs"
	"poktroll/runtime/di"
)

type Manager struct {
	config       *configs.Config
	genesisState json.RawMessage
}

func NewManager(config *configs.Config, genesis json.RawMessage) *Manager {
	mgr := new(Manager)
	mgr.config = config
	mgr.genesisState = genesis

	return mgr
}

func NewManagerFromFiles(configPath, genesisPath string) *Manager {
	cfg, genesisState, err := parseFiles(configPath, genesisPath)
	if err != nil {
		log.Fatalf("[ERROR] Failed to initialize runtime builder: %v", err)
	}
	return NewManager(cfg, genesisState)
}

func NewManagerFromReaders(configReader, genesisReader io.Reader) *Manager {
	cfg, genesisState, err := parseFromReaders(configReader, genesisReader)
	if err != nil {
		log.Fatalf("[ERROR] Failed to initialize runtime builder: %v", err)
	}
	return NewManager(cfg, genesisState)
}

func (m *Manager) Start() error                        { return nil }
func (m *Manager) CascadeStart() error                 { return nil }
func (m *Manager) Hydrate(_ *di.Injector, _ *[]string) {}

func (m *Manager) GetConfig() *configs.Config {
	return m.config
}

func (m *Manager) GetGenesis() json.RawMessage {
	return m.genesisState
}

func parseFiles(configJSONPath, genesisJSONPath string) (*configs.Config, json.RawMessage, error) {
	config := configs.ParseConfig(configJSONPath)
	genesisState, err := parseGenesis(genesisJSONPath)
	if err != nil {
		return nil, nil, err
	}
	return config, genesisState, nil
}

func parseFromReaders(configReader, genesisReader io.Reader) (*configs.Config, json.RawMessage, error) {
	cfgBz, err := io.ReadAll(configReader)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading config: %w", err)
	}
	cfg := new(configs.Config)
	if err := json.Unmarshal(cfgBz, cfg); err != nil {
		return nil, nil, fmt.Errorf("error unmarshaling config: %w", err)
	}
	genBz, err := io.ReadAll(genesisReader)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading genesis: %w", err)
	}
	return cfg, json.RawMessage(genBz), nil
}

// Manager option helpers

func WithClientDebugMode() func(*Manager) {
	return func(b *Manager) {
		b.config.ClientDebugMode = true
	}
}
