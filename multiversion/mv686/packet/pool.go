package packet

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// NewClientPool returns a new pool containing packets sent by a client.
// Packets may be retrieved from it simply by indexing it with the packet ID.
func NewClientPool() packet.Pool {
	pool := packet.NewClientPool()

	pool[packet.IDInventoryTransaction] = func() packet.Packet { return &InventoryTransaction{} }
	pool[packet.IDDisconnect] = func() packet.Packet { return &Disconnect{} }
	pool[packet.IDPlayerAuthInput] = func() packet.Packet { return &PlayerAuthInput{} }
	pool[packet.IDItemStackRequest] = func() packet.Packet { return &ItemStackRequest{} }

	delete(pool, packet.IDServerBoundDiagnostics)
	delete(pool, packet.IDServerBoundLoadingScreen)

	return pool
}

// NewServerPool returns a new pool containing packets sent by a server.
// Packets may be retrieved from it simply by indexing it with the packet ID.
func NewServerPool() packet.Pool {
	pool := packet.NewServerPool()

	pool[packet.IDSetActorLink] = func() packet.Packet { return &SetActorLink{} }
	pool[packet.IDResourcePacksInfo] = func() packet.Packet { return &ResourcePacksInfo{} }
	pool[packet.IDMobArmourEquipment] = func() packet.Packet { return &MobArmourEquipment{} }
	pool[packet.IDSetTitle] = func() packet.Packet { return &SetTitle{} }
	pool[packet.IDStopSound] = func() packet.Packet { return &StopSound{} }
	pool[packet.IDInventorySlot] = func() packet.Packet { return &InventorySlot{} }
	pool[packet.IDDisconnect] = func() packet.Packet { return &Disconnect{} }
	pool[packet.IDCameraInstruction] = func() packet.Packet { return &CameraInstruction{} }
	pool[packet.IDInventoryContent] = func() packet.Packet { return &InventoryContent{} }
	pool[packet.IDChangeDimension] = func() packet.Packet { return &ChangeDimension{} }
	pool[packet.IDAddPlayer] = func() packet.Packet { return &AddPlayer{} }
	pool[packet.IDAddActor] = func() packet.Packet { return &AddActor{} }
	pool[packet.IDCorrectPlayerMovePrediction] = func() packet.Packet { return &CorrectPlayerMovePrediction{} }
	pool[packet.IDItemStackResponse] = func() packet.Packet { return &ItemStackResponse{} }

	delete(pool, packet.IDCurrentStructureFeature)
	delete(pool, packet.IDJigsawStructureData)

	return pool
}
