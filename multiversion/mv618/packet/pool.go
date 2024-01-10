package packet

import (
	v622packet "github.com/oomph-ac/mv/multiversion/mv622/packet"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func NewClientPool() packet.Pool {
	pool := v622packet.NewClientPool()
	return pool
}

func NewServerPool() packet.Pool {
	pool := v622packet.NewServerPool()
	pool[packet.IDDisconnect] = func() packet.Packet { return &Disconnect{} }
	return pool
}
