package mappings

import (
	_ "embed"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// MVItemMapping holds all data items related.
type MVItemMapping struct {
	// items holds a list of all existing items in the game.
	items []protocol.ItemEntry
	// itemRuntimeIDsToNames holds a map to translate item runtime IDs to string IDs.
	itemRuntimeIDsToNames map[int32]string
	// itemNamesToRuntimeIDs holds a map to translate item string IDs to runtime IDs.
	itemNamesToRuntimeIDs map[string]int32

	recipes []protocol.Recipe
}

// ItemMapping returns MVItemMapping instance of all item entries and runtime ID maps from the resource JSON.
func itemMapping(itemRuntimeIDData []byte) MVItemMapping {
	var m map[string]int32
	err := nbt.Unmarshal(itemRuntimeIDData, &m)
	if err != nil {
		panic(err)
	}

	var items []protocol.ItemEntry
	var itemRuntimeIDsToNames = make(map[int32]string)
	var itemNamesToRuntimeIDs = make(map[string]int32)

	for name, rid := range m {
		items = append(items, protocol.ItemEntry{
			Name:      name,
			RuntimeID: int16(rid),
		})
		itemNamesToRuntimeIDs[name] = rid
		itemRuntimeIDsToNames[rid] = name
	}

	return MVItemMapping{
		items:                 items,
		itemRuntimeIDsToNames: itemRuntimeIDsToNames,
		itemNamesToRuntimeIDs: itemNamesToRuntimeIDs,
	}
}

// ItemNameByID returns an item's name by its legacy ID.
func (m MVItemMapping) ItemNameByID(id int32) (string, bool) {
	// TODO: Properly handle item aliases.
	name, ok := m.itemRuntimeIDsToNames[id]
	return name, ok
}

// ItemIDByName returns an item's ID by its name.
func (m MVItemMapping) ItemIDByName(name string) (int32, bool) {
	// TODO: Properly handle item aliases.
	id, ok := m.itemNamesToRuntimeIDs[name]
	if !ok {
		id = m.itemNamesToRuntimeIDs["minecraft:name_tag"]
	}
	return id, ok
}

// Items returns a slice of all item entries.
func (m MVItemMapping) Items() []protocol.ItemEntry {
	return m.items
}

// Recipes returns a slice of all recipes.
func (m MVItemMapping) Recipes() []protocol.Recipe {
	return m.recipes
}
