package ken

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken/state"
)

type Ctx struct {
	st state.State

	Session *discordgo.Session
	Event   *discordgo.InteractionCreate
	Command Command
}

func newCtx() *Ctx {
	return &Ctx{}
}

func (c *Ctx) Respond(r *discordgo.InteractionResponse) error {
	return c.Session.InteractionRespond(c.Event.Interaction, r)
}

func (c *Ctx) FollowUp(wait bool, data *discordgo.WebhookParams) (err error, fum *FollowUpMessage) {
	self, err := c.st.SelfUser(c.Session)
	if err != nil {
		return
	}
	msg, err := c.Session.FollowupMessageCreate(self.ID, c.Event.Interaction, wait, data)
	if err != nil {
		return
	}
	fum = &FollowUpMessage{
		Message: msg,
		self:    self.ID,
		s:       c.Session,
		i:       c.Event.Interaction,
	}
	return
}
