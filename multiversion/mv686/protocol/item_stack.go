package protocol

import (
	"fmt"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// ItemStackRequest represents a single request present in an ItemStackRequest packet sent by the client to
// change an item in an inventory.
// Item stack requests are either approved or rejected by the server using the ItemStackResponse packet.
type ItemStackRequest struct {
	// RequestID is a unique ID for the request. This ID is used by the server to send a response for this
	// specific request in the ItemStackResponse packet.
	RequestID int32
	// Actions is a list of actions performed by the client. The actual type of the actions depends on which
	// ID was present, and is one of the concrete types below.
	Actions []protocol.StackRequestAction
	// FilterStrings is a list of filter strings involved in the request. This is typically filled with one string
	// when an anvil or cartography is used.
	FilterStrings []string
	// FilterCause represents the cause of any potential filtering. This is one of the constants above.
	FilterCause int32
}

// Marshal encodes/decodes an ItemStackRequest.
func (x *ItemStackRequest) Marshal(r protocol.IO) {
	r.Varint32(&x.RequestID)
	protocol.FuncSlice(r, &x.Actions, func(p *protocol.StackRequestAction) {
		if _, ok := r.(*protocol.Reader); ok {
			var id uint8
			r.Uint8(&id)
			if !lookupStackRequestAction(id, p) {
				r.UnknownEnumOption(id, "stack request action type")
				return
			}
			(*p).Marshal(r)
		} else {
			var id byte
			if !lookupStackRequestActionType(*p, &id) {
				r.UnknownEnumOption(fmt.Sprintf("%T", *p), "stack request action type")
			}
			r.Uint8(&id)
			(*p).Marshal(r)
		}
	})
	protocol.FuncSlice(r, &x.FilterStrings, r.String)
	r.Int32(&x.FilterCause)
}

// lookupStackRequestActionType looks up the ID of a StackRequestAction.
func lookupStackRequestActionType(x protocol.StackRequestAction, id *uint8) bool {
	switch x.(type) {
	case *TakeStackRequestAction:
		*id = protocol.StackRequestActionTake
	case *PlaceStackRequestAction:
		*id = protocol.StackRequestActionPlace
	case *SwapStackRequestAction:
		*id = protocol.StackRequestActionSwap
	case *DropStackRequestAction:
		*id = protocol.StackRequestActionDrop
	case *DestroyStackRequestAction:
		*id = protocol.StackRequestActionDestroy
	case *ConsumeStackRequestAction:
		*id = protocol.StackRequestActionConsume
	case *CraftRecipeStackRequestAction:
		*id = protocol.StackRequestActionCraftRecipe
	case *AutoCraftRecipeStackRequestAction:
		*id = protocol.StackRequestActionCraftRecipeAuto
	case *CraftCreativeStackRequestAction:
		*id = protocol.StackRequestActionCraftCreative
	case *CraftRecipeOptionalStackRequestAction:
		*id = protocol.StackRequestActionCraftRecipeOptional
	case *CraftGrindstoneRecipeStackRequestAction:
		*id = protocol.StackRequestActionCraftGrindstone
	case *protocol.CreateStackRequestAction:
		*id = protocol.StackRequestActionCreate
	case *protocol.LabTableCombineStackRequestAction:
		*id = protocol.StackRequestActionLabTableCombine
	case *protocol.BeaconPaymentStackRequestAction:
		*id = protocol.StackRequestActionBeaconPayment
	case *protocol.MineBlockStackRequestAction:
		*id = protocol.StackRequestActionMineBlock
	case *protocol.CraftLoomRecipeStackRequestAction:
		*id = protocol.StackRequestActionCraftLoom
	case *protocol.CraftNonImplementedStackRequestAction:
		*id = protocol.StackRequestActionCraftNonImplementedDeprecated
	case *protocol.CraftResultsDeprecatedStackRequestAction:
		*id = protocol.StackRequestActionCraftResultsDeprecated
	default:
		return false
	}
	return true
}

// lookupStackRequestAction looks up the StackRequestAction matching an ID.
func lookupStackRequestAction(id uint8, x *protocol.StackRequestAction) bool {
	switch id {
	case protocol.StackRequestActionTake:
		*x = &TakeStackRequestAction{}
	case protocol.StackRequestActionPlace:
		*x = &PlaceStackRequestAction{}
	case protocol.StackRequestActionSwap:
		*x = &SwapStackRequestAction{}
	case protocol.StackRequestActionDrop:
		*x = &DropStackRequestAction{}
	case protocol.StackRequestActionDestroy:
		*x = &DestroyStackRequestAction{}
	case protocol.StackRequestActionConsume:
		*x = &ConsumeStackRequestAction{}
	case protocol.StackRequestActionCraftRecipe:
		*x = &CraftRecipeStackRequestAction{}
	case protocol.StackRequestActionCraftRecipeAuto:
		*x = &AutoCraftRecipeStackRequestAction{}
	case protocol.StackRequestActionCraftCreative:
		*x = &CraftCreativeStackRequestAction{}
	case protocol.StackRequestActionCraftRecipeOptional:
		*x = &CraftRecipeOptionalStackRequestAction{}
	case protocol.StackRequestActionCraftGrindstone:
		*x = &CraftGrindstoneRecipeStackRequestAction{}
	case protocol.StackRequestActionCreate:
		*x = &protocol.CreateStackRequestAction{}
	case protocol.StackRequestActionPlaceInContainer:
		*x = &protocol.PlaceInContainerStackRequestAction{}
	case protocol.StackRequestActionTakeOutContainer:
		*x = &protocol.TakeOutContainerStackRequestAction{}
	case protocol.StackRequestActionLabTableCombine:
		*x = &protocol.LabTableCombineStackRequestAction{}
	case protocol.StackRequestActionBeaconPayment:
		*x = &protocol.BeaconPaymentStackRequestAction{}
	case protocol.StackRequestActionMineBlock:
		*x = &protocol.MineBlockStackRequestAction{}
	case protocol.StackRequestActionCraftLoom:
		*x = &protocol.CraftLoomRecipeStackRequestAction{}
	case protocol.StackRequestActionCraftNonImplementedDeprecated:
		*x = &protocol.CraftNonImplementedStackRequestAction{}
	case protocol.StackRequestActionCraftResultsDeprecated:
		*x = &protocol.CraftResultsDeprecatedStackRequestAction{}
	default:
		return false
	}
	return true
}

// CraftRecipeStackRequestAction is sent by the client the moment it begins crafting an item. This is the
// first action sent, before the Consume and Create item stack request actions.
// This action is also sent when an item is enchanted. Enchanting should be treated mostly the same way as
// crafting, where the old item is consumed.
type CraftRecipeStackRequestAction struct {
	// RecipeNetworkID is the network ID of the recipe that is about to be crafted. This network ID matches
	// one of the recipes sent in the CraftingData packet, where each of the recipes have a RecipeNetworkID as
	// of 1.16.
	RecipeNetworkID uint32
}

// Marshal ...
func (a *CraftRecipeStackRequestAction) Marshal(r protocol.IO) {
	r.Varuint32(&a.RecipeNetworkID)
}

// AutoCraftRecipeStackRequestAction is sent by the client similarly to the CraftRecipeStackRequestAction. The
// only difference is that the recipe is automatically created and crafted by shift clicking the recipe book.
type AutoCraftRecipeStackRequestAction struct {
	// RecipeNetworkID is the network ID of the recipe that is about to be crafted. This network ID matches
	// one of the recipes sent in the CraftingData packet, where each of the recipes have a RecipeNetworkID as
	// of 1.16.
	RecipeNetworkID uint32
	// TimesCrafted is how many times the recipe was crafted.
	TimesCrafted byte
	// Ingredients is a slice of ItemDescriptorCount that contains the ingredients that were used to craft the recipe.
	// It is not exactly clear what this is used for, but it is sent by the vanilla client.
	Ingredients []protocol.ItemDescriptorCount
}

// Marshal ...
func (a *AutoCraftRecipeStackRequestAction) Marshal(r protocol.IO) {
	r.Varuint32(&a.RecipeNetworkID)
	r.Uint8(&a.TimesCrafted)
	protocol.FuncSlice(r, &a.Ingredients, r.ItemDescriptorCount)
}

// CraftCreativeStackRequestAction is sent by the client when it takes an item out fo the creative inventory.
// The item is thus not really crafted, but instantly created.
type CraftCreativeStackRequestAction struct {
	// CreativeItemNetworkID is the network ID of the creative item that is being created. This is one of the
	// creative item network IDs sent in the CreativeContent packet.
	CreativeItemNetworkID uint32
}

// Marshal ...
func (a *CraftCreativeStackRequestAction) Marshal(r protocol.IO) {
	r.Varuint32(&a.CreativeItemNetworkID)
}

// CraftRecipeOptionalStackRequestAction is sent when using an anvil. When this action is sent, the
// FilterStrings field in the respective stack request is non-empty and contains the name of the item created
// using the anvil or cartography table.
type CraftRecipeOptionalStackRequestAction struct {
	// RecipeNetworkID is the network ID of the multi-recipe that is about to be crafted. This network ID matches
	// one of the multi-recipes sent in the CraftingData packet, where each of the recipes have a RecipeNetworkID as
	// of 1.16.
	RecipeNetworkID uint32
	// FilterStringIndex is the index of a filter string sent in a ItemStackRequest.
	FilterStringIndex int32
}

// Marshal ...
func (c *CraftRecipeOptionalStackRequestAction) Marshal(r protocol.IO) {
	r.Varuint32(&c.RecipeNetworkID)
	r.Int32(&c.FilterStringIndex)
}

// CraftGrindstoneRecipeStackRequestAction is sent when a grindstone recipe is crafted. It contains the RecipeNetworkID
// to identify the recipe crafted, and the cost for crafting the recipe.
type CraftGrindstoneRecipeStackRequestAction struct {
	// RecipeNetworkID is the network ID of the recipe that is about to be crafted. This network ID matches
	// one of the recipes sent in the CraftingData packet, where each of the recipes have a RecipeNetworkID as
	// of 1.16.
	RecipeNetworkID uint32
	// Cost is the cost of the recipe that was crafted.
	Cost int32
}

// Marshal ...
func (c *CraftGrindstoneRecipeStackRequestAction) Marshal(r protocol.IO) {
	r.Varuint32(&c.RecipeNetworkID)
	r.Varint32(&c.Cost)
}

// ItemStackResponse is a response to an individual ItemStackRequest.
type ItemStackResponse struct {
	// Status specifies if the request with the RequestID below was successful. If this is the case, the
	// ContainerInfo below will have information on what slots ended up changing. If not, the container info
	// will be empty.
	// A non-0 status means an error occurred and will result in the action being reverted.
	Status uint8
	// RequestID is the unique ID of the request that this response is in reaction to. If rejected, the client
	// will undo the actions from the request with this ID.
	RequestID int32
	// ContainerInfo holds information on the containers that had their contents changed as a result of the
	// request.
	ContainerInfo []StackResponseContainerInfo
}

// Marshal encodes/decodes an ItemStackResponse.
func (x *ItemStackResponse) Marshal(r protocol.IO) {
	r.Uint8(&x.Status)
	r.Varint32(&x.RequestID)
	if x.Status == protocol.ItemStackResponseStatusOK {
		protocol.Slice(r, &x.ContainerInfo)
	}
}

// StackResponseContainerInfo holds information on what slots in a container have what item stack in them.
type StackResponseContainerInfo struct {
	// ContainerID is the container ID of the container that the slots that follow are in. For the main
	// inventory, this value seems to be 0x1b. For the cursor, this value seems to be 0x3a. For the crafting
	// grid, this value seems to be 0x0d.
	ContainerID byte
	// SlotInfo holds information on what item stack should be present in specific slots in the container.
	SlotInfo []StackResponseSlotInfo
}

// Marshal encodes/decodes a StackResponseContainerInfo.
func (x *StackResponseContainerInfo) Marshal(r protocol.IO) {
	r.Uint8(&x.ContainerID)
	protocol.Slice(r, &x.SlotInfo)
}

// StackResponseSlotInfo holds information on what item stack should be present in a specific slot.
type StackResponseSlotInfo struct {
	// Slot and HotbarSlot seem to be the same value every time: The slot that was actually changed. I'm not
	// sure if these slots ever differ.
	Slot, HotbarSlot byte
	// Count is the total count of the item stack. This count will be shown client-side after the response is
	// sent to the client.
	Count byte
	// StackNetworkID is the network ID of the new stack at a specific slot.
	StackNetworkID int32
	// CustomName is the custom name of the item stack. It is used in relation to text filtering.
	CustomName string
	// DurabilityCorrection is the current durability of the item stack. This durability will be shown
	// client-side after the response is sent to the client.
	DurabilityCorrection int32
}

// Marshal encodes/decodes a StackResponseSlotInfo.
func (x *StackResponseSlotInfo) Marshal(r protocol.IO) {
	r.Uint8(&x.Slot)
	r.Uint8(&x.HotbarSlot)
	r.Uint8(&x.Count)
	r.Varint32(&x.StackNetworkID)
	if x.Slot != x.HotbarSlot {
		r.InvalidValue(x.HotbarSlot, "hotbar slot", "hot bar slot must be equal to normal slot")
	}
	r.String(&x.CustomName)
	r.Varint32(&x.DurabilityCorrection)
}

// transferStackRequestAction is the structure shared by StackRequestActions that transfer items from one
// slot into another.
type transferStackRequestAction struct {
	// Count is the count of the item in the source slot that was taken towards the destination slot.
	Count byte
	// Source and Destination point to the source slot from which Count of the item stack were taken and the
	// destination slot to which this item was moved.
	Source, Destination StackRequestSlotInfo
}

// Marshal ...
func (a *transferStackRequestAction) Marshal(r protocol.IO) {
	r.Uint8(&a.Count)
	StackReqSlotInfo(r, &a.Source)
	StackReqSlotInfo(r, &a.Destination)
}

// TakeStackRequestAction is sent by the client to the server to take x amount of items from one slot in a
// container to the cursor.
type TakeStackRequestAction struct {
	transferStackRequestAction
}

// PlaceStackRequestAction is sent by the client to the server to place x amount of items from one slot into
// another slot, such as when shift clicking an item in the inventory to move it around or when moving an item
// in the cursor into a slot.
type PlaceStackRequestAction struct {
	transferStackRequestAction
}

// SwapStackRequestAction is sent by the client to swap the item in its cursor with an item present in another
// container. The two item stacks swap places.
type SwapStackRequestAction struct {
	// Source and Destination point to the source slot from which Count of the item stack were taken and the
	// destination slot to which this item was moved.
	Source, Destination StackRequestSlotInfo
}

// Marshal ...
func (a *SwapStackRequestAction) Marshal(r protocol.IO) {
	StackReqSlotInfo(r, &a.Source)
	StackReqSlotInfo(r, &a.Destination)
}

// DropStackRequestAction is sent by the client when it drops an item out of the inventory when it has its
// inventory opened. This action is not sent when a player drops an item out of the hotbar using the Q button
// (or the equivalent on mobile). The InventoryTransaction packet is still used for that action, regardless of
// whether the item stack network IDs are used or not.
type DropStackRequestAction struct {
	// Count is the count of the item in the source slot that was taken towards the destination slot.
	Count byte
	// Source is the source slot from which items were dropped to the ground.
	Source StackRequestSlotInfo
	// Randomly seems to be set to false in most cases. I'm not entirely sure what this does, but this is what
	// vanilla calls this field.
	Randomly bool
}

// Marshal ...
func (a *DropStackRequestAction) Marshal(r protocol.IO) {
	r.Uint8(&a.Count)
	StackReqSlotInfo(r, &a.Source)
	r.Bool(&a.Randomly)
}

// DestroyStackRequestAction is sent by the client when it destroys an item in creative mode by moving it
// back into the creative inventory.
type DestroyStackRequestAction struct {
	// Count is the count of the item in the source slot that was destroyed.
	Count byte
	// Source is the source slot from which items came that were destroyed by moving them into the creative
	// inventory.
	Source StackRequestSlotInfo
}

// Marshal ...
func (a *DestroyStackRequestAction) Marshal(r protocol.IO) {
	r.Uint8(&a.Count)
	StackReqSlotInfo(r, &a.Source)
}

// ConsumeStackRequestAction is sent by the client when it uses an item to craft another item. The original
// item is 'consumed'.
type ConsumeStackRequestAction struct {
	DestroyStackRequestAction
}

// StackRequestSlotInfo holds information on a specific slot client-side.
type StackRequestSlotInfo struct {
	// ContainerID is the ID of the container that the slot was in.
	ContainerID byte
	// Slot is the index of the slot within the container with the ContainerID above.
	Slot byte
	// StackNetworkID is the unique stack ID that the client assumes to be present in this slot. The server
	// must check if these IDs match. If they do not match, servers should reject the stack request that the
	// action holding this info was in.
	StackNetworkID int32
}

// StackReqSlotInfo reads/writes a StackRequestSlotInfo x using IO r.
func StackReqSlotInfo(r protocol.IO, x *StackRequestSlotInfo) {
	r.Uint8(&x.ContainerID)
	r.Uint8(&x.Slot)
	r.Varint32(&x.StackNetworkID)
}
