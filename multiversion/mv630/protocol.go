package mv630

import (
	"github.com/oomph-ac/mv/multiversion/mv630/packet"
	"github.com/oomph-ac/mv/multiversion/mv649"
	"github.com/oomph-ac/mv/multiversion/util"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"

	gtpacket "github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type Protocol struct{}

func (Protocol) ID() int32 {
	return 630
}

func (Protocol) Ver() string {
	return "1.20.50"
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

func (Protocol) Encryption(key [32]byte) gtpacket.Encryption {
	return gtpacket.NewCTREncryption(key[:])
}

func (Protocol) ConvertToLatest(pk gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	if upgraded, ok := util.DefaultUpgrade(conn, pk, Mapping); ok {
		if upgraded == nil {
			return []gtpacket.Packet{}
		}

		return Upgrade([]gtpacket.Packet{upgraded}, conn)
	}

	return Upgrade([]gtpacket.Packet{pk}, conn)
}

func (Protocol) ConvertFromLatest(pk gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	if downgraded, ok := util.DefaultDowngrade(conn, pk, Mapping); ok {
		return Downgrade([]gtpacket.Packet{downgraded}, conn)
	}

	return Downgrade([]gtpacket.Packet{pk}, conn)
}

func Upgrade(pks []gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	packets := make([]gtpacket.Packet, 0, len(pks))
	for _, pk := range mv649.Upgrade(pks, conn) {
		switch pk := pk.(type) {
		case *packet.PlayerAuthInput:
			packets = append(packets, &gtpacket.PlayerAuthInput{
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
		case *packet.LevelChunk:
			packets = append(packets, &gtpacket.LevelChunk{
				Position:        pk.Position,
				HighestSubChunk: pk.HighestSubChunk,
				SubChunkCount:   pk.SubChunkCount,
				CacheEnabled:    pk.CacheEnabled,
				BlobHashes:      pk.BlobHashes,
				RawPayload:      pk.RawPayload,
			})
		default:
			packets = append(packets, pk)
		}
	}

	pks = nil
	return packets
}

func Downgrade(pks []gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	packets := make([]gtpacket.Packet, 0, len(pks))
	for _, pk := range mv649.Downgrade(pks, conn) {
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
		case *gtpacket.PlayerAuthInput:
			packets = append(packets, &packet.PlayerAuthInput{
				Pitch:               pk.Pitch,
				Yaw:                 pk.Yaw,
				Position:            pk.Position,
				MoveVector:          pk.MoveVector,
				HeadYaw:             pk.HeadYaw,
				InputData:           pk.InputData,
				InputMode:           pk.InputMode,
				PlayMode:            pk.PlayMode,
				InteractionModel:    pk.InteractionModel,
				GazeDirection:       pk.GazeDirection,
				Tick:                pk.Tick,
				Delta:               pk.Delta,
				ItemInteractionData: pk.ItemInteractionData,
				ItemStackRequest:    pk.ItemStackRequest,
				BlockActions:        pk.BlockActions,
				AnalogueMoveVector:  pk.AnalogueMoveVector,
			})
		case *gtpacket.PlayerList:
			packets = append(packets, &packet.PlayerList{
				Entries: packet.DowngradePlayerEntries(pk.Entries),
			})
		default:
			packets = append(packets, pk)
		}
	}

	pks = nil
	return packets
}
