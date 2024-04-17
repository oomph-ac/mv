package packet

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

func NewClientPool() packet.Pool {
	pool := packet.NewClientPool()
	pool[packet.IDPlayerAuthInput] = func() packet.Packet { return &PlayerAuthInput{} }
	pool[packet.IDLecternUpdate] = func() packet.Packet { return &LecternUpdate{} }

	return pool
}

func NewServerPool() packet.Pool {
	pool := packet.NewServerPool()
	pool[packet.IDMobEffect] = func() packet.Packet { return &MobEffect{} }
	pool[packet.IDResourcePacksInfo] = func() packet.Packet { return &ResourcePacksInfo{} }
	pool[packet.IDSetActorMotion] = func() packet.Packet { return &SetActorMotion{} }

	return pool
}
