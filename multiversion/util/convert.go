package util

import (
	"github.com/oomph-ac/mv/multiversion/latest"
	"github.com/oomph-ac/mv/multiversion/mappings"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// LatestAirRID is the runtime ID of the air block in the latest version of the game.
var LatestAirRID, _ = latest.StateToRuntimeID("minecraft:air", nil)

// DowngradeItem downgrades the input item stack to a legacy item stack. It returns a boolean indicating if the item was
// downgraded successfully.
func DowngradeItem(input protocol.ItemStack, mappings mappings.MVMapping) protocol.ItemStack {
	name, _ := latest.ItemRuntimeIDToName(input.NetworkID)
	networkID, _ := mappings.ItemIDByName(name)
	input.ItemType.NetworkID = networkID
	if input.BlockRuntimeID > 0 {
		input.BlockRuntimeID = int32(DowngradeBlockRuntimeID(uint32(input.BlockRuntimeID), mappings))
	}
	return input
}

// UpgradeItem upgrades the input item stack to a latest item stack. It returns a boolean indicating if the item was
// upgraded successfully.
func UpgradeItem(input protocol.ItemStack, mappings mappings.MVMapping) protocol.ItemStack {
	if input.ItemType.NetworkID == 0 {
		return protocol.ItemStack{}
	}
	name, _ := mappings.ItemNameByID(input.ItemType.NetworkID)
	networkID, _ := latest.ItemNameToRuntimeID(name)
	input.ItemType.NetworkID = networkID
	if input.BlockRuntimeID > 0 {
		input.BlockRuntimeID = int32(UpgradeBlockRuntimeID(uint32(input.BlockRuntimeID), mappings))
	}
	return input
}

// DowngradeBlockRuntimeID downgrades a latest block runtime ID to a legacy block runtime ID.
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
