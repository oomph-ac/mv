package mv662

import (
	"github.com/oomph-ac/mv/multiversion/mv662/packet"
	"github.com/oomph-ac/mv/multiversion/mv671"
	"github.com/oomph-ac/mv/multiversion/util"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	gtpacket "github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// Protocol is a protocol implementation for
// Minecraft: Bedrock Edition 1.20.70.
type Protocol struct{}

// ID ...
func (Protocol) ID() int32 {
	return 662
}

// Ver ...
func (Protocol) Ver() string {
	return "1.20.70"
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
	var packets []gtpacket.Packet
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
				InteractionModel:       uint32(pk.InteractionModel),
				GazeDirection:          pk.GazeDirection,
				Tick:                   pk.Tick,
				Delta:                  pk.Delta,
				ItemInteractionData:    pk.ItemInteractionData,
				ItemStackRequest:       pk.ItemStackRequest,
				BlockActions:           pk.BlockActions,
				VehicleRotation:        pk.VehicleRotation,
				ClientPredictedVehicle: pk.ClientPredictedVehicle,
				AnalogueMoveVector:     pk.AnalogueMoveVector,
			})
		default:
			packets = append(packets, pk)
		}
	}

	return mv671.Upgrade(packets, conn)
}

// Downgrade ...
func Downgrade(pks []gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	downgraded := mv671.Downgrade(pks, conn)
	packets := make([]gtpacket.Packet, 0, len(downgraded))

	for _, pk := range downgraded {
		switch pk := pk.(type) {
		case *gtpacket.ResourcePackStack:
			packets = append(packets, &packet.ResourcePackStack{
				TexturePackRequired:          pk.TexturePackRequired,
				BehaviourPacks:               pk.BehaviourPacks,
				TexturePacks:                 pk.TexturePacks,
				BaseGameVersion:              pk.BaseGameVersion,
				Experiments:                  pk.Experiments,
				ExperimentsPreviouslyToggled: pk.ExperimentsPreviouslyToggled,
			})
		case *gtpacket.StartGame:
			packets = append(packets, &packet.StartGame{
				EntityUniqueID:                 pk.EntityUniqueID,
				EntityRuntimeID:                pk.EntityRuntimeID,
				PlayerGameMode:                 pk.PlayerGameMode,
				PlayerPosition:                 pk.PlayerPosition,
				Pitch:                          pk.Pitch,
				Yaw:                            pk.Yaw,
				WorldSeed:                      pk.WorldSeed,
				SpawnBiomeType:                 pk.SpawnBiomeType,
				UserDefinedBiomeName:           pk.UserDefinedBiomeName,
				Dimension:                      pk.Dimension,
				Generator:                      pk.Generator,
				WorldGameMode:                  pk.WorldGameMode,
				Difficulty:                     pk.Difficulty,
				WorldSpawn:                     pk.WorldSpawn,
				AchievementsDisabled:           pk.AchievementsDisabled,
				EditorWorldType:                pk.EditorWorldType,
				CreatedInEditor:                pk.CreatedInEditor,
				ExportedFromEditor:             pk.ExportedFromEditor,
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
				PersonaDisabled:                pk.PersonaDisabled,
				CustomSkinsDisabled:            pk.CustomSkinsDisabled,
				EmoteChatMuted:                 pk.EmoteChatMuted,
				BaseGameVersion:                pk.BaseGameVersion,
				LimitedWorldWidth:              pk.LimitedWorldWidth,
				LimitedWorldDepth:              pk.LimitedWorldDepth,
				NewNether:                      pk.NewNether,
				EducationSharedResourceURI:     pk.EducationSharedResourceURI,
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
				PropertyData:                   pk.PropertyData,
				ServerBlockStateChecksum:       pk.ServerBlockStateChecksum,
				ClientSideGeneration:           pk.ClientSideGeneration,
				WorldTemplateID:                pk.WorldTemplateID,
				ChatRestrictionLevel:           pk.ChatRestrictionLevel,
				DisablePlayerInteractions:      pk.DisablePlayerInteractions,
				UseBlockNetworkIDHashes:        pk.UseBlockNetworkIDHashes,
				ServerAuthoritativeSound:       pk.ServerAuthoritativeSound,
			})
		case *gtpacket.UpdateBlockSynced:
			packets = append(packets, &packet.UpdateBlockSynced{
				Position:          pk.Position,
				NewBlockRuntimeID: pk.NewBlockRuntimeID,
				Flags:             pk.Flags,
				Layer:             pk.Layer,
				EntityUniqueID:    int64(pk.EntityUniqueID),
				TransitionType:    pk.TransitionType,
			})
		case *gtpacket.UpdatePlayerGameType:
			packets = append(packets, &packet.UpdatePlayerGameType{
				GameType:       pk.GameType,
				PlayerUniqueID: pk.PlayerUniqueID,
			})
		case *gtpacket.ClientBoundDebugRenderer:
			packets = append(packets, &packet.ClientBoundDebugRenderer{
				Type:     pk.Type,
				Text:     pk.Text,
				Position: pk.Position,
				Red:      pk.Red,
				Green:    pk.Green,
				Blue:     pk.Blue,
				Alpha:    pk.Alpha,
				Duration: pk.Duration,
			})
		case *gtpacket.CraftingData:
			recipes := make([]protocol.Recipe, 0, len(pk.Recipes))
			for _, r := range pk.Recipes {
				switch r := r.(type) {
				case *protocol.ShapedRecipe:
					recipes = append(recipes, &packet.ShapedRecipe{
						RecipeID:        r.RecipeID,
						Width:           r.Width,
						Height:          r.Height,
						Input:           r.Input,
						Output:          r.Output,
						UUID:            r.UUID,
						Block:           r.Block,
						Priority:        r.Priority,
						RecipeNetworkID: r.RecipeNetworkID,
					})
				case *protocol.ShapedChemistryRecipe:
					recipes = append(recipes, &packet.ShapedChemistryRecipe{
						ShapedRecipe: packet.ShapedRecipe{
							RecipeID:        r.RecipeID,
							Width:           r.Width,
							Height:          r.Height,
							Input:           r.Input,
							Output:          r.Output,
							UUID:            r.UUID,
							Block:           r.Block,
							Priority:        r.Priority,
							RecipeNetworkID: r.RecipeNetworkID,
						},
					})
				default:
					recipes = append(recipes, r)
				}
			}

			packets = append(packets, &packet.CraftingData{
				Recipes:                      recipes,
				PotionRecipes:                pk.PotionRecipes,
				PotionContainerChangeRecipes: pk.PotionContainerChangeRecipes,
				MaterialReducers:             pk.MaterialReducers,
				ClearRecipes:                 pk.ClearRecipes,
			})
		default:
			packets = append(packets, pk)
		}
	}

	return packets
}
