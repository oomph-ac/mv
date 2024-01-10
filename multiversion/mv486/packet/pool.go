package packet

import (
	v589packet "github.com/oomph-ac/mv/multiversion/mv618/packet"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func NewClientPool() packet.Pool {
	pool := v589packet.NewClientPool()
	pool[packet.IDRequestChunkRadius] = func() packet.Packet { return &RequestChunkRadius{} }
	return pool
}

func NewServerPool() packet.Pool {
	pool := v589packet.NewServerPool()
	pool[packet.IDStartGame] = func() packet.Packet { return &StartGame{} }
	//pool[packet.IDResourcePacksInfo] = func() packet.Packet { return &ResourcePacksInfo{} }
	pool[packet.IDAddActor] = func() packet.Packet { return &AddActor{} }
	pool[packet.IDAddPlayer] = func() packet.Packet { return &AddPlayer{} }
	pool[packet.IDAddVolumeEntity] = func() packet.Packet { return &AddVolumeEntity{} }
	pool[packet.IDRequestChunkRadius] = func() packet.Packet { return &RequestChunkRadius{} }
	pool[packet.IDCommandRequest] = func() packet.Packet { return &CommandRequest{} }
	pool[packet.IDNetworkChunkPublisherUpdate] = func() packet.Packet { return &NetworkChunkPublisherUpdate{} }
	pool[packet.IDPlayerAction] = func() packet.Packet { return &PlayerAction{} }
	pool[packet.IDPlayerAuthInput] = func() packet.Packet { return &PlayerAuthInput{} }
	pool[packet.IDPlayerList] = func() packet.Packet { return &PlayerList{} }
	pool[packet.IDPlayerSkin] = func() packet.Packet { return &PlayerSkin{} }
	pool[packet.IDRemoveVolumeEntity] = func() packet.Packet { return &RemoveVolumeEntity{} }
	pool[packet.IDSpawnParticleEffect] = func() packet.Packet { return &SpawnParticleEffect{} }
	pool[packet.IDStartGame] = func() packet.Packet { return &StartGame{} }
	pool[packet.IDStructureBlockUpdate] = func() packet.Packet { return &StructureBlockUpdate{} }
	pool[packet.IDStructureTemplateDataRequest] = func() packet.Packet { return &StructureTemplateDataRequest{} }
	pool[packet.IDUpdateAttributes] = func() packet.Packet { return &UpdateAttributes{} }
	pool[packet.IDItemStackRequest] = func() packet.Packet { return &ItemStackRequest{} }
	pool[packet.IDModalFormResponse] = func() packet.Packet { return &ModalFormResponse{} }
	return pool
}
