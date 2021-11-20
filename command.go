package ken

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Command interface {
	// Name returns the unique name of the command.
	Name() string

	// Description returns a brief text which concisely
	// describes the commands purpose.
	Description() string

	// Run is called on command invokation getting
	// passed the invocation context.
	//
	// When something goes wrong during command
	// execution, you can return an error which is
	// then handled by Ken's OnCommandError handler.
	Run(ctx *Ctx) (err error)
}

func toApplicationCommand(c Command) *discordgo.ApplicationCommand {
	switch cm := c.(type) {
	case SlashCommand:
		return &discordgo.ApplicationCommand{
			Name:        cm.Name(),
			Type:        discordgo.ChatApplicationCommand,
			Description: cm.Description(),
			Version:     cm.Version(),
			Options:     cm.Options(),
		}
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
	default:
		panic(fmt.Sprintf("Command type not implemented for command: %s", cm.Name()))
	}
}
