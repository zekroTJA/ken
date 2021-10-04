package ken

import (
	"github.com/bwmarrin/discordgo"
)

// Ctx holds the invokation context of
// a command.
type Ctx struct {
	ObjectMap

	k *Ken

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

// Respond to an interaction event with the given
// interaction response payload.
func (c *Ctx) Respond(r *discordgo.InteractionResponse) error {
	return c.Session.InteractionRespond(c.Event.Interaction, r)
}

// FollowUp creates a follow up message to the
// interaction event and returns a FollowUpMessage
// object containing the created message as well as
// an error instance, if an error occurred.
//
// This way it allows to be chained in one call with
// subsequent FollowUpMessage method calls.
func (c *Ctx) FollowUp(wait bool, data *discordgo.WebhookParams) (fum *FollowUpMessage) {
	fum = &FollowUpMessage{
		s: c.Session,
		i: c.Event.Interaction,
	}
	fum.self, fum.Error = c.k.opt.State.SelfUser(c.Session)
	if fum.Error != nil {
		return
	}
	fum.Message, fum.Error = c.Session.FollowupMessageCreate(fum.self.ID, c.Event.Interaction, wait, data)
	return
}

// FollowUpEmbed is shorthand for FollowUp with an
// embed payload as passed.
func (c *Ctx) FollowUpEmbed(emb *discordgo.MessageEmbed) (fum *FollowUpMessage) {
	if emb.Color <= 0 {
		emb.Color = c.k.opt.EmbedColors.Default
	}
	return c.FollowUp(true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			emb,
		},
	})
}

// FollowUpError is shorthand for FollowUpEmbed with an
// error embed as message with the passed content and
// title.
func (c *Ctx) FollowUpError(content, title string) (fum *FollowUpMessage) {
	return c.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: content,
		Title:       title,
		Color:       c.k.opt.EmbedColors.Error,
	})
}

// Defer is shorthand for Respond with an InteractionResponse
// of the type InteractionResponseDeferredChannelMessageWithSource.
//
// It should be used when the interaction response can not be
// instantly returned.
func (c *Ctx) Defer() (err error) {
	return c.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
}

// Get either returns an instance from the internal object map -
// if existent. Otherwise, the object is looked up in the specified
// dependency provider, if available. When no object was found in
// either of both maps, nil is returned.
func (c *Ctx) Get(key string) (v interface{}) {
	if v = c.ObjectMap.Get(key); v == nil && c.k.opt.DependencyProvider != nil {
		v = c.k.opt.DependencyProvider.Get(key)
	}
	return
}

// Channel tries to fetch the channel object from the contained
// channel ID using the specified state manager.
func (c *Ctx) Channel() (*discordgo.Channel, error) {
	return c.k.opt.State.Channel(c.Session, c.Event.ChannelID)
}

// Channel tries to fetch the guild object from the contained
// guild ID using the specified state manager.
func (c *Ctx) Guild() (*discordgo.Guild, error) {
	return c.k.opt.State.Guild(c.Session, c.Event.GuildID)
}

// User returns the User object of the executor either from
// the events User object or from the events Member object.
func (c *Ctx) User() (u *discordgo.User) {
	u = c.Event.User
	if u == nil && c.Event.Member != nil {
		u = c.Event.Member.User
	}
	return
}
