package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

type DeleteMessageCommand struct{}

var (
	_ ken.MessageCommand = (*DeleteMessageCommand)(nil)
)

func (c *DeleteMessageCommand) TypeMessage() {}

func (c *DeleteMessageCommand) Name() string {
	return "delete"
}

func (c *DeleteMessageCommand) Description() string {
	return "Delete the selected message"
}

func (c *DeleteMessageCommand) Run(ctx *ken.Ctx) (err error) {
	var msg *discordgo.Message
	for _, msg = range ctx.Event.ApplicationCommandData().Resolved.Messages {
		break
	}

	if err = ctx.Session.ChannelMessageDelete(msg.ChannelID, msg.ID); err != nil {
		return
	}

	err = ctx.RespondEmbed(&discordgo.MessageEmbed{
		Description: "Message deleted.",
	})
	return
}
