package packet

import (
	v685packet "github.com/oomph-ac/mv/multiversion/mv685/packet"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// NewClientPool returns a new pool containing packets sent by a client.
// Packets may be retrieved from it simply by indexing it with the packet ID.
func NewClientPool() packet.Pool {
	pool := v685packet.NewClientPool()

	pool[packet.IDContainerClose] = func() packet.Packet { return &ContainerClose{} }
	pool[packet.IDText] = func() packet.Packet { return &Text{} }
	pool[packet.IDCodeBuilderSource] = func() packet.Packet { return &CodeBuilderSource{} }

	return pool
}

// NewServerPool returns a new pool containing packets sent by a server.
// Packets may be retrieved from it simply by indexing it with the packet ID.
func NewServerPool() packet.Pool {
	pool := v685packet.NewServerPool()

	pool[packet.IDCraftingData] = func() packet.Packet { return &CraftingData{} }
	pool[packet.IDContainerClose] = func() packet.Packet { return &ContainerClose{} }
	pool[packet.IDText] = func() packet.Packet { return &Text{} }
	pool[packet.IDStartGame] = func() packet.Packet { return &StartGame{} }

	return pool
}
