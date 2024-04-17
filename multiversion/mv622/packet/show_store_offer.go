package packet

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	gtpacket "github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type ShowStoreOffer struct {
	// OfferID is a string that identifies the offer for which a window should be opened. While typically a
	// UUID, the ID could be anything.
	OfferID string
	// ShowAll specifies if all other offers of the same 'author' as the one of the offer associated with the
	// OfferID should also be displayed, alongside the target offer.
	ShowAll bool
}

func (*ShowStoreOffer) ID() uint32 {
	return gtpacket.IDShowStoreOffer
}

func (pk *ShowStoreOffer) Marshal(io protocol.IO) {
	io.String(&pk.OfferID)
	io.Bool(&pk.ShowAll)
}
