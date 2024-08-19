package util

import (
	"github.com/oomph-ac/mv/multiversion/latest"
	"github.com/oomph-ac/mv/multiversion/mappings"
)

// DowngradeBlockRuntimeID downgrades the latest block runtime ID to a legacy block runtime ID.
func DowngradeBlockRuntimeID(input uint32, mappings mappings.MVMapping) uint32 {
	name, properties, ok := latest.RuntimeIDToState(input)
	if !ok {
		return mappings.LegacyAirRID
	}

	return mappings.StateToRuntimeID(name, properties)
}

// UpgradeBlockRuntimeID upgrades a legacy block runtime ID to a latest block runtime ID.
func UpgradeBlockRuntimeID(input uint32, mappings mappings.MVMapping) uint32 {
	name, properties, ok := mappings.RuntimeIDToState(input)
	if !ok {
		return LatestAirRID
	}

	runtimeID, ok := latest.StateToRuntimeID(name, properties)
	if !ok {
		return LatestAirRID
	}
	return runtimeID
}
