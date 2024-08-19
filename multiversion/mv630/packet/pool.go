package packet

import (
	v649packet "github.com/oomph-ac/mv/multiversion/mv649/packet"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// NewClientPool returns a new pool containing packets sent by a client.
// Packets may be retrieved from it simply by indexing it with the packet ID.
func NewClientPool() packet.Pool {
	pool := v649packet.NewClientPool()

	pool[packet.IDPlayerAuthInput] = func() packet.Packet { return &PlayerAuthInput{} }

	delete(pool, packet.IDSetHud)

	return pool
}

// NewServerPool returns a new pool containing packets sent by a server.
// Packets may be retrieved from it simply by indexing it with the packet ID.
func NewServerPool() packet.Pool {
	pool := v649packet.NewServerPool()

	delete(pool, packet.IDSetHud)

	pool[packet.IDLevelChunk] = func() packet.Packet { return &LevelChunk{} }
	pool[packet.IDPlayerList] = func() packet.Packet { return &PlayerList{} }

	return pool
}
