package mv685

import (
	"github.com/oomph-ac/mv/multiversion/mv685/packet"
	"github.com/oomph-ac/mv/multiversion/mv686"
	"github.com/oomph-ac/mv/multiversion/util"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	gtpacket "github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// Protocol is a protocol implementation for
// Minecraft: Bedrock Edition 1.21.0.
type Protocol struct{}

// ID ...
func (Protocol) ID() int32 {
	return 685
}

// Ver ...
func (Protocol) Ver() string {
	return "1.21.0"
}

// NewReader ...
func (Protocol) NewReader(r minecraft.ByteReader, shieldID int32, enableLimits bool) protocol.IO {
	return protocol.NewReader(r, shieldID, enableLimits)
}

// NewWriter ...
func (Protocol) NewWriter(r minecraft.ByteWriter, shieldID int32) protocol.IO {
	return protocol.NewWriter(r, shieldID)
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
	if upgraded, ok := util.DefaultUpgrade(conn, pk, mv686.Mapping); ok {
		return Upgrade([]gtpacket.Packet{upgraded}, conn)
	}
	return Upgrade([]gtpacket.Packet{pk}, conn)
}

// ConvertFromLatest ...
func (Protocol) ConvertFromLatest(pk gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	if downgraded, ok := util.DefaultDowngrade(conn, pk, mv686.Mapping); ok {
		return Downgrade([]gtpacket.Packet{downgraded}, conn)
	}
	return Downgrade([]gtpacket.Packet{pk}, conn)
}

// Upgrade ...
func Upgrade(pks []gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	return mv686.Upgrade(pks, conn)
}

// Downgrade ...
func Downgrade(pks []gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	return mv686.Downgrade(pks, conn)
}
