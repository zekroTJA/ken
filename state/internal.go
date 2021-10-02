package state

import "github.com/bwmarrin/discordgo"

var _ State = (*Internal)(nil)

// Internal implements the state Interface for
// the internal discordgo.State instance.
type Internal struct{}

// NewInternal returns a new instance of Internal.
func NewInternal() *Internal {
	return &Internal{}
}

func (*Internal) SelfUser(s *discordgo.Session) (u *discordgo.User, err error) {
	u = s.State.User
	return
}

func (*Internal) Channel(s *discordgo.Session, id string) (c *discordgo.Channel, err error) {
	if c, err = s.State.Channel(id); err != nil {
		return
	}
	if c == nil {
		c, err = s.Channel(id)
	}
	return
}

func (*Internal) Guild(s *discordgo.Session, id string) (g *discordgo.Guild, err error) {
	if g, err = s.State.Guild(id); err != nil {
		return
	}
	if g == nil {
		g, err = s.Guild(id)
	}
	return
}
