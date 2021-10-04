package middleware

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

type PermissionsMiddleware struct {
	count int
}

var (
	_ ken.MiddlewareBefore = (*PermissionsMiddleware)(nil)
)

func (c *PermissionsMiddleware) Before(ctx *ken.Ctx) (next bool, err error) {
	guildRoles, err := ctx.Session.GuildRoles(ctx.Event.GuildID)
	if err != nil {
		return
	}

roleLoop:
	for _, rid := range ctx.Event.Member.Roles {
		for _, r := range guildRoles {
			if rid == r.ID && r.Name == "Admin" {
				next = true
				break roleLoop
			}
		}
	}

	if !next {
		err = ctx.Respond(&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Color:       0xF44336,
						Description: "You must have the role \"Admim\" to perform this command!",
					},
				},
			},
		})
	}
	return
}
