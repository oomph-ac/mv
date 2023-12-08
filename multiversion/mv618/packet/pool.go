package packet

import (
	"github.com/oomph-ac/mv/multiversion/mv622/packet"
	gtpacket "github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func NewClientPool() gtpacket.Pool {
	pool := packet.NewClientPool()
	return pool
}

func NewServerPool() gtpacket.Pool {
	pool := packet.NewServerPool()
	pool[gtpacket.IDDisconnect] = func() gtpacket.Packet { return &Disconnect{} }
	return pool
}
