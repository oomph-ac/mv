package mv618

import (
	"bytes"
	"fmt"

	"github.com/df-mc/dragonfly/server/world"
	"github.com/oomph-ac/mv/internal/pointer"
	"github.com/oomph-ac/mv/multiversion/chunk"
	"github.com/oomph-ac/mv/multiversion/mv618/packet"
	"github.com/oomph-ac/mv/multiversion/util"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	gtpacket "github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type Protocol struct{}

func (Protocol) ID() int32 {
	return 618
}

func (Protocol) Ver() string {
	return "1.20.30"
}

func (Protocol) NewReader(r minecraft.ByteReader, shieldID int32, enableLimits bool) protocol.IO {
	return protocol.NewReader(r, shieldID, enableLimits)
}

func (Protocol) NewWriter(w minecraft.ByteWriter, shieldID int32) protocol.IO {
	return protocol.NewWriter(w, shieldID)
}

func (Protocol) Packets(listener bool) gtpacket.Pool {
	if listener {
		return packet.NewClientPool()
	}
	return packet.NewServerPool()
}

func (Protocol) ConvertToLatest(pk gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	return []gtpacket.Packet{pk}
}

func (Protocol) ConvertFromLatest(pk gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	packets := []gtpacket.Packet{}
	switch pk := pk.(type) {
	case *gtpacket.Disconnect:
		packets = append(packets, &packet.Disconnect{
			HideDisconnectionScreen: pk.HideDisconnectionScreen,
			Message:                 pk.Message,
		})
	case *gtpacket.LevelChunk:
		if pk.SubChunkCount == protocol.SubChunkRequestModeLimited || pk.SubChunkCount == protocol.SubChunkRequestModeLimitless {
			packets = append(packets, pk)
			return packets
		}

		r := world.Overworld.Range()
		buff := bytes.NewBuffer(pk.RawPayload)
		c, err := chunk.NetworkDecode(util.LatestAirRID, buff, int(pk.SubChunkCount), conn.GameData().BaseGameVersion == "1.17.40", r)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		downgraded := chunk.New(Mapping.LegacyAirRID, r)
		for subInd, sub := range c.Sub() {
			for layerInd, layer := range sub.Layers() {
				downgradedLayer := downgraded.Sub()[subInd].Layer(uint8(layerInd))
				for x := uint8(0); x < 16; x++ {
					for z := uint8(0); z < 16; z++ {
						for y := uint8(0); y < 16; y++ {
							latestRuntimeID := layer.At(x, y, z)
							downgradedLayer.Set(x, y, z, util.DowngradeBlockRuntimeID(latestRuntimeID, Mapping))
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
	case *gtpacket.SubChunk:
		for i, entry := range pk.SubChunkEntries {
			if entry.Result == protocol.SubChunkResultSuccess && !pk.CacheEnabled {
				buff := bytes.NewBuffer(entry.RawPayload)
				subChunk, err := chunk.DecodeSubChunk(util.LatestAirRID, world.Overworld.Range(), buff, pointer.Make(uint8(0)), chunk.NetworkEncoding)
				if err != nil {
					fmt.Println(err)
					return nil
				}
				downgraded := chunk.NewSubChunk(Mapping.LegacyAirRID)
				for layerInd, layer := range subChunk.Layers() {
					downgradedLayer := downgraded.Layer(uint8(layerInd))
					for x := uint8(0); x < 16; x++ {
						for z := uint8(0); z < 16; z++ {
							for y := uint8(0); y < 16; y++ {
								latestRuntimeID := layer.At(x, y, z)
								downgradedLayer.Set(x, y, z, util.DowngradeBlockRuntimeID(latestRuntimeID, Mapping))
							}
						}
					}
				}
				ind := int16(pk.Position.Y()) + int16(entry.Offset[1]) - int16(world.Overworld.Range()[0])>>4
				serialised := chunk.EncodeSubChunk(downgraded, chunk.NetworkEncoding, world.Overworld.Range(), int(ind))
				pk.SubChunkEntries[i].RawPayload = append(serialised, buff.Bytes()...)
			}
		}
	case *gtpacket.UpdateBlock:
		pk.NewBlockRuntimeID = util.DowngradeBlockRuntimeID(pk.NewBlockRuntimeID, Mapping)
	}

	packets = append(packets, pk)
	return packets
}
