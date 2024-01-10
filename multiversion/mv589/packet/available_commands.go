package packet

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// AvailableCommands is sent by the server to send a list of all commands that
// the player is able to use on the server. This packet holds all the arguments
// of each commands as well, making it possible for the client to provide
// auto-completion and command usages. AvailableCommands packets can be resent,
// but the packet is often very big, so doing this very often should be avoided.
type AvailableCommands struct {
	// EnumValues is a slice of all enum values of any enum in the
	// AvailableCommands packet. EnumValues generally should contain each
	// possible value only once. Enums are built by pointing to entries in this
	// slice.
	EnumValues []string
	// Suffixes, like EnumValues, is a slice of all suffix values of any command
	// parameter in the AvailableCommands packet.
	Suffixes []string
	// Enums is a slice of all (fixed) command enums present in any of the
	// commands.
	Enums []protocol.CommandEnum
	// Commands is a list of all commands that the client should show
	// client-side. The AvailableCommands packet replaces any commands sent
	// before. It does not only add the commands that are sent in it.
	Commands []Command
	// DynamicEnums is a slice of dynamic command enums. These command enums can
	// be changed during runtime without having to resend an AvailableCommands
	// packet.
	DynamicEnums []protocol.DynamicEnum
	// Constraints is a list of constraints that should be applied to certain
	// options of enums in the commands above.
	Constraints []protocol.CommandEnumConstraint
}

// ID ...
func (*AvailableCommands) ID() uint32 {
	return packet.IDAvailableCommands
}

func (pk *AvailableCommands) Marshal(io protocol.IO) {
	protocol.FuncSlice(io, &pk.EnumValues, io.String)
	protocol.FuncSlice(io, &pk.Suffixes, io.String)
	protocol.FuncIOSlice(io, &pk.Enums, protocol.CommandEnumContext{EnumValues: pk.EnumValues}.Marshal)
	protocol.Slice(io, &pk.Commands)
	protocol.Slice(io, &pk.DynamicEnums)
	protocol.Slice(io, &pk.Constraints)
}

// Command holds the data that a command requires to be shown to a player client-side. The command is shown in
// the /help command and auto-completed using this data.
type Command struct {
	// Name is the name of the command. The command may be executed using this name, and will be shown in the
	// /help list with it. It currently seems that the client crashes if the Name contains uppercase letters.
	Name string
	// Description is the description of the command. It is shown in the /help list and when starting to write
	// a command.
	Description string
	// Flags is a combination of flags not currently known. Leaving the Flags field empty appears to work.
	Flags uint16
	// PermissionLevel is the command permission level that the player required to execute this command. The
	// field no longer seems to serve a purpose, as the client does not handle the execution of commands
	// anymore: The permissions should be checked server-side.
	PermissionLevel byte
	// AliasesOffset is the offset to a CommandEnum that holds the values that
	// should be used as aliases for this command.
	AliasesOffset uint32
	// Overloads is a list of command overloads that specify the ways in which a command may be executed. The
	// overloads may be completely different.
	Overloads []CommandOverload
}

func (c *Command) Marshal(r protocol.IO) {
	r.String(&c.Name)
	r.String(&c.Description)
	r.Uint16(&c.Flags)
	r.Uint8(&c.PermissionLevel)
	r.Uint32(&c.AliasesOffset)
	protocol.Slice(r, &c.Overloads)
}

func DowngradeCommands(c []protocol.Command) []Command {
	new := []Command{}
	for _, o := range c {
		new = append(new, Command{
			Name:            o.Name,
			Description:     o.Description,
			Flags:           o.Flags,
			PermissionLevel: o.PermissionLevel,
			AliasesOffset:   o.AliasesOffset,
			Overloads:       latestCommandOverloadsToSupported(o.Overloads),
		})
	}

	return new
}

// CommandOverload represents an overload of a command. This overload can be compared to function overloading
// in languages such as java. It represents a single usage of the command. A command may have multiple
// different overloads, which are handled differently.
type CommandOverload struct {
	// Parameters is a list of command parameters that are part of the overload. These parameters specify the
	// usage of the command when this overload is applied.
	Parameters []protocol.CommandParameter
}

func (c *CommandOverload) Marshal(r protocol.IO) {
	protocol.Slice(r, &c.Parameters)
}

func latestCommandOverloadsToSupported(c []protocol.CommandOverload) []CommandOverload {
	new := []CommandOverload{}
	for _, o := range c {
		new = append(new, CommandOverload{
			Parameters: o.Parameters,
		})
	}

	return new
}
