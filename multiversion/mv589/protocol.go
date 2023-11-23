package mv589

import (
	"github.com/oomph-ac/mv/multiversion/mv589/packet"
	"github.com/oomph-ac/mv/multiversion/mv594"
	v594packet "github.com/oomph-ac/mv/multiversion/mv594/packet"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	gtpacket "github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type Protocol struct{}

func (Protocol) ID() int32 {
	return 589
}

func (Protocol) Ver() string {
	return "1.20.0"
}

func (Protocol) NewReader(r minecraft.ByteReader, shieldID int32, enableLimits bool) protocol.IO {
	return protocol.NewReader(r, shieldID, enableLimits)
}

func (Protocol) NewWriter(w minecraft.ByteWriter, shieldID int32) protocol.IO {
	return protocol.NewWriter(w, shieldID)
}

func (Protocol) Packets(listener bool) gtpacket.Pool {
	if listener {
		return v594packet.NewServerPool()
	}
	return v594packet.NewClientPool()
}

func (Protocol) ConvertToLatest(pk gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	return []gtpacket.Packet{pk}
}

func (Protocol) ConvertFromLatest(pk gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	return Downgrade([]gtpacket.Packet{pk}, conn)
}

func Downgrade(pks []gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	packets := mv594.Downgrade(pks, conn)
	for _, pk := range pks {
		switch pk := pk.(type) {
		case *gtpacket.AvailableCommands:
			packets = append(packets, &packet.AvailableCommands{
				EnumValues:   pk.EnumValues,
				Suffixes:     pk.Suffixes,
				Enums:        pk.Enums,
				Commands:     pk.Commands,
				DynamicEnums: pk.DynamicEnums,
				Constraints:  pk.Constraints,
			})
		default:
			packets = append(packets, pk)
		}
	}

	pks = nil
	return packets
}
