package ken

import "github.com/bwmarrin/discordgo"

type Command interface {
	Name() string
	Description() string
	Version() string
	Type() discordgo.ApplicationCommandType
	Options() []*discordgo.ApplicationCommandOption

	Run(ctx *Ctx) (err error)
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
