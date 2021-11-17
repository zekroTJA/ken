package commands

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
	"github.com/zekrotja/ken/middlewares/ratelimit"
)

type TestCommand struct{}

var (
	_ ken.Command              = (*TestCommand)(nil)
	_ ken.DmCapable            = (*TestCommand)(nil)
	_ ratelimit.LimitedCommand = (*TestCommand)(nil)
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

func (c *TestCommand) LimiterBurst() int {
	return 2
}

func (c *TestCommand) LimiterRestoration() time.Duration {
	return 30 * time.Second
}

func (c *TestCommand) IsLimiterGlobal() bool {
	return false
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
