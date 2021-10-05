package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
	"github.com/zekrotja/ken/example/middlewares/middlewares"
)

type KickCommand struct{}

var (
	_ ken.Command                     = (*KickCommand)(nil)
	_ middlewares.RequiresRoleCommand = (*KickCommand)(nil)
)

func (c *KickCommand) Name() string {
	return "kick"
}

func (c *KickCommand) Description() string {
	return "Kick a member from the guild."
}

func (c *KickCommand) Version() string {
	return "1.0.0"
}

func (c *KickCommand) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *KickCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "member",
			Description: "The member to be kicked.",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "reason",
			Description: "The kick reason.",
			Required:    true,
		},
	}
}

func (c *KickCommand) RequiresRole() string {
	return "Admin"
}

func (c *KickCommand) Run(ctx *ken.Ctx) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	user := ctx.Options().GetByName("member").UserValue(ctx)
	reason := ctx.Options().GetByName("reason").StringValue()

	if err = ctx.Session.GuildMemberDeleteWithReason(ctx.Event.GuildID, user.ID, reason); err != nil {
		return
	}

	err = ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: fmt.Sprintf("Kicked member <@%s> with reason\n```\n%s```", user.ID, reason),
	}).Error

	return
}
