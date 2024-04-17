package packet

import (
	"github.com/oomph-ac/mv/multiversion/mv649/packet"
	gtpacket "github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func NewClientPool() gtpacket.Pool {
	pool := packet.NewClientPool()
	pool[gtpacket.IDPlayerAuthInput] = func() gtpacket.Packet { return &PlayerAuthInput{} }
	delete(pool, gtpacket.IDSetHud)
	return pool
}

func NewServerPool() gtpacket.Pool {
	pool := packet.NewServerPool()
	delete(pool, gtpacket.IDSetHud)
	pool[gtpacket.IDLevelChunk] = func() gtpacket.Packet { return &LevelChunk{} }
	pool[gtpacket.IDPlayerList] = func() gtpacket.Packet { return &PlayerList{} }
	return pool
}
