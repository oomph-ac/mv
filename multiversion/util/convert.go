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

// DefaultUpgrade translates a packet from the legacy version to the latest version.
func DefaultUpgrade(conn *minecraft.Conn, pk packet.Packet, mapping mappings.MVMapping) (packet.Packet, bool) {
	handled := true
	switch pk := pk.(type) {
	case *packet.InventoryTransaction:
		for i, action := range pk.Actions {
			pk.Actions[i].OldItem.Stack = UpgradeItem(action.OldItem.Stack, mapping)
			pk.Actions[i].NewItem.Stack = UpgradeItem(action.NewItem.Stack, mapping)
		}
		switch data := pk.TransactionData.(type) {
		case *protocol.UseItemTransactionData:
			if data.BlockRuntimeID > 0 {
				data.BlockRuntimeID = UpgradeBlockRuntimeID(data.BlockRuntimeID, mapping)
			}
			pk.TransactionData = data
		case *protocol.UseItemOnEntityTransactionData:
			data.HeldItem.Stack = UpgradeItem(data.HeldItem.Stack, mapping)
			pk.TransactionData = data
		case *protocol.ReleaseItemTransactionData:
			data.HeldItem.Stack = UpgradeItem(data.HeldItem.Stack, mapping)
			pk.TransactionData = data
		}
	case *packet.ItemStackRequest:
		for i, request := range pk.Requests {
			var actions = make([]protocol.StackRequestAction, 0)
			for _, action := range request.Actions {
				switch data := action.(type) {
				case *protocol.CraftResultsDeprecatedStackRequestAction:
					for k, item := range data.ResultItems {
						data.ResultItems[k] = UpgradeItem(item, mapping)
					}
					action = data
				}
				actions = append(actions, action)
			}
			pk.Requests[i].Actions = actions
		}
	case *packet.MobArmourEquipment:
		pk.Helmet.Stack = UpgradeItem(pk.Helmet.Stack, mapping)
		pk.Chestplate.Stack = UpgradeItem(pk.Chestplate.Stack, mapping)
		pk.Leggings.Stack = UpgradeItem(pk.Leggings.Stack, mapping)
		pk.Boots.Stack = UpgradeItem(pk.Boots.Stack, mapping)
	case *packet.MobEquipment:
		pk.NewItem.Stack = UpgradeItem(pk.NewItem.Stack, mapping)
	default:
		if pk.ID() == 53 {
			return nil, true
		}

		handled = false
	}

	return pk, handled
}

// DefaultDowngrade translates a packet from the latest version to the legacy version.
func DefaultDowngrade(conn *minecraft.Conn, pk packet.Packet, mapping mappings.MVMapping) (packet.Packet, bool) {
	handled := true
	switch pk := pk.(type) {
	case *packet.AddItemActor:
		pk.Item.Stack = DowngradeItem(pk.Item.Stack, mapping)
	case *packet.AddPlayer:
		pk.HeldItem.Stack = DowngradeItem(pk.HeldItem.Stack, mapping)
	case *packet.CreativeContent:
		for i, item := range pk.Items {
			pk.Items[i].Item = DowngradeItem(item.Item, mapping)
		}
	case *packet.InventoryContent:
		for i, item := range pk.Content {
			pk.Content[i].Stack = DowngradeItem(item.Stack, mapping)
		}
	case *packet.InventorySlot:
		pk.NewItem.Stack = DowngradeItem(pk.NewItem.Stack, mapping)
	case *packet.LevelEvent:
		if pk.EventType == packet.LevelEventParticlesDestroyBlock || pk.EventType == packet.LevelEventParticlesCrackBlock {
			pk.EventData = int32(DowngradeBlockRuntimeID(uint32(pk.EventData), mapping))
		}
	case *packet.LevelSoundEvent:
		if pk.SoundType == packet.SoundEventPlace || pk.SoundType == packet.SoundEventHit || pk.SoundType == packet.SoundEventItemUseOn || pk.SoundType == packet.SoundEventLand {
			pk.ExtraData = int32(DowngradeBlockRuntimeID(uint32(pk.ExtraData), mapping))
		}
	case *packet.LevelChunk:
		if pk.SubChunkCount == protocol.SubChunkRequestModeLimited || pk.SubChunkCount == protocol.SubChunkRequestModeLimitless {
			return nil, false
		}

		r := world.Overworld.Range()
		buff := bytes.NewBuffer(pk.RawPayload)
		c, err := chunk.NetworkDecode(LatestAirRID, buff, int(pk.SubChunkCount), conn.GameData().BaseGameVersion == "1.17.40", r)
		if err != nil {
			fmt.Println(err)
			return nil, false
		}
		downgraded := chunk.New(mapping.LegacyAirRID, r)
		for subInd, sub := range c.Sub() {
			for layerInd, layer := range sub.Layers() {
				downgradedLayer := downgraded.Sub()[subInd].Layer(uint8(layerInd))
				for x := uint8(0); x < 16; x++ {
					for z := uint8(0); z < 16; z++ {
						for y := uint8(0); y < 16; y++ {
							latestRuntimeID := layer.At(x, y, z)
							if latestRuntimeID == LatestAirRID {
								// Don't bother with air.
								continue
							}
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
			_, _ = chunkBuf.Write(data.SubChunks[i])
		}
		_, _ = chunkBuf.Write(data.Biomes)

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
								if latestRuntimeID == LatestAirRID {
									// Don't bother with air.
									continue
								}
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
	case *packet.UpdateBlockSynced:
		pk.NewBlockRuntimeID = DowngradeBlockRuntimeID(pk.NewBlockRuntimeID, mapping)
	case *packet.UpdateSubChunkBlocks:
		for i, block := range pk.Blocks {
			pk.Blocks[i].BlockRuntimeID = DowngradeBlockRuntimeID(block.BlockRuntimeID, mapping)
		}
		for i, block := range pk.Extra {
			pk.Blocks[i].BlockRuntimeID = DowngradeBlockRuntimeID(block.BlockRuntimeID, mapping)
		}
	default:
		handled = false
	}

	return pk, handled
}
