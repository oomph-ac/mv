package packet

import (
	gtpacket "github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// NOTE: CorrectPlayerMovementPrediction is not included in here, since changes
// to the packet were made late, and it was updated around 1.20.50 (630).

const (
	IDResourcePackStack        uint32 = 7
	IDStartGame                uint32 = 11
	IDCraftingData             uint32 = 52
	IDUpdateBlockSynced        uint32 = 110
	IDPlayerAuthInput          uint32 = 144
	IDUpdatePlayerGameType     uint32 = 151
	IDClientBoundDebugRenderer uint32 = 164
)

func NewClientPool() gtpacket.Pool {
	pool := gtpacket.NewClientPool()
	pool[IDPlayerAuthInput] = func() gtpacket.Packet { return &PlayerAuthInput{} }

	return pool
}

func NewServerPool() gtpacket.Pool {
	pool := gtpacket.NewServerPool()
	pool[IDResourcePackStack] = func() gtpacket.Packet { return &ResourcePackStack{} }
	pool[IDStartGame] = func() gtpacket.Packet { return &StartGame{} }
	pool[IDCraftingData] = func() gtpacket.Packet { return &CraftingData{} }
	pool[IDUpdateBlockSynced] = func() gtpacket.Packet { return &UpdateBlockSynced{} }
	pool[IDUpdatePlayerGameType] = func() gtpacket.Packet { return &UpdatePlayerGameType{} }
	pool[IDClientBoundDebugRenderer] = func() gtpacket.Packet { return &ClientBoundDebugRenderer{} }

	return pool
}
