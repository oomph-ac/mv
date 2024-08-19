package mv686

import (
	"github.com/oomph-ac/mv/multiversion/mv686/packet"
	v686protocol "github.com/oomph-ac/mv/multiversion/mv686/protocol"
	"github.com/oomph-ac/mv/multiversion/util"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	gtpacket "github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// Protocol is a protocol implementation for
// Minecraft: Bedrock Edition 1.21.2.
type Protocol struct {
}

// ID ...
func (p Protocol) ID() int32 {
	return 686
}

// Ver ...
func (p Protocol) Ver() string {
	return "1.21.2"
}

// Packets ...
func (p Protocol) Packets(listener bool) gtpacket.Pool {
	if listener {
		return packet.NewClientPool()
	}
	return packet.NewServerPool()
}

// NewReader ...
func (Protocol) NewReader(r minecraft.ByteReader, shieldID int32, enableLimits bool) protocol.IO {
	return protocol.NewReader(r, shieldID, enableLimits)
}

// NewWriter ...
func (Protocol) NewWriter(r minecraft.ByteWriter, shieldID int32) protocol.IO {
	return protocol.NewWriter(r, shieldID)
}

// ConvertToLatest ...
func (Protocol) ConvertToLatest(pk gtpacket.Packet, conn *minecraft.Conn) []gtpacket.Packet {
	if upgraded, ok := util.DefaultUpgrade(conn, pk, Mapping); ok {
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
func Upgrade(pks []gtpacket.Packet, _ *minecraft.Conn) []gtpacket.Packet {
	packets := make([]gtpacket.Packet, 0, len(pks))
	for _, pk := range pks {
		switch pk := pk.(type) {
		case *packet.ItemStackRequest:
			requests := make([]protocol.ItemStackRequest, len(pk.Requests))
			for ri, request := range pk.Requests {
				actions := make([]protocol.StackRequestAction, len(request.Actions))
				for ai, action := range request.Actions {
					var newAction protocol.StackRequestAction
					switch action := action.(type) {
					case *v686protocol.TakeStackRequestAction:
						a := &protocol.TakeStackRequestAction{}
						a.Count = action.Count
						a.Source = upgradeStackRequestSlotInfo(action.Source)
						a.Destination = upgradeStackRequestSlotInfo(action.Destination)

						newAction = a
					case *v686protocol.PlaceStackRequestAction:
						a := &protocol.PlaceStackRequestAction{}
						a.Count = action.Count
						a.Source = upgradeStackRequestSlotInfo(action.Source)
						a.Destination = upgradeStackRequestSlotInfo(action.Destination)

						newAction = a
					case *v686protocol.SwapStackRequestAction:
						a := &protocol.SwapStackRequestAction{}
						a.Source = upgradeStackRequestSlotInfo(action.Source)
						a.Destination = upgradeStackRequestSlotInfo(action.Destination)

						newAction = a
					case *v686protocol.DropStackRequestAction:
						a := &protocol.DropStackRequestAction{}
						a.Count = action.Count
						a.Source = upgradeStackRequestSlotInfo(action.Source)
						a.Randomly = action.Randomly

						newAction = a
					case *v686protocol.DestroyStackRequestAction:
						a := &protocol.DestroyStackRequestAction{}
						a.Count = action.Count
						a.Source = upgradeStackRequestSlotInfo(action.Source)

						newAction = a
					case *v686protocol.CraftRecipeStackRequestAction:
						newAction = &protocol.CraftRecipeStackRequestAction{
							RecipeNetworkID: action.RecipeNetworkID,
							NumberOfCrafts:  1,
						}
					case *v686protocol.AutoCraftRecipeStackRequestAction:
						newAction = &protocol.AutoCraftRecipeStackRequestAction{
							RecipeNetworkID: action.RecipeNetworkID,
							NumberOfCrafts:  1,
							TimesCrafted:    action.TimesCrafted,
							Ingredients:     action.Ingredients,
						}
					case *v686protocol.CraftCreativeStackRequestAction:
						newAction = &protocol.CraftCreativeStackRequestAction{
							CreativeItemNetworkID: action.CreativeItemNetworkID,
							NumberOfCrafts:        1,
						}
					case *v686protocol.CraftRecipeOptionalStackRequestAction:
						newAction = &protocol.CraftRecipeOptionalStackRequestAction{
							RecipeNetworkID:   action.RecipeNetworkID,
							NumberOfCrafts:    1,
							FilterStringIndex: action.FilterStringIndex,
						}
					case *v686protocol.CraftGrindstoneRecipeStackRequestAction:
						newAction = &protocol.CraftGrindstoneRecipeStackRequestAction{
							RecipeNetworkID: action.RecipeNetworkID,
							NumberOfCrafts:  1,
							Cost:            action.Cost,
						}
					default:
						newAction = action
					}
					actions[ai] = newAction
				}
				requests[ri] = protocol.ItemStackRequest{
					RequestID:     request.RequestID,
					Actions:       actions,
					FilterStrings: request.FilterStrings,
					FilterCause:   request.FilterCause,
				}
			}
			packets = append(packets, &gtpacket.ItemStackRequest{
				Requests: requests,
			})
		case *packet.InventoryTransaction:
			data := pk.TransactionData
			if u, ok := data.(*v686protocol.UseItemTransactionData); ok {
				data = &protocol.UseItemTransactionData{
					LegacyRequestID:    u.LegacyRequestID,
					LegacySetItemSlots: u.LegacySetItemSlots,
					Actions:            u.Actions,
					ActionType:         u.ActionType,
					BlockPosition:      u.BlockPosition,
					BlockFace:          u.BlockFace,
					HotBarSlot:         u.HotBarSlot,
					HeldItem:           u.HeldItem,
					Position:           u.Position,
					ClickedPosition:    u.ClickedPosition,
					BlockRuntimeID:     u.BlockRuntimeID,
					TriggerType:        0,
					ClientPrediction:   0,
				}
			}

			packets = append(packets, &gtpacket.InventoryTransaction{
				LegacyRequestID:    pk.LegacyRequestID,
				LegacySetItemSlots: pk.LegacySetItemSlots,
				Actions:            pk.Actions,
				TransactionData:    data,
			})
		case *packet.Disconnect:
			packets = append(packets, &gtpacket.Disconnect{
				Reason:                  pk.Reason,
				HideDisconnectionScreen: pk.HideDisconnectionScreen,
				Message:                 pk.Message,
				FilteredMessage:         "",
			})
		case *packet.PlayerAuthInput:
			packets = append(packets, &gtpacket.PlayerAuthInput{
				Pitch:            pk.Pitch,
				Yaw:              pk.Yaw,
				Position:         pk.Position,
				MoveVector:       pk.MoveVector,
				HeadYaw:          pk.HeadYaw,
				InputData:        pk.InputData,
				InputMode:        pk.InputMode,
				PlayMode:         pk.PlayMode,
				InteractionModel: pk.InteractionModel,
				GazeDirection:    pk.GazeDirection,
				Tick:             pk.Tick,
				Delta:            pk.Delta,
				ItemInteractionData: protocol.UseItemTransactionData{
					LegacyRequestID:    pk.ItemInteractionData.LegacyRequestID,
					LegacySetItemSlots: pk.ItemInteractionData.LegacySetItemSlots,
					Actions:            pk.ItemInteractionData.Actions,
					ActionType:         pk.ItemInteractionData.ActionType,
					BlockPosition:      pk.ItemInteractionData.BlockPosition,
					BlockFace:          pk.ItemInteractionData.BlockFace,
					HotBarSlot:         pk.ItemInteractionData.HotBarSlot,
					HeldItem:           pk.ItemInteractionData.HeldItem,
					Position:           pk.ItemInteractionData.Position,
					ClickedPosition:    pk.ItemInteractionData.ClickedPosition,
					BlockRuntimeID:     pk.ItemInteractionData.BlockRuntimeID,
					TriggerType:        protocol.TriggerTypeUnknown,
					ClientPrediction:   protocol.ClientPredictionSuccess,
				},
				ItemStackRequest: protocol.ItemStackRequest{
					RequestID:     pk.ItemStackRequest.RequestID,
					Actions:       pk.ItemStackRequest.Actions,
					FilterStrings: pk.ItemStackRequest.FilterStrings,
					FilterCause:   pk.ItemStackRequest.FilterCause,
				},
				BlockActions:           pk.BlockActions,
				VehicleRotation:        pk.VehicleRotation,
				ClientPredictedVehicle: pk.ClientPredictedVehicle,
				AnalogueMoveVector:     pk.AnalogueMoveVector,
			})
		default:
			packets = append(packets, pk)
		}
	}
	return packets
}

// Downgrade ...
func Downgrade(pks []gtpacket.Packet, _ *minecraft.Conn) []gtpacket.Packet {
	packets := make([]gtpacket.Packet, 0, len(pks))
	for _, pk := range pks {
		switch pk := pk.(type) {
		case *gtpacket.ItemStackResponse:
			responses := make([]v686protocol.ItemStackResponse, len(pk.Responses))
			for i, response := range pk.Responses {
				newResponse := v686protocol.ItemStackResponse{
					Status:        response.Status,
					RequestID:     response.RequestID,
					ContainerInfo: make([]v686protocol.StackResponseContainerInfo, len(response.ContainerInfo)),
				}
				for containerIndex, containerInfo := range response.ContainerInfo {
					newInfo := v686protocol.StackResponseContainerInfo{
						ContainerID: containerInfo.Container.ContainerID,
						SlotInfo:    make([]v686protocol.StackResponseSlotInfo, len(containerInfo.SlotInfo)),
					}
					for slotIndex, slotInfo := range containerInfo.SlotInfo {
						newInfo.SlotInfo[slotIndex] = v686protocol.StackResponseSlotInfo{
							Slot:                 slotInfo.Slot,
							HotbarSlot:           slotInfo.HotbarSlot,
							Count:                slotInfo.Count,
							StackNetworkID:       slotInfo.StackNetworkID,
							CustomName:           slotInfo.CustomName,
							DurabilityCorrection: slotInfo.DurabilityCorrection,
						}
					}
					newResponse.ContainerInfo[containerIndex] = newInfo
				}
				responses[i] = newResponse
			}

			packets = append(packets, &packet.ItemStackResponse{
				Responses: responses,
			})
		case *gtpacket.CorrectPlayerMovePrediction:
			packets = append(packets, &packet.CorrectPlayerMovePrediction{
				PredictionType: pk.PredictionType,
				Position:       pk.Position,
				Delta:          pk.Delta,
				Rotation:       pk.Rotation,
				OnGround:       pk.OnGround,
				Tick:           pk.Tick,
			})
		case *gtpacket.ResourcePacksInfo:
			packets = append(packets, &packet.ResourcePacksInfo{
				TexturePackRequired: pk.TexturePackRequired,
				HasAddons:           pk.HasAddons,
				HasScripts:          pk.HasScripts,
				BehaviourPacks:      downgradeBehaviourPacks(pk.BehaviourPacks),
				TexturePacks:        downgradeTexturePacks(pk.TexturePacks),
				ForcingServerPacks:  pk.ForcingServerPacks,
				PackURLs:            pk.PackURLs,
			})
		case *gtpacket.SetActorLink:
			packets = append(packets, &packet.SetActorLink{
				EntityLink: downgradeEntityLink(pk.EntityLink),
			})
		case *gtpacket.AddActor:
			packets = append(packets, &packet.AddActor{
				EntityUniqueID:   pk.EntityUniqueID,
				EntityRuntimeID:  pk.EntityRuntimeID,
				EntityType:       pk.EntityType,
				Position:         pk.Position,
				Velocity:         pk.Velocity,
				Pitch:            pk.Pitch,
				Yaw:              pk.Yaw,
				HeadYaw:          pk.HeadYaw,
				BodyYaw:          pk.BodyYaw,
				Attributes:       pk.Attributes,
				EntityMetadata:   pk.EntityMetadata,
				EntityProperties: pk.EntityProperties,
				EntityLinks:      downgradeEntityLinks(pk.EntityLinks),
			})
		case *gtpacket.AddPlayer:
			packets = append(packets, &packet.AddPlayer{
				UUID:             pk.UUID,
				Username:         pk.Username,
				EntityRuntimeID:  pk.EntityRuntimeID,
				PlatformChatID:   pk.PlatformChatID,
				Position:         pk.Position,
				Velocity:         pk.Velocity,
				Pitch:            pk.Pitch,
				Yaw:              pk.Yaw,
				HeadYaw:          pk.HeadYaw,
				HeldItem:         pk.HeldItem,
				GameType:         pk.GameType,
				EntityMetadata:   pk.EntityMetadata,
				EntityProperties: pk.EntityProperties,
				AbilityData:      pk.AbilityData,
				EntityLinks:      downgradeEntityLinks(pk.EntityLinks),
				DeviceID:         pk.DeviceID,
				BuildPlatform:    pk.BuildPlatform,
			})
		case *gtpacket.MobArmourEquipment:
			packets = append(packets, &packet.MobArmourEquipment{
				EntityRuntimeID: pk.EntityRuntimeID,
				Helmet:          pk.Helmet,
				Chestplate:      pk.Chestplate,
				Leggings:        pk.Leggings,
				Boots:           pk.Boots,
			})
		case *gtpacket.SetTitle:
			packets = append(packets, &packet.SetTitle{
				ActionType:       pk.ActionType,
				Text:             pk.Text,
				FadeInDuration:   pk.FadeInDuration,
				RemainDuration:   pk.RemainDuration,
				FadeOutDuration:  pk.FadeOutDuration,
				XUID:             pk.XUID,
				PlatformOnlineID: pk.PlatformOnlineID,
			})
		case *gtpacket.StopSound:
			packets = append(packets, &packet.StopSound{
				SoundName: pk.SoundName,
				StopAll:   pk.StopAll,
			})
		case *gtpacket.InventorySlot:
			packets = append(packets, &packet.InventorySlot{
				WindowID: pk.WindowID,
				Slot:     pk.Slot,
				NewItem:  pk.NewItem,
			})
		case *gtpacket.Disconnect:
			packets = append(packets, &packet.Disconnect{
				Reason:                  pk.Reason,
				HideDisconnectionScreen: pk.HideDisconnectionScreen,
				Message:                 pk.Message,
			})
		case *gtpacket.CameraInstruction:
			packets = append(packets, &packet.CameraInstruction{
				Set:   pk.Set,
				Clear: pk.Clear,
				Fade:  pk.Fade,
			})
		case *gtpacket.InventoryContent:
			packets = append(packets, &packet.InventoryContent{
				WindowID: pk.WindowID,
				Content:  pk.Content,
			})
		case *gtpacket.ChangeDimension:
			packets = append(packets, &packet.ChangeDimension{
				Dimension: pk.Dimension,
				Position:  pk.Position,
				Respawn:   pk.Respawn,
			})
		default:
			packets = append(packets, pk)
		}
	}
	return packets
}

// downgradeBehaviourPacks transforms latest []protocol.TexturePackInfo to []v686protocol.TexturePackInfo.
func downgradeTexturePacks(infos []protocol.TexturePackInfo) []v686protocol.TexturePackInfo {
	newTexturePacks := make([]v686protocol.TexturePackInfo, len(infos))
	for i, info := range infos {
		newTexturePacks[i] = v686protocol.TexturePackInfo{
			UUID:            info.UUID,
			Version:         info.Version,
			Size:            info.Size,
			ContentKey:      info.ContentKey,
			SubPackName:     info.SubPackName,
			ContentIdentity: info.ContentIdentity,
			HasScripts:      info.HasScripts,
		}
	}
	return newTexturePacks
}

// downgradeBehaviourPacks transforms latest []protocol.BehaviourPackInfo to []v686protocol.BehaviourPackInfo.
func downgradeBehaviourPacks(infos []protocol.BehaviourPackInfo) []v686protocol.BehaviourPackInfo {
	newTexturePacks := make([]v686protocol.BehaviourPackInfo, len(infos))
	for i, info := range infos {
		newTexturePacks[i] = v686protocol.BehaviourPackInfo{
			UUID:            info.UUID,
			Version:         info.Version,
			Size:            info.Size,
			ContentKey:      info.ContentKey,
			SubPackName:     info.SubPackName,
			ContentIdentity: info.ContentIdentity,
			HasScripts:      info.HasScripts,
		}
	}
	return newTexturePacks
}

// upgradeStackRequestSlotInfo transforms v686protocol.StackRequestSlotInfo to the latest version.
func upgradeStackRequestSlotInfo(info v686protocol.StackRequestSlotInfo) protocol.StackRequestSlotInfo {
	return protocol.StackRequestSlotInfo{
		Container: protocol.FullContainerName{
			ContainerID: info.ContainerID,
		},
		Slot:           info.Slot,
		StackNetworkID: info.StackNetworkID,
	}
}

// downgradeEntityLinks transforms array of latest []protocol.EntityLink to []v686protocol.EntityLink.
func downgradeEntityLinks(links []protocol.EntityLink) []v686protocol.EntityLink {
	newLinks := make([]v686protocol.EntityLink, 0, len(links))
	for _, link := range links {
		newLinks = append(newLinks, downgradeEntityLink(link))
	}
	return newLinks
}

// downgradeEntityLink transforms latest protocol.EntityLink to v686protocol.EntityLink.
func downgradeEntityLink(link protocol.EntityLink) v686protocol.EntityLink {
	return v686protocol.EntityLink{
		RiddenEntityUniqueID: link.RiddenEntityUniqueID,
		RiderEntityUniqueID:  link.RiderEntityUniqueID,
		Type:                 link.Type,
		Immediate:            link.Immediate,
		RiderInitiated:       link.RiderInitiated,
	}
}
