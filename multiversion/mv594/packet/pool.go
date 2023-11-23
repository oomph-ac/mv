package packet

import (
	v618packet "github.com/oomph-ac/mv/multiversion/mv618/packet"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func NewClientPool() packet.Pool {
	pool := v618packet.NewClientPool()
	return pool
}

func NewServerPool() packet.Pool {
	pool := v618packet.NewServerPool()
	pool[packet.IDStartGame] = func() packet.Packet { return &StartGame{} }
	pool[packet.IDResourcePacksInfo] = func() packet.Packet { return &ResourcePacksInfo{} }

	return pool
}
