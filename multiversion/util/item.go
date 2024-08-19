package util

import (
	"github.com/oomph-ac/mv/multiversion/latest"
	"github.com/oomph-ac/mv/multiversion/mappings"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// DowngradeItem downgrades the input item stack to a legacy item stack. It returns a boolean indicating if the item was
// downgraded successfully.
func DowngradeItem(input protocol.ItemStack, mappings mappings.MVMapping) protocol.ItemStack {
	name, _ := latest.ItemRuntimeIDToName(input.NetworkID)
	networkID, ok := mappings.ItemIDByName(name)
	if !ok {
		return input
	}

	input.ItemType.NetworkID = networkID
	if input.BlockRuntimeID > 0 {
		input.BlockRuntimeID = int32(DowngradeBlockRuntimeID(uint32(input.BlockRuntimeID), mappings))
	}
	return input
}

// UpgradeItem upgrades the input item stack to the latest item stack. It returns a boolean indicating if the item was
// upgraded successfully.
func UpgradeItem(input protocol.ItemStack, mappings mappings.MVMapping) protocol.ItemStack {
	if input.ItemType.NetworkID == 0 {
		return protocol.ItemStack{}
	}

	name, _ := mappings.ItemNameByID(input.ItemType.NetworkID)
	networkID, ok := latest.ItemNameToRuntimeID(name)
	if !ok {
		return input
	}

	input.ItemType.NetworkID = networkID
	if input.BlockRuntimeID > 0 {
		input.BlockRuntimeID = int32(UpgradeBlockRuntimeID(uint32(input.BlockRuntimeID), mappings))
	}
	return input
}
