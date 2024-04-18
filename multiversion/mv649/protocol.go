package mv649

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"

	"github.com/oomph-ac/mv/multiversion/mv649/packet"
	"github.com/oomph-ac/mv/multiversion/util"
	gtpacket "github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type Protocol struct{}

func (Protocol) ID() int32 {
	return 649
}

func (Protocol) Ver() string {
	return "1.20.60"
}

func (Protocol) NewReader(r minecraft.ByteReader, shieldID int32, enableLimits bool) protocol.IO {
	return protocol.NewReader(r, shieldID, enableLimits)
}

func (Protocol) NewWriter(r minecraft.ByteWriter, shieldID int32) protocol.IO {
	return protocol.NewWriter(r, shieldID)
}

func (Protocol) Packets(listener bool) gtpacket.Pool {
	if listener {
		return packet.NewClientPool()
	}
	return packet.NewServerPool()
}

func (Protocol) Encryption(key [32]byte) gtpacket.Encryption {
	return gtpacket.NewCTREncryption(key[:])
}

func (Protocol) ConvertToLatest(pk gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	if upgraded, ok := util.DefaultUpgrade(conn, pk, Mapping); ok {
		if upgraded == nil {
			return []gtpacket.Packet{}
		}

		return Upgrade([]gtpacket.Packet{upgraded}, conn)
	}

	return Upgrade([]gtpacket.Packet{pk}, conn)
}

func (Protocol) ConvertFromLatest(pk gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	if downgraded, ok := util.DefaultDowngrade(conn, pk, Mapping); ok {
		return Downgrade([]gtpacket.Packet{downgraded}, conn)
	}

	return Downgrade([]gtpacket.Packet{pk}, conn)
}

func Upgrade(pks []gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	packets := []gtpacket.Packet{}
	for _, pk := range pks {
		switch pk := pk.(type) {
		case *packet.PlayerAuthInput:
			packets = append(packets, &gtpacket.PlayerAuthInput{
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
		case *gtpacket.LecternUpdate:
			packets = append(packets, &packet.LecternUpdate{
				Page:      pk.Page,
				PageCount: pk.PageCount,
				Position:  pk.Position,
				DropBook:  false,
			})
		default:
			packets = append(packets, pk)
		}
	}

	return packets
}

func Downgrade(pks []gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	packets := []gtpacket.Packet{}
	for _, pk := range pks {
		switch pk := pk.(type) {
		case *gtpacket.AvailableCommands:
			for _, c := range pk.Commands {
				for _, o := range c.Overloads {
					for _, p := range o.Parameters {
						var newT uint32 = 0
						// check if p.type has the arg type, and add it to newT
						if p.Type&protocol.CommandArgValid != 0 {
							newT |= protocol.CommandArgValid
						}
						if p.Type&protocol.CommandArgEnum != 0 {
							newT |= protocol.CommandArgEnum
						}
						if p.Type&protocol.CommandArgSuffixed != 0 {
							newT |= protocol.CommandArgSuffixed
						}
						if p.Type&protocol.CommandArgSoftEnum != 0 {
							newT |= protocol.CommandArgSoftEnum
						}

						// Downgrade arg types.
						if p.Type&protocol.CommandArgTypeEquipmentSlots != 0 {
							newT |= packet.CommandArgTypeEquipmentSlots
						}
						if p.Type&protocol.CommandArgTypeString != 0 {
							newT |= packet.CommandArgTypeString
						}
						if p.Type&protocol.CommandArgTypeBlockPosition != 0 {
							newT |= packet.CommandArgTypeBlockPosition
						}
						if p.Type&protocol.CommandArgTypePosition != 0 {
							newT |= packet.CommandArgTypePosition
						}
						if p.Type&protocol.CommandArgTypeMessage != 0 {
							newT |= packet.CommandArgTypeMessage
						}
						if p.Type&protocol.CommandArgTypeRawText != 0 {
							newT |= packet.CommandArgTypeRawText
						}
						if p.Type&protocol.CommandArgTypeJSON != 0 {
							newT |= packet.CommandArgTypeJSON
						}
						if p.Type&protocol.CommandArgTypeBlockStates != 0 {
							newT |= packet.CommandArgTypeBlockStates
						}
						if p.Type&protocol.CommandArgTypeCommand != 0 {
							newT |= packet.CommandArgTypeCommand
						}

						p.Type = newT
					}
				}
			}

			packets = append(packets, pk)
		case *gtpacket.SetActorMotion:
			packets = append(packets, &packet.SetActorMotion{
				Velocity:        pk.Velocity,
				EntityRuntimeID: pk.EntityRuntimeID,
			})
		case *gtpacket.ResourcePacksInfo:
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
