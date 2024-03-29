package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

type InfoUserCommand struct{}

var (
	_ ken.UserCommand = (*InfoUserCommand)(nil)
)

func (c *InfoUserCommand) TypeUser() {}

func (c *InfoUserCommand) Name() string {
	return "userinfo"
}

func (c *InfoUserCommand) Description() string {
	return "Display user information."
}

func (c *InfoUserCommand) Run(ctx ken.Context) (err error) {
	err = ctx.RespondEmbed(&discordgo.MessageEmbed{
		Description: ctx.User().String(),
	})
	return
}
