package packet

import (
	v594packet "github.com/oomph-ac/mv/multiversion/mv594/packet"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func NewClientPool() packet.Pool {
	pool := v594packet.NewClientPool()
	return pool
}

func NewServerPool() packet.Pool {
	pool := v594packet.NewServerPool()
	pool[packet.IDAvailableCommands] = func() packet.Packet { return &AvailableCommands{} }

	return pool
}
