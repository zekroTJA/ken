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

func (c *DeleteMessageCommand) Run(ctx ken.Context) (err error) {
	var msg *discordgo.Message
	for _, msg = range ctx.GetEvent().ApplicationCommandData().Resolved.Messages {
		break
	}

	if err = ctx.GetSession().ChannelMessageDelete(msg.ChannelID, msg.ID); err != nil {
		return
	}

	err = ctx.RespondEmbed(&discordgo.MessageEmbed{
		Description: "Message deleted.",
	})
	return
}
