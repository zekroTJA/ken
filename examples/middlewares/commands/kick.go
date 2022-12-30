package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
	"github.com/zekrotja/ken/examples/middlewares/middlewares"
)

type KickCommand struct{}

var (
	_ ken.SlashCommand                = (*KickCommand)(nil)
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

func (c *KickCommand) Run(ctx ken.Context) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	user := ctx.Options().GetByName("member").UserValue(ctx)
	reason := ctx.Options().GetByName("reason").StringValue()

	if err = ctx.GetSession().GuildMemberDeleteWithReason(ctx.GetEvent().GuildID, user.ID, reason); err != nil {
		return
	}

	err = ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: fmt.Sprintf("Kicked member <@%s> with reason\n```\n%s```", user.ID, reason),
	}).Send().Error

	return
}
