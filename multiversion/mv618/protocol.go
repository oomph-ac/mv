package mv618

import (
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
	packets := util.DowngradeBlockPackets(conn, pk, Mapping)

	for i, pk := range packets {
		dc, ok := pk.(*gtpacket.Disconnect)
		if !ok {
			continue
		}

		packets[i] = &packet.Disconnect{
			Message:                 dc.Message,
			HideDisconnectionScreen: dc.HideDisconnectionScreen,
		}
	}

	packets = append(packets, pk)
	return packets
}
