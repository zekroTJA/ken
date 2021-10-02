package ratelimit

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

// Middleware command implements the ratelimit middleware.
type Middleware struct {
	manager Manager
}

var _ ken.MiddlewareBefore = (*Middleware)(nil)

// New returns a new instance of Middleware.
//
// Optionally, you can pass a custom Manager instance
// if you want to handle limiters differently than the
// standard Manager implementation.
func New(manager ...Manager) *Middleware {
	var m Manager

	if len(manager) > 0 && manager[0] != nil {
		m = manager[0]
	} else {
		m = newInternalManager()
	}

	return &Middleware{m}
}

func (m *Middleware) Before(ctx *ken.Ctx) (next bool, err error) {
	c, ok := ctx.Command.(LimitedCommand)
	if !ok {
		return true, nil
	}

	var guildID string
	if c.IsLimiterGlobal() {
		guildID = "__global__"
	} else {
		var ch *discordgo.Channel
		ch, err = ctx.Channel()
		if err != nil {
			return
		}
		if ch.Type == discordgo.ChannelTypeDM || ch.Type == discordgo.ChannelTypeGroupDM {
			guildID = "__dm__"
		} else {
			guildID = ctx.Event.GuildID
		}
	}

	limiter := m.manager.GetLimiter(ctx.Command, ctx.Event.User.ID, guildID)
	if ok, next := limiter.Take(); !ok {
		err := ctx.FollowUpError(fmt.Sprintf(
			"You are being ratelimited.\nWait %s until you can use this command again.",
			next.String()), "Rate Limited").Error
		return false, err
	}

	return true, nil
}
