package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

type TestCommand struct{}

var (
	_ ken.SlashCommand = (*TestCommand)(nil)
	_ ken.DmCapable    = (*TestCommand)(nil)
	_ CmdWithExtras    = (*TestCommand)(nil)
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
			Type:        discordgo.ApplicationCommandOptionBoolean,
			Name:        "pog",
			Required:    true,
			Description: "pog",
		},
	}
}

func (c *TestCommand) IsDmCapable() bool {
	return true
}

func (c *TestCommand) ExtraString() string {
	return "extra string"
}

func (c *TestCommand) ExtraInt() int {
	return 2
}

func (c *TestCommand) Run(ctx *ken.Ctx) (err error) {
	val := ctx.Options().GetByName("pog").BoolValue()

	msg := "not poggers"
	if val {
		msg = "poggers"
	}

	err = ctx.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	return
}
