package packet

import (
	v686protocol "github.com/oomph-ac/mv/multiversion/mv686/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	gtpacket "github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// ItemStackResponse is sent by the server in response to an ItemStackRequest packet from the client. This
// packet is used to either approve or reject ItemStackRequests from the client. If a request is approved, the
// client will simply continue as normal. If rejected, the client will undo the actions so that the inventory
// should be in sync with the server again.
type ItemStackResponse struct {
	// Responses is a list of responses to ItemStackRequests sent by the client before. Responses either
	// approve or reject a request from the client.
	// Vanilla limits the size of this slice to 4096.
	Responses []v686protocol.ItemStackResponse
}

// ID ...
func (*ItemStackResponse) ID() uint32 {
	return gtpacket.IDItemStackResponse
}

func (pk *ItemStackResponse) Marshal(io protocol.IO) {
	protocol.Slice(io, &pk.Responses)
}
