package ken

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// FollowUpMessage wraps an interaction follow
// up message and collected errors.
type FollowUpMessage struct {
	*discordgo.Message

	// Error contains the error instance of
	// error occurrences during method execution.
	Error error

	self *discordgo.User
	s    *discordgo.Session
	i    *discordgo.Interaction
}

// Edit overwrites the given follow up message with the
// data specified.
func (m *FollowUpMessage) Edit(data *discordgo.WebhookEdit) (err error) {
	if m.Error != nil {
		err = m.Error
		return
	}

	inter, err := m.s.FollowupMessageEdit(m.self.ID, m.i, m.ID, data)
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
		Embeds: []*discordgo.MessageEmbed{emb},
	})
}

// Delete removes the follow up message.
func (m *FollowUpMessage) Delete() (err error) {
	if m.Error != nil {
		err = m.Error
		return
	}

	err = m.s.FollowupMessageDelete(m.self.ID, m.i, m.ID)
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
