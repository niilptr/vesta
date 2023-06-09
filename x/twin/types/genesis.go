package types

import (
	"fmt"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		TwinList:      []Twin{},
		TrainingState: nil,
		Params:        DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in twin
	twinIndexMap := make(map[string]struct{})

	for _, elem := range gs.TwinList {
		index := string(TwinKey(elem.Name))
		if _, ok := twinIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for twin")
		}
		twinIndexMap[index] = struct{}{}
	}

	return gs.Params.Validate()
}
