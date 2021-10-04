package ken

import "github.com/bwmarrin/discordgo"

// Command defines a callable slash command.
type Command interface {
	// Name returns the unique name of the command.
	Name() string
	// Description returns a brief text which concisely
	// describes the commands purpose.
	Description() string
	// Version returns the commands semantic version.
	Version() string
	// Type returns the commands command type.
	Type() discordgo.ApplicationCommandType
	// Options returns an array of application
	// command options.
	Options() []*discordgo.ApplicationCommandOption

	// Run is called on command invokation getting
	// passed the invocation context.
	//
	// When something goes wrong during command
	// execution, you can return an error which is
	// then handled by Ken's OnCommandError handler.
	Run(ctx *Ctx) (err error)
}

// DmCapable extends a command to specify if it is
// able to be executed in DMs or not.
type DmCapable interface {
	// IsDmCapable returns true if the command can
	// be used in DMs.
	IsDmCapable() bool
}

func toApplicationCommand(c Command) *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        c.Name(),
		Type:        c.Type(),
		Description: c.Description(),
		Version:     c.Version(),
		Options:     c.Options(),
	}
}
