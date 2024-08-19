package packet

import (
	v671packet "github.com/oomph-ac/mv/multiversion/mv671/packet"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// NOTE: CorrectPlayerMovementPrediction is not included in here, since changes
// to the packet were made late, and it was updated around 1.20.50 (630).

// NewClientPool returns a new pool containing packets sent by a client.
// Packets may be retrieved from it simply by indexing it with the packet ID.
func NewClientPool() packet.Pool {
	pool := v671packet.NewClientPool()

	pool[packet.IDPlayerAuthInput] = func() packet.Packet { return &PlayerAuthInput{} }

	return pool
}

// NewServerPool returns a new pool containing packets sent by a server.
// Packets may be retrieved from it simply by indexing it with the packet ID.
func NewServerPool() packet.Pool {
	pool := v671packet.NewServerPool()

	pool[packet.IDResourcePackStack] = func() packet.Packet { return &ResourcePackStack{} }
	pool[packet.IDStartGame] = func() packet.Packet { return &StartGame{} }
	pool[packet.IDCraftingData] = func() packet.Packet { return &CraftingData{} }
	pool[packet.IDUpdateBlockSynced] = func() packet.Packet { return &UpdateBlockSynced{} }
	pool[packet.IDUpdatePlayerGameType] = func() packet.Packet { return &UpdatePlayerGameType{} }
	pool[packet.IDClientBoundDebugRenderer] = func() packet.Packet { return &ClientBoundDebugRenderer{} }

	return pool
}
