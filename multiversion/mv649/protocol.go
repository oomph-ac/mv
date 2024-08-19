package mv649

import (
	"github.com/go-gl/mathgl/mgl32"
	v686packet "github.com/oomph-ac/mv/multiversion/mv686/packet"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"

	"github.com/oomph-ac/mv/multiversion/mv649/packet"
	"github.com/oomph-ac/mv/multiversion/mv662"
	v662packet "github.com/oomph-ac/mv/multiversion/mv662/packet"
	"github.com/oomph-ac/mv/multiversion/util"
	gtpacket "github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// Protocol is a protocol implementation for
// Minecraft: Bedrock Edition 1.20.60.
type Protocol struct{}

// ID ...
func (Protocol) ID() int32 {
	return 649
}

// Ver ...
func (Protocol) Ver() string {
	return "1.20.60"
}

// NewReader ...
func (Protocol) NewReader(r minecraft.ByteReader, shieldID int32, enableLimits bool) protocol.IO {
	return protocol.NewReader(r, shieldID, enableLimits)
}

// NewWriter ...
func (Protocol) NewWriter(r minecraft.ByteWriter, shieldID int32) protocol.IO {
	return protocol.NewWriter(r, shieldID)
}

// Packets ...
func (Protocol) Packets(listener bool) gtpacket.Pool {
	if listener {
		return packet.NewClientPool()
	}
	return packet.NewServerPool()
}

// ConvertToLatest ...
func (Protocol) ConvertToLatest(pk gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	if upgraded, ok := util.DefaultUpgrade(conn, pk, Mapping); ok {
		if upgraded == nil {
			return []gtpacket.Packet{}
		}

		return Upgrade([]gtpacket.Packet{upgraded}, conn)
	}

	return Upgrade([]gtpacket.Packet{pk}, conn)
}

// ConvertFromLatest ...
func (Protocol) ConvertFromLatest(pk gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	if downgraded, ok := util.DefaultDowngrade(conn, pk, Mapping); ok {
		return Downgrade([]gtpacket.Packet{downgraded}, conn)
	}

	return Downgrade([]gtpacket.Packet{pk}, conn)
}

// Upgrade ...
func Upgrade(pks []gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	packets := make([]gtpacket.Packet, 0, len(pks))
	for _, pk := range pks {
		switch pk := pk.(type) {
		case *packet.PlayerAuthInput:
			packets = append(packets, &v662packet.PlayerAuthInput{
				Pitch:                  pk.Pitch,
				Yaw:                    pk.Yaw,
				Position:               pk.Position,
				MoveVector:             pk.MoveVector,
				HeadYaw:                pk.HeadYaw,
				InputData:              pk.InputData,
				InputMode:              pk.InputMode,
				PlayMode:               pk.PlayMode,
				InteractionModel:       pk.InteractionModel,
				GazeDirection:          pk.GazeDirection,
				Tick:                   pk.Tick,
				Delta:                  pk.Delta,
				ItemInteractionData:    pk.ItemInteractionData,
				ItemStackRequest:       pk.ItemStackRequest,
				BlockActions:           pk.BlockActions,
				ClientPredictedVehicle: pk.ClientPredictedVehicle,
				AnalogueMoveVector:     pk.AnalogueMoveVector,
				VehicleRotation:        mgl32.Vec2{},
			})
		case *packet.LecternUpdate:
			packets = append(packets, &gtpacket.LecternUpdate{
				Page:      pk.Page,
				PageCount: pk.PageCount,
				Position:  pk.Position,
			})
		default:
			packets = append(packets, pk)
		}
	}

	return mv662.Upgrade(packets, conn)
}

// Downgrade ...
func Downgrade(pks []gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	downgraded := mv662.Downgrade(pks, conn)
	packets := make([]gtpacket.Packet, 0, len(downgraded))

	for _, pk := range downgraded {
		switch pk := pk.(type) {
		case *gtpacket.AvailableCommands:
			for _, cmd := range pk.Commands {
				for _, overload := range cmd.Overloads {
					for _, param := range overload.Parameters {
						var newT uint32 = protocol.CommandArgValid

						switch param.Type | protocol.CommandArgValid {
						case protocol.CommandArgTypeEquipmentSlots:
							newT |= packet.CommandArgTypeEquipmentSlots
						case protocol.CommandArgTypeString:
							newT |= packet.CommandArgTypeString
						case protocol.CommandArgTypeBlockPosition:
							newT |= packet.CommandArgTypeBlockPosition
						case protocol.CommandArgTypePosition:
							newT |= packet.CommandArgTypePosition
						case protocol.CommandArgTypeMessage:
							newT |= packet.CommandArgTypeMessage
						case protocol.CommandArgTypeRawText:
							newT |= packet.CommandArgTypeRawText
						case protocol.CommandArgTypeJSON:
							newT |= packet.CommandArgTypeJSON
						case protocol.CommandArgTypeBlockStates:
							newT |= packet.CommandArgTypeBlockStates
						case protocol.CommandArgTypeCommand:
							newT |= packet.CommandArgTypeCommand
						}

						param.Type = newT
					}
				}
			}
			packets = append(packets, pk)
		case *gtpacket.SetActorMotion:
			packets = append(packets, &packet.SetActorMotion{
				Velocity:        pk.Velocity,
				EntityRuntimeID: pk.EntityRuntimeID,
			})
		case *v686packet.ResourcePacksInfo:
			packets = append(packets, &packet.ResourcePacksInfo{
				TexturePackRequired: pk.TexturePackRequired,
				HasScripts:          pk.HasScripts,
				BehaviourPacks:      pk.BehaviourPacks,
				TexturePacks:        pk.TexturePacks,
				ForcingServerPacks:  pk.ForcingServerPacks,
				PackURLs:            pk.PackURLs,
			})
		case *gtpacket.MobEffect:
			packets = append(packets, &packet.MobEffect{
				EntityRuntimeID: pk.EntityRuntimeID,
				Operation:       pk.Operation,
				EffectType:      pk.EffectType,
				Amplifier:       pk.Amplifier,
				Particles:       pk.Particles,
				Duration:        pk.Duration,
			})
		default:
			packets = append(packets, pk)
		}
	}

	return packets
}
