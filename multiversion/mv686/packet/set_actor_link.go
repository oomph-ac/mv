package packet

import (
	v686protocol "github.com/oomph-ac/mv/multiversion/mv686/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	gtpacket "github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// SetActorLink is sent by the server to initiate an entity link client-side, meaning one entity will start
// riding another.
type SetActorLink struct {
	// EntityLink is the link to be set client-side. It links two entities together, so that one entity rides
	// another. Note that players that see those entities later will not see the link, unless it is also sent
	// in the AddActor and AddPlayer packets.
	EntityLink v686protocol.EntityLink
}

// ID ...
func (*SetActorLink) ID() uint32 {
	return gtpacket.IDSetActorLink
}

func (pk *SetActorLink) Marshal(io protocol.IO) {
	protocol.Single(io, &pk.EntityLink)
}
