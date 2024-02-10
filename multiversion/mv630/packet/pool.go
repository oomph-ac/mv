package packet

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

func NewClientPool() packet.Pool {
	pool := packet.NewClientPool()
	pool[packet.IDPlayerAuthInput] = func() packet.Packet { return &PlayerAuthInput{} }
	delete(pool, packet.IDSetHud)
	return pool
}

func NewServerPool() packet.Pool {
	pool := packet.NewServerPool()
	delete(pool, packet.IDSetHud)
	pool[packet.IDLevelChunk] = func() packet.Packet { return &LevelChunk{} }
	pool[packet.IDPlayerList] = func() packet.Packet { return &PlayerList{} }
	return pool
}
