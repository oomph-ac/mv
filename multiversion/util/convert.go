package util

import (
	"bytes"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/oomph-ac/mv/multiversion/chunk"
	"github.com/oomph-ac/mv/multiversion/latest"
	"github.com/oomph-ac/mv/multiversion/mappings"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sirupsen/logrus"
)

var overWorldRange = cube.Range{-64, 320}

// LatestAirRID is the runtime ID of the air block in the latest version of the game.
var LatestAirRID, _ = latest.StateToRuntimeID("minecraft:air", nil)

// DefaultUpgrade translates a packet from the legacy version to the latest version.
func DefaultUpgrade(_ *minecraft.Conn, pk packet.Packet, mapping mappings.MVMapping) (packet.Packet, bool) {
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
			data.HeldItem.Stack = UpgradeItem(data.HeldItem.Stack, mapping)

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
	case *packet.LevelChunk:
		if pk.SubChunkCount == protocol.SubChunkRequestModeLimited || pk.SubChunkCount == protocol.SubChunkRequestModeLimitless {
			return pk, true
		}

		buff := bytes.NewBuffer(pk.RawPayload)
		c, err := chunk.NetworkDecode(mapping.LegacyAirRID, buff, int(pk.SubChunkCount), false, overWorldRange)
		if err != nil {
			logrus.Error(err)
			return pk, true
		}

		newChunk := chunk.New(LatestAirRID, overWorldRange)
		for si, sub := range c.Sub() {
			for li, layer := range sub.Layers() {
				upgradedLayer := newChunk.Sub()[si].Layer(uint8(li))
				for x := uint8(0); x < 16; x++ {
					for z := uint8(0); z < 16; z++ {
						for y := uint8(0); y < 16; y++ {
							upgradedLayer.Set(x, y, z, UpgradeBlockRuntimeID(layer.At(x, y, z), mapping))
						}
					}
				}
			}
		}
		for x := uint8(0); x < 16; x++ {
			for z := uint8(0); z < 16; z++ {
				y := c.HighestBlock(x, z)
				newChunk.SetBiome(x, y, z, c.Biome(x, y, z))
			}
		}

		data := chunk.Encode(newChunk, chunk.NetworkEncoding, overWorldRange)
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
				var index byte = 0
				subChunk, err := chunk.DecodeSubChunk(mapping.LegacyAirRID, overWorldRange, buff, &index, chunk.NetworkEncoding)
				if err != nil {
					logrus.Error(err)
					return pk, true
				}

				newSub := chunk.NewSubChunk(LatestAirRID)
				for i, layer := range subChunk.Layers() {
					newLayer := newSub.Layer(uint8(i))
					for x := uint8(0); x < 16; x++ {
						for z := uint8(0); z < 16; z++ {
							for y := uint8(0); y < 16; y++ {
								newLayer.Set(x, y, z, UpgradeBlockRuntimeID(layer.At(x, y, z), mapping))
							}
						}
					}
				}
				newSub.Compact()

				pk.SubChunkEntries[i].RawPayload = append(chunk.EncodeSubChunk(newSub, chunk.NetworkEncoding, overWorldRange, int(index)), buff.Bytes()...)
			}
		}
	case *packet.UpdateBlock:
		pk.NewBlockRuntimeID = UpgradeBlockRuntimeID(pk.NewBlockRuntimeID, mapping)
	case *packet.UpdateBlockSynced:
		pk.NewBlockRuntimeID = UpgradeBlockRuntimeID(pk.NewBlockRuntimeID, mapping)
	case *packet.UpdateSubChunkBlocks:
		for i, block := range pk.Blocks {
			pk.Blocks[i].BlockRuntimeID = UpgradeBlockRuntimeID(block.BlockRuntimeID, mapping)
		}
		for i, block := range pk.Extra {
			pk.Blocks[i].BlockRuntimeID = UpgradeBlockRuntimeID(block.BlockRuntimeID, mapping)
		}
	case *packet.ClientCacheMissResponse:
		for i, blob := range pk.Blobs {
			if blob.Payload[0] != chunk.SubChunkVersion {
				continue
			}

			var index byte = 0
			sub, err := chunk.DecodeSubChunk(mapping.LegacyAirRID, overWorldRange, bytes.NewBuffer(blob.Payload), &index, chunk.NetworkEncoding)
			if err != nil {
				logrus.Error(err)
				return pk, true
			}

			newSub := chunk.NewSubChunk(LatestAirRID)
			for li, layer := range sub.Layers() {
				newLayer := newSub.Layer(uint8(li))
				for x := uint8(0); x < 16; x++ {
					for z := uint8(0); z < 16; z++ {
						for y := uint8(0); y < 16; y++ {
							newLayer.Set(x, y, z, UpgradeBlockRuntimeID(layer.At(x, y, z), mapping))
						}
					}
				}
			}
			newSub.Compact()

			pk.Blobs[i].Payload = chunk.EncodeSubChunk(newSub, chunk.NetworkEncoding, overWorldRange, int(index))
		}
	case *packet.SetHud:
		break
	default:
		if pk.ID() == 53 {
			return pk, true
		}
		handled = false
	}

	return pk, handled
}

// DefaultDowngrade translates a packet from the latest version to the legacy version.
func DefaultDowngrade(_ *minecraft.Conn, pk packet.Packet, mapping mappings.MVMapping) (packet.Packet, bool) {
	handled := true
	switch pk := pk.(type) {
	case *packet.MobEquipment:
		pk.NewItem.Stack = DowngradeItem(pk.NewItem.Stack, mapping)
	case *packet.MobArmourEquipment:
		pk.Leggings.Stack = DowngradeItem(pk.Leggings.Stack, mapping)
		pk.Boots.Stack = DowngradeItem(pk.Boots.Stack, mapping)
		pk.Helmet.Stack = DowngradeItem(pk.Helmet.Stack, mapping)
		pk.Chestplate.Stack = DowngradeItem(pk.Chestplate.Stack, mapping)
	case *packet.SetActorData:
		variant, ok := pk.EntityMetadata[protocol.EntityDataKeyVariant]
		if ok {
			pk.EntityMetadata[protocol.EntityDataKeyVariant] = int32(DowngradeBlockRuntimeID(uint32(variant.(int32)), mapping))
		}
	case *packet.AddActor:
		variant, ok := pk.EntityMetadata[protocol.EntityDataKeyVariant]
		if ok {
			pk.EntityMetadata[protocol.EntityDataKeyVariant] = int32(DowngradeBlockRuntimeID(uint32(variant.(int32)), mapping))
		}
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
	case *packet.CraftingData: // TODO: Fix crafting later, this keeps crashing the client.
		return &packet.CraftingData{
			ClearRecipes: true,
		}, true
	case *packet.StartGame:
		items := make([]protocol.ItemEntry, 0, len(pk.Items))
		for _, item := range pk.Items {
			id, ok := latest.ItemNameToRuntimeID(item.Name)
			if !ok {
				items = append(items, item)
				continue
			}

			name, ok := mapping.ItemNameByID(id)
			if !ok {
				items = append(items, item)
				continue
			}

			items = append(items, protocol.ItemEntry{
				Name:           name,
				RuntimeID:      item.RuntimeID,
				ComponentBased: item.ComponentBased,
			})
		}

		pk.Items = items
	case *packet.LevelChunk:

		if pk.SubChunkCount == protocol.SubChunkRequestModeLimited || pk.SubChunkCount == protocol.SubChunkRequestModeLimitless {
			return pk, true
		}

		buff := bytes.NewBuffer(pk.RawPayload)
		c, err := chunk.NetworkDecode(mapping.LegacyAirRID, buff, int(pk.SubChunkCount), false, overWorldRange)
		if err != nil {
			logrus.Error(err)
			return pk, true
		}

		newChunk := chunk.New(LatestAirRID, overWorldRange)
		for si, sub := range c.Sub() {
			for li, layer := range sub.Layers() {
				newLayer := newChunk.Sub()[si].Layer(uint8(li))
				for x := uint8(0); x < 16; x++ {
					for z := uint8(0); z < 16; z++ {
						for y := uint8(0); y < 16; y++ {
							newLayer.Set(x, y, z, DowngradeBlockRuntimeID(layer.At(x, y, z), mapping))
						}
					}
				}
			}
		}
		for x := uint8(0); x < 16; x++ {
			for z := uint8(0); z < 16; z++ {
				y := c.HighestBlock(x, z)
				newChunk.SetBiome(x, y, z, c.Biome(x, y, z))
			}
		}

		data := chunk.Encode(newChunk, chunk.NetworkEncoding, overWorldRange)
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
				var index byte = 0
				sub, err := chunk.DecodeSubChunk(LatestAirRID, overWorldRange, buff, &index, chunk.NetworkEncoding)
				if err != nil {
					logrus.Error(err)
					return pk, true
				}

				newSub := chunk.NewSubChunk(mapping.LegacyAirRID)
				for li, layer := range sub.Layers() {
					newLayer := newSub.Layer(uint8(li))
					for x := uint8(0); x < 16; x++ {
						for z := uint8(0); z < 16; z++ {
							for y := uint8(0); y < 16; y++ {
								newLayer.Set(x, y, z, DowngradeBlockRuntimeID(layer.At(x, y, z), mapping))
							}
						}
					}
				}
				newSub.Compact()

				pk.SubChunkEntries[i].RawPayload = append(chunk.EncodeSubChunk(newSub, chunk.NetworkEncoding, overWorldRange, int(index)), buff.Bytes()...)
			}
		}
	default:
		handled = false
	}

	return pk, handled
}
