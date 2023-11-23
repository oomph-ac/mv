package packet

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

func NewClientPool() packet.Pool {
	pool := packet.NewClientPool()
	return pool
}

func NewServerPool() packet.Pool {
	pool := packet.NewServerPool()
	pool[packet.IDDisconnect] = func() packet.Packet { return &Disconnect{} }
	return pool
}
