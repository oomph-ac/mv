package mappings

import (
	"bytes"
	_ "embed"

	"github.com/df-mc/worldupgrader/blockupgrader"
	"github.com/oomph-ac/mv/multiversion/latest"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// MVBlockMapping holds all data blocks related.
type MVBlockMapping struct {
	// blocks holds a list of all existing v in the game.
	blocks []protocol.BlockEntry
	// stateToRuntimeID maps a block state hash to a runtime ID.
	stateToRuntimeID map[latest.StateHash]uint32
	// runtimeIDToState maps a runtime ID to a state.
	runtimeIDToState map[uint32]blockupgrader.BlockState
	// LegacyAirRID is the runtime ID of the air block of that mapping.
	LegacyAirRID uint32
}

// blockMapping returns MVBlockMapping instance of all block entries and values in the maps from the resource JSON.
func blockMapping(blockStateData []byte) MVBlockMapping {
	dec := nbt.NewDecoder(bytes.NewBuffer(blockStateData))

	// Register all block states present in the block_states.nbt file. These are all possible options registered
	// blocks may encode to.
	var s blockupgrader.BlockState
	var blocks []protocol.BlockEntry
	var stateToRuntimeID = make(map[latest.StateHash]uint32)
	var runtimeIDToState = make(map[uint32]blockupgrader.BlockState)

	for {
		if err := dec.Decode(&s); err != nil {
			break
		}

		s = blockupgrader.Upgrade(s)
		rid := uint32(len(blocks))
		blocks = append(blocks, protocol.BlockEntry{
			Name:       s.Name,
			Properties: s.Properties,
		})

		stateToRuntimeID[latest.HashState(s)] = rid
		runtimeIDToState[rid] = s
	}

	mappings := MVBlockMapping{
		blocks:           blocks,
		stateToRuntimeID: stateToRuntimeID,
		runtimeIDToState: runtimeIDToState,
	}
	mappings.LegacyAirRID = mappings.StateToRuntimeID("minecraft:air", nil)

	return mappings
}

// StateToRuntimeID converts a name and its state properties to a runtime ID.
func (m MVBlockMapping) StateToRuntimeID(name string, properties map[string]any) uint32 {
	rid, ok := m.stateToRuntimeID[latest.HashState(blockupgrader.BlockState{Name: name, Properties: properties})]
	if !ok {
		rid = m.stateToRuntimeID[latest.HashState(blockupgrader.BlockState{Name: "minecraft:info_update"})]
	}

	return rid
}

// RuntimeIDToState converts a runtime ID to a name and its state properties.
func (m MVBlockMapping) RuntimeIDToState(runtimeID uint32) (name string, properties map[string]any, found bool) {
	s := m.runtimeIDToState[runtimeID]
	return s.Name, s.Properties, true
}

// Blocks returns a slice of all block entries.
func (m MVBlockMapping) Blocks() []protocol.BlockEntry {
	return m.blocks
}
