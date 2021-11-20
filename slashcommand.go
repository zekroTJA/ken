package ken

import "github.com/bwmarrin/discordgo"

// SlashCommand defines a callable slash command.
type SlashCommand interface {
	Command

	// Version returns the commands semantic version.
	Version() string

	// Options returns an array of application
	// command options.
	Options() []*discordgo.ApplicationCommandOption
}

// DmCapable extends a command to specify if it is
// able to be executed in DMs or not.
type DmCapable interface {
	// IsDmCapable returns true if the command can
	// be used in DMs.
	IsDmCapable() bool
}
