package packet

import (
	v630packet "github.com/oomph-ac/mv/multiversion/mv630/packet"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func NewClientPool() packet.Pool {
	pool := v630packet.NewClientPool()
	return pool
}

func NewServerPool() packet.Pool {
	pool := v630packet.NewServerPool()
	pool[packet.IDShowStoreOffer] = func() packet.Packet { return &ShowStoreOffer{} }
	return pool
}
