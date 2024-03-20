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

func (Protocol) Name() string {
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
		return Upgrade(upgraded, conn)
	}

	return Upgrade(pk, conn)
}

func (Protocol) ConvertFromLatest(pk gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	if downgraded, ok := util.DefaultDowngrade(conn, pk, Mapping); ok {
		return Downgrade(downgraded, conn)
	}

	return Downgrade(pk, conn)
}

func Upgrade(pk gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	packets := []gtpacket.Packet{}
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
	}

	return packets
}

func Downgrade(pk gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	packets := []gtpacket.Packet{}

	switch pk := pk.(type) {
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
	}

	return packets
}
