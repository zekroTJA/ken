package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

type InfoUserCommand struct{}

var (
	_ ken.UserCommand = (*InfoUserCommand)(nil)
)

func (c *InfoUserCommand) Name() string {
	return "info"
}

func (c *InfoUserCommand) Description() string {
	return "Dispaly user information."
}

func (c *InfoUserCommand) Run(ctx *ken.Ctx) (err error) {
	err = ctx.RespondEmbed(&discordgo.MessageEmbed{
		Description: ctx.User().String(),
	})
	return
}
