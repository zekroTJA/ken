package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

type SubsCommand struct{}

var (
	_ ken.SlashCommand = (*SubsCommand)(nil)
	_ ken.DmCapable    = (*SubsCommand)(nil)
)

func (c *SubsCommand) Name() string {
	return "subs"
}

func (c *SubsCommand) Description() string {
	return "An example command with sub commands."
}

func (c *SubsCommand) Version() string {
	return "1.0.0"
}

func (c *SubsCommand) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *SubsCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "one",
			Description: "First sub command",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "arg",
					Description: "Argument",
					Required:    true,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "two",
			Description: "Second sub command",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "arg",
					Description: "Argument",
					Required:    false,
				},
			},
		},
	}
}

func (c *SubsCommand) IsDmCapable() bool {
	return true
}

func (c *SubsCommand) Run(ctx ken.Context) (err error) {
	err = ctx.HandleSubCommands(
		ken.SubCommandHandler{"one", c.one},
		ken.SubCommandHandler{"two", c.two},
	)

	return
}

func (c *SubsCommand) one(ctx *ken.SubCommandCtx) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}
	arg := ctx.Options().GetByName("arg").StringValue()
	err = ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: "one: " + arg,
	}).Error
	return
}

func (c *SubsCommand) two(ctx *ken.SubCommandCtx) (err error) {
	var arg int
	if argV, ok := ctx.Options().GetByNameOptional("arg"); ok {
		arg = int(argV.IntValue())
	}
	err = ctx.RespondEmbed(&discordgo.MessageEmbed{
		Description: fmt.Sprintf("two: %d", arg),
	})
	return
}
