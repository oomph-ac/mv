package mv630

import (
	"github.com/oomph-ac/mv/multiversion/mv630/packet"
	v630protocol "github.com/oomph-ac/mv/multiversion/mv630/protocol"
	"github.com/oomph-ac/mv/multiversion/mv649"
	v649packet "github.com/oomph-ac/mv/multiversion/mv649/packet"
	"github.com/oomph-ac/mv/multiversion/util"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	gtpacket "github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// Protocol is a protocol implementation for
// Minecraft: Bedrock Edition 1.20.50.
type Protocol struct{}

// ID ...
func (Protocol) ID() int32 {
	return 630
}

// Ver ...
func (Protocol) Ver() string {
	return "1.20.50"
}

// NewReader ...
func (Protocol) NewReader(r minecraft.ByteReader, shieldID int32, enableLimits bool) protocol.IO {
	return protocol.NewReader(r, shieldID, enableLimits)
}

// NewWriter ...
func (Protocol) NewWriter(w minecraft.ByteWriter, shieldID int32) protocol.IO {
	return protocol.NewWriter(w, shieldID)
}

// Packets ...
func (Protocol) Packets(listener bool) gtpacket.Pool {
	if listener {
		return packet.NewClientPool()
	}
	return packet.NewServerPool()
}

// ConvertToLatest ...
func (Protocol) ConvertToLatest(pk gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	if upgraded, ok := util.DefaultUpgrade(conn, pk, Mapping); ok {
		if upgraded == nil {
			return []gtpacket.Packet{}
		}

		return Upgrade([]gtpacket.Packet{upgraded}, conn)
	}

	return Upgrade([]gtpacket.Packet{pk}, conn)
}

// ConvertFromLatest ...
func (Protocol) ConvertFromLatest(pk gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	if downgraded, ok := util.DefaultDowngrade(conn, pk, Mapping); ok {
		return Downgrade([]gtpacket.Packet{downgraded}, conn)
	}

	return Downgrade([]gtpacket.Packet{pk}, conn)
}

// Upgrade ...
func Upgrade(pks []gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	packets := make([]gtpacket.Packet, 0, len(pks))
	for _, pk := range pks {
		switch pk := pk.(type) {
		case *packet.PlayerAuthInput:
			packets = append(packets, &v649packet.PlayerAuthInput{
				Pitch:                  pk.Pitch,
				Yaw:                    pk.Yaw,
				Position:               pk.Position,
				MoveVector:             pk.MoveVector,
				HeadYaw:                pk.HeadYaw,
				InputData:              pk.InputData,
				InputMode:              pk.InputMode,
				PlayMode:               pk.PlayMode,
				InteractionModel:       pk.InteractionModel,
				GazeDirection:          pk.GazeDirection,
				Tick:                   pk.Tick,
				Delta:                  pk.Delta,
				ItemInteractionData:    pk.ItemInteractionData,
				ItemStackRequest:       pk.ItemStackRequest,
				BlockActions:           pk.BlockActions,
				AnalogueMoveVector:     pk.AnalogueMoveVector,
				ClientPredictedVehicle: 0,
			})
		default:
			packets = append(packets, pk)
		}
	}

	return mv649.Upgrade(packets, conn)
}

// Downgrade ...
func Downgrade(pks []gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	downgraded := mv649.Downgrade(pks, conn)
	packets := make([]gtpacket.Packet, 0, len(downgraded))

	for _, pk := range downgraded {
		switch pk := pk.(type) {
		case *gtpacket.LevelChunk:
			packets = append(packets, &packet.LevelChunk{
				Position:        pk.Position,
				HighestSubChunk: pk.HighestSubChunk,
				SubChunkCount:   pk.SubChunkCount,
				CacheEnabled:    pk.CacheEnabled,
				BlobHashes:      pk.BlobHashes,
				RawPayload:      pk.RawPayload,
			})
		case *gtpacket.PlayerList:
			packets = append(packets, &packet.PlayerList{
				Entries: downgradePlayerListEntries(pk.Entries),
			})
		default:
			packets = append(packets, pk)
		}
	}
	return packets
}

// downgradePlayerListEntries downgrades a slice of protocol.PlayerListEntry.
func downgradePlayerListEntries(entries []protocol.PlayerListEntry) []v630protocol.PlayerListEntry {
	newEntries := make([]v630protocol.PlayerListEntry, 0, len(entries))
	for _, e := range entries {
		newEntries = append(newEntries, v630protocol.PlayerListEntry{
			UUID:           e.UUID,
			EntityUniqueID: e.EntityUniqueID,
			Username:       e.Username,
			XUID:           e.XUID,
			PlatformChatID: e.PlatformChatID,
			BuildPlatform:  e.BuildPlatform,
			Skin:           e.Skin,
			Teacher:        e.Teacher,
			Host:           e.Host,
		})
	}
	return newEntries
}
