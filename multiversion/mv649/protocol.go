package mv649

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"

	"github.com/oomph-ac/mv/multiversion/mv649/packet"
	"github.com/oomph-ac/mv/multiversion/mv662"
	v662packet "github.com/oomph-ac/mv/multiversion/mv662/packet"
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

func Downgrade(pks []gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	downgraded := mv662.Downgrade(pks, conn)
	packets := make([]gtpacket.Packet, 0, len(downgraded))

	for _, pk := range downgraded {
		switch pk := pk.(type) {
		case *gtpacket.AvailableCommands:
			// HACK!!! Why??!?!?! because GOLANG doesn't like it when i just replace p.Type :////
			cmds := make([]protocol.Command, 0, len(pk.Commands))
			for _, c := range pk.Commands {
				cmd := protocol.Command{}
				cmd.Name = c.Name
				cmd.Description = c.Description
				cmd.Flags = c.Flags
				cmd.PermissionLevel = c.PermissionLevel
				cmd.AliasesOffset = c.AliasesOffset
				cmd.ChainedSubcommandOffsets = c.ChainedSubcommandOffsets
				cmd.Overloads = make([]protocol.CommandOverload, 0, len(c.Overloads))

				for _, o := range c.Overloads {
					overload := protocol.CommandOverload{}
					overload.Chaining = o.Chaining
					overload.Parameters = make([]protocol.CommandParameter, 0, len(o.Parameters))

					for _, p := range o.Parameters {
						param := protocol.CommandParameter{}
						param.Name = p.Name
						param.Optional = p.Optional
						param.Options = p.Options

						var newT uint32 = protocol.CommandArgValid
						if p.Type == (protocol.CommandArgTypeEquipmentSlots | protocol.CommandArgValid) {
							newT |= packet.CommandArgTypeEquipmentSlots
						} else if p.Type == (protocol.CommandArgTypeString | protocol.CommandArgValid) {
							newT |= packet.CommandArgTypeString
						} else if p.Type == (protocol.CommandArgTypeBlockPosition | protocol.CommandArgValid) {
							newT |= packet.CommandArgTypeBlockPosition
						} else if p.Type == (protocol.CommandArgTypePosition | protocol.CommandArgValid) {
							newT |= packet.CommandArgTypePosition
						} else if p.Type == (protocol.CommandArgTypeMessage | protocol.CommandArgValid) {
							newT |= packet.CommandArgTypeMessage
						} else if p.Type == (protocol.CommandArgTypeRawText | protocol.CommandArgValid) {
							newT |= packet.CommandArgTypeRawText
						} else if p.Type == (protocol.CommandArgTypeJSON | protocol.CommandArgValid) {
							newT |= packet.CommandArgTypeJSON
						} else if p.Type == (protocol.CommandArgTypeBlockStates | protocol.CommandArgValid) {
							newT |= packet.CommandArgTypeBlockStates
						} else if p.Type == (protocol.CommandArgTypeCommand | protocol.CommandArgValid) {
							newT |= packet.CommandArgTypeCommand
						} else {
							// We don't need to downgrade these.
							continue
						}

						param.Type = newT
						overload.Parameters = append(overload.Parameters, param)
					}

					cmd.Overloads = append(cmd.Overloads, overload)
				}

				cmds = append(cmds, cmd)
			}
			pk.Commands = cmds
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
