package ratelimit

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

// Middleware command implements the ratelimit middleware.
type Middleware struct {
	manager Manager
	force   bool
}

var _ ken.Middleware = (*Middleware)(nil)

// New returns a new instance of Middleware.
func New(cfg ...Config) *Middleware {
	c := defaultConfig

	if len(cfg) != 0 {
		_c := cfg[0]

		if _c.Manager != nil {
			c.Manager = _c.Manager
		}
		if _c.Force {
			c.Force = true
		}
	}

	if c.Manager == nil {
		c.Manager = newInternalManager()
	}

	return &Middleware{
		manager: c.Manager,
		force:   c.Force,
	}
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
			guildID = ctx.GetEvent().GuildID
		}
	}

	limiter := m.manager.GetLimiter(ctx.Command, ctx.User().ID, guildID)
	if ok, next := limiter.Take(); !ok {
		err := ctx.RespondError(fmt.Sprintf(
			"You are being ratelimited.\nWait %s until you can use this command again.",
			next.Round(1*time.Second).String()), "Rate Limited")
		return false, err
	}

	if !m.force {
		ctx.Set(limiterKey, limiter)
	}

	return true, nil
}

func (m *Middleware) After(ctx *ken.Ctx, cmdError error) (err error) {
	if !m.force {
		sk, ok := ctx.Get(skipKey).(bool)
		if cmdError != nil || sk && ok {
			if l, ok := ctx.Get(limiterKey).(*Limiter); ok {
				l.Restore()
			}
		}
	}
	return
}
