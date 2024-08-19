package packet

import (
	v686packet "github.com/oomph-ac/mv/multiversion/mv686/packet"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// NewClientPool returns a new pool containing packets sent by a client.
// Packets may be retrieved from it simply by indexing it with the packet ID.
func NewClientPool() packet.Pool {
	pool := v686packet.NewClientPool()

	return pool
}

// NewServerPool returns a new pool containing packets sent by a server.
// Packets may be retrieved from it simply by indexing it with the packet ID.
func NewServerPool() packet.Pool {
	pool := v686packet.NewServerPool()

	delete(pool, packet.IDClientBoundCloseForm)

	return pool
}
