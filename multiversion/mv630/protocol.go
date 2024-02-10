package mv630

import (
	"github.com/oomph-ac/mv/multiversion/mv630/packet"
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

func (Protocol) ConvertToLatest(pk gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	return []gtpacket.Packet{pk}
}

func (Protocol) ConvertFromLatest(pk gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	return []gtpacket.Packet{pk}
}

func Downgrade(pks []gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	packets := []gtpacket.Packet{}
	for _, pk := range pks {
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
		default:
			packets = append(packets, pk)
		}
	}

	pks = nil
	return packets
}
