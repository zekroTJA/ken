package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

type Guild1Command struct{}

var (
	_ ken.SlashCommand       = (*Guild1Command)(nil)
	_ ken.GuildScopedCommand = (*Guild1Command)(nil)
)

func (c *Guild1Command) Name() string {
	return "guild1"
}

func (c *Guild1Command) Guild() string {
	return "362162947738566657"
}

func (c *Guild1Command) Description() string {
	return "Basic Test Command - guild1"
}

func (c *Guild1Command) Version() string {
	return "1.0.0"
}

func (c *Guild1Command) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Guild1Command) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}

func (c *Guild1Command) Run(ctx *ken.Ctx) (err error) {
	err = ctx.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "This command is only available on guild 1.",
		},
	})
	return
}
