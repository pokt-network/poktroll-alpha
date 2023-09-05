package runtime

import (
	"encoding/json"
	"fmt"
	"os"

	nm "github.com/cometbft/cometbft/node"
	"github.com/cometbft/cometbft/types"
)

// parseGenesis parses the genesis file in JSON format and returns it as a json.RawMessage
func parseGenesis(genesisJSONPath string) (json.RawMessage, error) {
	data, err := os.ReadFile(genesisJSONPath)
	if err != nil {
		return nil, fmt.Errorf("error reading genesis file %s: %w", genesisJSONPath, err)
	}
	return json.RawMessage(data), nil
}

func GenesisDocProvider(cmtGenesisJSONPath, genesisJSONPath string) nm.GenesisDocProvider {
	return func() (*types.GenesisDoc, error) {
		config, err := types.GenesisDocFromFile(cmtGenesisJSONPath)
		if err != nil {
			return nil, err
		}
		genesisState, err := parseGenesis(genesisJSONPath)
		if err != nil {
			return nil, err
		}
		config.AppState = genesisState
		return config, nil
	}
}
