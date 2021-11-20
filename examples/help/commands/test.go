package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
	"github.com/zekrotja/ken/middlewares/cmdhelp"
)

type TestCommand struct{}

var (
	_ ken.SlashCommand     = (*TestCommand)(nil)
	_ ken.DmCapable        = (*TestCommand)(nil)
	_ cmdhelp.HelpProvider = (*TestCommand)(nil)
)

func (c *TestCommand) Name() string {
	return "test"
}

func (c *TestCommand) Description() string {
	return "Basic Test Command"
}

func (c *TestCommand) Version() string {
	return "1.0.0"
}

func (c *TestCommand) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *TestCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "help",
			Description: "Show help",
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "pog",
			Description: "Pog",
		},
	}
}

func (c *TestCommand) IsDmCapable() bool {
	return true
}

func (C *TestCommand) Help(ctx *ken.SubCommandCtx) (emb *discordgo.MessageEmbed, err error) {
	emb = &discordgo.MessageEmbed{
		Color:       0x00ff00,
		Description: "This is a help message that describes how to use this command!",
	}
	return
}

func (c *TestCommand) Run(ctx *ken.Ctx) (err error) {
	err = ctx.HandleSubCommands(
		ken.SubCommandHandler{"pog", c.pog},
	)

	return
}

func (c *TestCommand) pog(ctx *ken.SubCommandCtx) (err error) {
	return ctx.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "https://i1.sndcdn.com/avatars-ypCd5dE5YbGkyF0p-Y59d9w-t500x500.jpg",
		},
	})
}
