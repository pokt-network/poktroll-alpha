package types

import (
	"fmt"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		ApplicationList: []Application{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in application
	applicationIndexMap := make(map[string]struct{})

	for _, app := range gs.ApplicationList {
		addr := string(ApplicationKey(app.Address))
		if _, ok := applicationIndexMap[addr]; ok {
			return fmt.Errorf("duplicated index for application")
		}
		applicationIndexMap[addr] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
