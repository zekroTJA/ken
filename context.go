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

func (c *Ctx) FollowUpEmbed(emb *discordgo.MessageEmbed) (fum *FollowUpMessage) {
	return c.FollowUp(true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			emb,
		},
	})
}

func (c *Ctx) FollowUpError(content, title string) (fum *FollowUpMessage) {
	return c.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: content,
		Title:       title,
		Color:       clrEmbedError,
	})
}

func (c *Ctx) Defer() (err error) {
	return c.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
}

func (c *Ctx) Get(key string) (v interface{}) {
	if v = c.ObjectMap.Get(key); v == nil && c.dp != nil {
		v = c.dp.Get(key)
	}
	return
}

func (c *Ctx) Channel() (*discordgo.Channel, error) {
	return c.st.Channel(c.Session, c.Event.ChannelID)
}

func (c *Ctx) Guild() (*discordgo.Guild, error) {
	return c.st.Guild(c.Session, c.Event.GuildID)
}

func (c *Ctx) User() (u *discordgo.User) {
	u = c.Event.User
	if u == nil && c.Event.Member != nil {
		u = c.Event.Member.User
	}
	return
}
