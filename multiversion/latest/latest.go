package latest

import (
	"bytes"
	_ "embed"

	"github.com/df-mc/worldupgrader/blockupgrader"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
)

var (
	//go:embed block_states.nbt
	BlockStateData []byte
	//go:embed item_runtime_ids.nbt
	ItemRuntimeIDData []byte

	// stateToRuntimeID maps a block state hash to a runtime ID.
	stateToRuntimeID = make(map[StateHash]uint32)
	// runtimeIDToState maps a runtime ID to a state.
	runtimeIDToState = make(map[uint32]blockupgrader.BlockState)

	// itemRuntimeIDsToNames holds a map to translate item runtime IDs to string IDs.
	itemRuntimeIDsToNames = make(map[int32]string)
	// itemNamesToRuntimeIDs holds a map to translate item string IDs to runtime IDs.
	itemNamesToRuntimeIDs = make(map[string]int32)
)

// init initializes the item and block state mappings.
func init() {
	dec := nbt.NewDecoder(bytes.NewBuffer(BlockStateData))

	// Register all block states present in the block_states.nbt file. These are all possible options registered
	// blocks may encode to.
	var s blockupgrader.BlockState
	var rid uint32
	for {
		if err := dec.Decode(&s); err != nil {
			break
		}

		upgraded := blockupgrader.Upgrade(s)
		stateToRuntimeID[HashState(upgraded)] = rid
		runtimeIDToState[rid] = s
		rid++
	}

	var m map[string]int32
	err := nbt.Unmarshal(ItemRuntimeIDData, &m)
	if err != nil {
		panic(err)
	}
	for name, rid := range m {
		itemNamesToRuntimeIDs[name] = rid
		itemRuntimeIDsToNames[rid] = name
	}
}

// StateToRuntimeID converts a name and its state properties to a runtime ID.
func StateToRuntimeID(name string, properties map[string]any) (runtimeID uint32, found bool) {
	upgraded := blockupgrader.Upgrade(blockupgrader.BlockState{Name: name, Properties: properties})
	rid, ok := stateToRuntimeID[HashState(upgraded)]
	return rid, ok
}

// RuntimeIDToState converts a runtime ID to a name and its state properties.
func RuntimeIDToState(runtimeID uint32) (name string, properties map[string]any, found bool) {
	s := runtimeIDToState[runtimeID]
	return s.Name, s.Properties, true
}

// ItemRuntimeIDToName converts an item runtime ID to a string ID.
func ItemRuntimeIDToName(runtimeID int32) (name string, found bool) {
	name, ok := itemRuntimeIDsToNames[runtimeID]
	return name, ok
}

// ItemNameToRuntimeID converts a string ID to an item runtime ID.
func ItemNameToRuntimeID(name string) (runtimeID int32, found bool) {
	rid, ok := itemNamesToRuntimeIDs[name]
	return rid, ok
}
