package types

import (
	"fmt"

	sharedtypes "poktroll/x/shared/types"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		ServicersList: []sharedtypes.Servicers{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in servicers
	servicersIndexMap := make(map[string]struct{})

	for _, elem := range gs.ServicersList {
		index := string(ServicersKey(elem.Address))
		if _, ok := servicersIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for servicers")
		}
		servicersIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
