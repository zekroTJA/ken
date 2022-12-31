package ken

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// FollowUpMessageBuilder builds a followup
// message interaction response.
type FollowUpMessageBuilder struct {
	ken *Ken
	i   *discordgo.Interaction

	data *discordgo.WebhookParams
	wait bool

	componentBuilder *ComponentBuilder
}

// Send builds the followup message and sends
// it as response to the interaction.
func (b *FollowUpMessageBuilder) Send() *FollowUpMessage {
	if b.componentBuilder != nil {
		b.data.Components = append(b.data.Components, b.componentBuilder.components...)
	}

	fum := &FollowUpMessage{
		ken: b.ken,
		i:   b.i,
	}
	fum.Message, fum.Error = b.ken.s.FollowupMessageCreate(b.i, b.wait, b.data)
	if fum.HasError() {
		return fum
	}

	if b.componentBuilder != nil {
		b.componentBuilder.chanId = fum.ChannelID
		b.componentBuilder.msgId = fum.ID
		fum.unregisterComponentHandlers, fum.Error = b.componentBuilder.build()
	}

	return fum
}

// AddComponents is getting passed a builder function
// where you can attach message components and handlers
// which will be applied to the followup message when
// sent.
func (b *FollowUpMessageBuilder) AddComponents(cb func(*ComponentBuilder)) *FollowUpMessageBuilder {
	if b.componentBuilder == nil {
		b.componentBuilder = newBuilder(b.ken.componentHandler)
	}
	cb(b.componentBuilder)
	return b
}

// FollowUpMessage wraps an interaction follow
// up message and collected errors.
type FollowUpMessage struct {
	*discordgo.Message

	// Error contains the error instance of
	// error occurrences during method execution.
	Error error

	ken *Ken
	i   *discordgo.Interaction

	unregisterComponentHandlers func() error
}

// Edit overwrites the given follow up message with the
// data specified.
func (m *FollowUpMessage) Edit(data *discordgo.WebhookEdit) (err error) {
	if m.Error != nil {
		err = m.Error
		return
	}

	inter, err := m.ken.s.FollowupMessageEdit(m.i, m.ID, data)
	if err != nil {
		return
	}
	// This is done to avoid setting m.Message to nil when
	// an error is returned above.
	m.Message = inter
	return
}

// EditEmbed is shorthand for edit with the passed embed as
// WebhookEdit data.
func (m *FollowUpMessage) EditEmbed(emb *discordgo.MessageEmbed) (err error) {
	return m.Edit(&discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{emb},
	})
}

// Delete removes the follow up message.
func (m *FollowUpMessage) Delete() (err error) {
	if m.Error != nil {
		err = m.Error
		return
	}

	err = m.ken.s.FollowupMessageDelete(m.i, m.ID)
	return
}

// DeleteAfter queues a deletion of the follow up
// message after the specified duration.
func (m *FollowUpMessage) DeleteAfter(d time.Duration) *FollowUpMessage {
	go func() {
		time.Sleep(d)
		m.Delete()
	}()
	return m
}

// HasError returns true if the value of Error
// is not nil.
func (m *FollowUpMessage) HasError() bool {
	return m.Error != nil
}

// AddComponents returns a new component builder to add
// message components with handlers to the FollowUpMessage.
func (m *FollowUpMessage) AddComponents() *ComponentBuilder {
	return m.ken.Components().Add(m.ID, m.ChannelID)
}

// UnregisterComponentHandlers removes all handlers of
// attached componets from the register.
func (m *FollowUpMessage) UnregisterComponentHandlers() error {
	if m.unregisterComponentHandlers != nil {
		return m.unregisterComponentHandlers()
	}
	return nil
}
