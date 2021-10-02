package ken

import "github.com/bwmarrin/discordgo"

type FollowUpMessage struct {
	*discordgo.Message

	self string
	s    *discordgo.Session
	i    *discordgo.Interaction
}

func (m *FollowUpMessage) Edit(data *discordgo.WebhookEdit) (err error) {
	inter, err := m.s.FollowupMessageEdit(m.self, m.i, m.ID, data)
	if err != nil {
		return
	}
	// This is done to avoid setting m.Message to nil when
	// an error is returned above.
	m.Message = inter
	return
}

func (m *FollowUpMessage) Delete() (err error) {
	err = m.s.FollowupMessageDelete(m.self, m.i, m.ID)
	return
}
