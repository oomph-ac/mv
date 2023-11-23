package util

import (
	"bytes"
	"fmt"

	"github.com/df-mc/dragonfly/server/world"
	"github.com/oomph-ac/mv/internal/pointer"
	"github.com/oomph-ac/mv/multiversion/chunk"
	"github.com/oomph-ac/mv/multiversion/latest"
	"github.com/oomph-ac/mv/multiversion/mappings"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
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

// DowngradeBlockPacket translates a block packet from the latest version to the legacy version.
func DowngradeBlockPacket(conn *minecraft.Conn, pk packet.Packet, mapping mappings.MVMapping) (packet.Packet, bool) {
	handled := true
	switch pk := pk.(type) {
	case *packet.LevelChunk:
		if pk.SubChunkCount == protocol.SubChunkRequestModeLimited || pk.SubChunkCount == protocol.SubChunkRequestModeLimitless {
			return pk, true
		}

		r := world.Overworld.Range()
		buff := bytes.NewBuffer(pk.RawPayload)
		c, err := chunk.NetworkDecode(LatestAirRID, buff, int(pk.SubChunkCount), conn.GameData().BaseGameVersion == "1.17.40", r)
		if err != nil {
			fmt.Println(err)
			return pk, true
		}

		downgraded := chunk.New(mapping.LegacyAirRID, r)
		for subInd, sub := range c.Sub() {
			for layerInd, layer := range sub.Layers() {
				downgradedLayer := downgraded.Sub()[subInd].Layer(uint8(layerInd))
				for x := uint8(0); x < 16; x++ {
					for z := uint8(0); z < 16; z++ {
						for y := uint8(0); y < 16; y++ {
							latestRuntimeID := layer.At(x, y, z)
							downgradedLayer.Set(x, y, z, DowngradeBlockRuntimeID(latestRuntimeID, mapping))
						}
					}
				}
			}
		}
		for x := uint8(0); x < 16; x++ {
			for z := uint8(0); z < 16; z++ {
				y := c.HighestBlock(x, z)
				downgraded.SetBiome(x, y, z, c.Biome(x, y, z))
			}
		}

		data := chunk.Encode(downgraded, chunk.NetworkEncoding, r)
		chunkBuf := bytes.NewBuffer(nil)
		for i := range data.SubChunks {
			chunkBuf.Write(data.SubChunks[i])
		}
		chunkBuf.Write(data.Biomes)

		pk.SubChunkCount = uint32(len(data.SubChunks))
		pk.RawPayload = append(chunkBuf.Bytes(), buff.Bytes()...)
	case *packet.SubChunk:
		for i, entry := range pk.SubChunkEntries {
			if entry.Result == protocol.SubChunkResultSuccess && !pk.CacheEnabled {
				buff := bytes.NewBuffer(entry.RawPayload)
				subChunk, err := chunk.DecodeSubChunk(LatestAirRID, world.Overworld.Range(), buff, pointer.Make(uint8(0)), chunk.NetworkEncoding)
				if err != nil {
					fmt.Println(err)
					return pk, true
				}

				downgraded := chunk.NewSubChunk(mapping.LegacyAirRID)
				for layerInd, layer := range subChunk.Layers() {
					downgradedLayer := downgraded.Layer(uint8(layerInd))
					for x := uint8(0); x < 16; x++ {
						for z := uint8(0); z < 16; z++ {
							for y := uint8(0); y < 16; y++ {
								latestRuntimeID := layer.At(x, y, z)
								downgradedLayer.Set(x, y, z, DowngradeBlockRuntimeID(latestRuntimeID, mapping))
							}
						}
					}
				}
				ind := int16(pk.Position.Y()) + int16(entry.Offset[1]) - int16(world.Overworld.Range()[0])>>4
				serialised := chunk.EncodeSubChunk(downgraded, chunk.NetworkEncoding, world.Overworld.Range(), int(ind))
				pk.SubChunkEntries[i].RawPayload = append(serialised, buff.Bytes()...)
			}
		}
	case *packet.UpdateBlock:
		pk.NewBlockRuntimeID = DowngradeBlockRuntimeID(pk.NewBlockRuntimeID, mapping)
	default:
		handled = false
	}

	return pk, handled
}
