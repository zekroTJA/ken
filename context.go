package ken

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken/state"
)

// Ctx holds the invokation context of
// a command.
type Ctx struct {
	ObjectMap

	st state.State
	dp ObjectProvider

	// Session holds the discordgo session instance.
	Session *discordgo.Session
	// Event provides the InteractionCreate event
	// instance.
	Event *discordgo.InteractionCreate
	// Command provides the called command instance.
	Command Command
}

func newCtx() *Ctx {
	return &Ctx{
		ObjectMap: make(simpleObjectMap),
	}
}

func (c *Ctx) Respond(r *discordgo.InteractionResponse) error {
	return c.Session.InteractionRespond(c.Event.Interaction, r)
}

func (c *Ctx) FollowUp(wait bool, data *discordgo.WebhookParams) (fum *FollowUpMessage) {
	fum = &FollowUpMessage{
		s: c.Session,
		i: c.Event.Interaction,
	}
	fum.self, fum.Error = c.st.SelfUser(c.Session)
	if fum.Error != nil {
		return
	}
	fum.Message, fum.Error = c.Session.FollowupMessageCreate(fum.self.ID, c.Event.Interaction, wait, data)
	return
}

func (c *Ctx) FollowUpError(content, title string) (fum *FollowUpMessage) {
	return c.FollowUp(true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			{
				Description: content,
				Title:       title,
				Color:       clrEmbedError,
			},
		},
	})
}

func (c *Ctx) Channel() (*discordgo.Channel, error) {
	return c.st.Channel(c.Session, c.Event.ChannelID)
}

func (c *Ctx) Get(key string) (v interface{}) {
	if v = c.ObjectMap.Load(key); v == nil && c.dp != nil {
		v = c.dp.Load(key)
	}
	return
}
