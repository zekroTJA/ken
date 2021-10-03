package ken

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type FollowUpMessage struct {
	*discordgo.Message

	Error error

	self *discordgo.User
	s    *discordgo.Session
	i    *discordgo.Interaction
}

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

func (m *FollowUpMessage) Delete() (err error) {
	if m.Error != nil {
		err = m.Error
		return
	}

	err = m.s.FollowupMessageDelete(m.self.ID, m.i, m.ID)
	return
}

func (m *FollowUpMessage) DeleteAfter(d time.Duration) *FollowUpMessage {
	go func() {
		time.Sleep(d)
		m.Delete()
	}()
	return m
}
