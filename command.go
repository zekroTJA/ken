package ken

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Command specifies the base interface for an
// application command.
type Command interface {
	// Name returns the unique name of the command.
	Name() string

	// Description returns a brief text which concisely
	// describes the commands purpose.
	//
	// Currently, this is ignored by user and message
	// commands, because the API currently does not
	// support descriptions for these types of
	// application commands.
	Description() string

	// Run is called on command invokation getting
	// passed the invocation context.
	//
	// When something goes wrong during command
	// execution, you can return an error which is
	// then handled by Ken's OnCommandError handler.
	Run(ctx Context) (err error)
}

// GuildScopedCommand can be implemented by your
// commands to scope them to specific guilds.
//
// The command then will be only registered on
// the guild returned by the Guild method.
type GuildScopedCommand interface {
	Guild() string
}

func toApplicationCommand(c Command) *discordgo.ApplicationCommand {
	switch cm := c.(type) {
	case UserCommand:
		return &discordgo.ApplicationCommand{
			Name: cm.Name(),
			Type: discordgo.UserApplicationCommand,
		}
	case MessageCommand:
		return &discordgo.ApplicationCommand{
			Name: cm.Name(),
			Type: discordgo.MessageApplicationCommand,
		}
	case SlashCommand:
		return &discordgo.ApplicationCommand{
			Name:        cm.Name(),
			Type:        discordgo.ChatApplicationCommand,
			Description: cm.Description(),
			Version:     cm.Version(),
			Options:     cm.Options(),
		}
	default:
		panic(fmt.Sprintf("Command type not implemented for command: %s", cm.Name()))
	}
}
