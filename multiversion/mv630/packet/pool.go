package packet

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

func NewClientPool() packet.Pool {
	pool := packet.NewClientPool()
	delete(pool, packet.IDSetHud)
	return pool
}

func NewServerPool() packet.Pool {
	pool := packet.NewServerPool()
	delete(pool, packet.IDSetHud)
	return pool
}
