package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

type Guild2Command struct{}

var (
	_ ken.SlashCommand       = (*Guild1Command)(nil)
	_ ken.GuildScopedCommand = (*Guild2Command)(nil)
)

func (c *Guild2Command) Name() string {
	return "guild2"
}

func (c *Guild2Command) Guild() string {
	return "526196711962705925"
}

func (c *Guild2Command) Description() string {
	return "Basic Test Command - guild2"
}

func (c *Guild2Command) Version() string {
	return "1.0.0"
}

func (c *Guild2Command) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Guild2Command) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}

func (c *Guild2Command) Run(ctx *ken.Ctx) (err error) {
	err = ctx.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "This command is only available on guild 2.",
		},
	})
	return
}
