package mv618

import (
	"github.com/oomph-ac/mv/multiversion/mv618/packet"
	"github.com/oomph-ac/mv/multiversion/mv622"
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

func (Protocol) Encryption(key [32]byte) gtpacket.Encryption {
	return gtpacket.NewCTREncryption(key[:])
}

func (Protocol) ConvertToLatest(pk gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	if upgraded, ok := util.DefaultUpgrade(conn, pk, Mapping); ok {
		if upgraded == nil {
			return []gtpacket.Packet{}
		}

		return []gtpacket.Packet{upgraded}
	}

	return []gtpacket.Packet{pk}
}

func (Protocol) ConvertFromLatest(pk gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	if downgraded, ok := util.DefaultDowngrade(conn, pk, Mapping); ok {
		return []gtpacket.Packet{downgraded}
	}

	return Downgrade([]gtpacket.Packet{pk}, conn)
}

func Downgrade(pks []gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	packets := []gtpacket.Packet{}
	for _, pk := range mv622.Downgrade(pks, conn) {
		switch pk := pk.(type) {
		case *gtpacket.Disconnect:
			packets = append(packets, &packet.Disconnect{
				HideDisconnectionScreen: pk.HideDisconnectionScreen,
				Message:                 pk.Message,
			})
		default:
			packets = append(packets, pk)
		}
	}

	pks = nil
	return packets
}
