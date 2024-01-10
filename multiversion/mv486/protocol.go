package mv486

import (
	"fmt"

	"github.com/oomph-ac/mv/multiversion/mv486/packet"
	"github.com/oomph-ac/mv/multiversion/mv486/types"
	"github.com/oomph-ac/mv/multiversion/mv589"
	mv594_packet "github.com/oomph-ac/mv/multiversion/mv594/packet"
	"github.com/oomph-ac/mv/multiversion/util"
	"github.com/samber/lo"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	gtpacket "github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type Protocol struct{}

func (Protocol) ID() int32 {
	return 486
}

func (Protocol) Ver() string {
	return "1.18.12"
}

func (Protocol) NewReader(r minecraft.ByteReader, shieldID int32, enableLimits bool) protocol.IO {
	return protocol.NewReader(r, shieldID, enableLimits)
}

func (Protocol) NewWriter(w minecraft.ByteWriter, shieldID int32) protocol.IO {
	return protocol.NewWriter(w, shieldID)
}

func (Protocol) Packets(listener bool) gtpacket.Pool {
	if listener {
		return packet.NewServerPool()
	}
	return packet.NewClientPool()
}

func (Protocol) Encryption(key [32]byte) gtpacket.Encryption {
	return gtpacket.NewCTREncryption(key[:])
}

func (Protocol) ConvertToLatest(pk gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	if pk, ok := pk.(*gtpacket.PacketViolationWarning); ok {
		fmt.Printf("Violation: %d", pk.PacketID)
	}
	if pk, ok := pk.(*packet.RequestChunkRadius); ok {
		return []gtpacket.Packet{&gtpacket.RequestChunkRadius{
			ChunkRadius:    pk.ChunkRadius,
			MaxChunkRadius: pk.ChunkRadius,
		}}
	}

	if pk, ok := pk.(*packet.PlayerAuthInput); ok {
		return []gtpacket.Packet{&gtpacket.PlayerAuthInput{
			Pitch:         pk.Pitch,
			Yaw:           pk.Yaw,
			Position:      pk.Position,
			MoveVector:    pk.MoveVector,
			HeadYaw:       pk.HeadYaw,
			InputData:     pk.InputData,
			InputMode:     pk.InputMode,
			PlayMode:      pk.PlayMode,
			GazeDirection: pk.GazeDirection,
			Tick:          pk.Tick,
			Delta:         pk.Delta,
			ItemInteractionData: func(data protocol.UseItemTransactionData) protocol.UseItemTransactionData {
				data.LegacySetItemSlots = lo.Map(data.LegacySetItemSlots, func(item protocol.LegacySetItemSlot, _ int) protocol.LegacySetItemSlot {
					if item.ContainerID > 21 { // RECIPE_BOOK
						item.ContainerID -= 1
					}
					return item
				})
				return data
			}(pk.ItemInteractionData),
			ItemStackRequest: protocol.ItemStackRequest{
				RequestID:     pk.ItemStackRequest.RequestID,
				Actions:       pk.ItemStackRequest.Actions,
				FilterStrings: pk.ItemStackRequest.FilterStrings,
				FilterCause:   pk.ItemStackRequest.FilterCause,
			},
			BlockActions: pk.BlockActions,
		}}
	}

	if upgraded, ok := util.DefaultUpgrade(conn, pk, Mapping); ok {
		if upgraded == nil {
			return []gtpacket.Packet{}
		}

		return []gtpacket.Packet{upgraded}
	}

	return []gtpacket.Packet{pk}
}

func (Protocol) ConvertFromLatest(pk gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	if downgraded, ok := util.DefaultDowngrade(conn, pk, Mapping); ok {
		return []gtpacket.Packet{downgraded}
	}

	return Downgrade([]gtpacket.Packet{pk}, conn)
}

func Downgrade(pks []gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	packets := []gtpacket.Packet{}
	for _, pk := range mv589.Downgrade(pks, conn) {
		switch pk := pk.(type) {
		case *mv594_packet.StartGame:
			packets = append(packets, &packet.StartGame{
				EntityUniqueID:                 pk.EntityUniqueID,
				EntityRuntimeID:                pk.EntityRuntimeID,
				PlayerGameMode:                 pk.PlayerGameMode,
				PlayerPosition:                 pk.PlayerPosition,
				Pitch:                          pk.Pitch,
				Yaw:                            pk.Yaw,
				WorldSeed:                      int32(pk.WorldSeed),
				SpawnBiomeType:                 pk.SpawnBiomeType,
				UserDefinedBiomeName:           pk.UserDefinedBiomeName,
				Dimension:                      pk.Dimension,
				Generator:                      pk.Generator,
				WorldGameMode:                  pk.WorldGameMode,
				Difficulty:                     pk.Difficulty,
				WorldSpawn:                     pk.WorldSpawn,
				AchievementsDisabled:           pk.AchievementsDisabled,
				DayCycleLockTime:               pk.DayCycleLockTime,
				EducationEditionOffer:          pk.EducationEditionOffer,
				EducationFeaturesEnabled:       pk.EducationFeaturesEnabled,
				EducationProductID:             pk.EducationProductID,
				RainLevel:                      pk.RainLevel,
				LightningLevel:                 pk.LightningLevel,
				ConfirmedPlatformLockedContent: pk.ConfirmedPlatformLockedContent,
				MultiPlayerGame:                pk.MultiPlayerGame,
				LANBroadcastEnabled:            pk.LANBroadcastEnabled,
				XBLBroadcastMode:               pk.XBLBroadcastMode,
				PlatformBroadcastMode:          pk.PlatformBroadcastMode,
				CommandsEnabled:                pk.CommandsEnabled,
				TexturePackRequired:            pk.TexturePackRequired,
				GameRules:                      pk.GameRules,
				Experiments:                    pk.Experiments,
				ExperimentsPreviouslyToggled:   pk.ExperimentsPreviouslyToggled,
				BonusChestEnabled:              pk.BonusChestEnabled,
				StartWithMapEnabled:            pk.StartWithMapEnabled,
				PlayerPermissions:              pk.PlayerPermissions,
				ServerChunkTickRadius:          pk.ServerChunkTickRadius,
				HasLockedBehaviourPack:         pk.HasLockedBehaviourPack,
				HasLockedTexturePack:           pk.HasLockedTexturePack,
				FromLockedWorldTemplate:        pk.FromLockedWorldTemplate,
				MSAGamerTagsOnly:               pk.MSAGamerTagsOnly,
				FromWorldTemplate:              pk.FromWorldTemplate,
				WorldTemplateSettingsLocked:    pk.WorldTemplateSettingsLocked,
				OnlySpawnV1Villagers:           pk.OnlySpawnV1Villagers,
				BaseGameVersion:                pk.BaseGameVersion,
				LimitedWorldWidth:              pk.LimitedWorldWidth,
				LimitedWorldDepth:              pk.LimitedWorldDepth,
				NewNether:                      pk.NewNether,
				EducationSharedResourceURI:     pk.EducationSharedResourceURI,
				ForceExperimentalGameplay:      false,
				LevelID:                        pk.LevelID,
				WorldName:                      pk.WorldName,
				TemplateContentIdentity:        pk.TemplateContentIdentity,
				Trial:                          pk.Trial,
				PlayerMovementSettings:         pk.PlayerMovementSettings,
				Time:                           pk.Time,
				EnchantmentSeed:                pk.EnchantmentSeed,
				Blocks:                         pk.Blocks,
				Items:                          pk.Items,
				MultiPlayerCorrelationID:       pk.MultiPlayerCorrelationID,
				ServerAuthoritativeInventory:   pk.ServerAuthoritativeInventory,
				GameVersion:                    pk.GameVersion,
				ServerBlockStateChecksum:       pk.ServerBlockStateChecksum,
			})
		case *gtpacket.NetworkChunkPublisherUpdate:
			packets = append(packets, &packet.NetworkChunkPublisherUpdate{
				Position: pk.Position,
				Radius:   pk.Radius,
			})
		case *gtpacket.SetActorData:
			packets = append(packets, &gtpacket.SetActorData{
				EntityRuntimeID:  pk.EntityRuntimeID,
				EntityMetadata:   pk.EntityMetadata,
				EntityProperties: pk.EntityProperties,
				Tick:             pk.Tick,
			})
		case *gtpacket.UpdateAttributes:
			packets = append(packets, &packet.UpdateAttributes{
				EntityRuntimeID: pk.EntityRuntimeID,
				Attributes: lo.Map(pk.Attributes, func(item protocol.Attribute, _ int) types.Attribute {
					return types.Attribute{Attribute: item}
				}),
				Tick: pk.Tick,
			})
		default:
			packets = append(packets, pk)
		}
	}

	return packets
}
