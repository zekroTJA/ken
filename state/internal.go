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
	if c, err = s.State.Channel(id); err != nil && err != discordgo.ErrStateNotFound {
		return
	}
	if c == nil {
		c, err = s.Channel(id)
	}
	return
}

func (*Internal) Guild(s *discordgo.Session, id string) (g *discordgo.Guild, err error) {
	if g, err = s.State.Guild(id); err != nil && err != discordgo.ErrStateNotFound {
		return
	}
	if g == nil {
		g, err = s.Guild(id)
	}
	return
}

func (*Internal) Role(s *discordgo.Session, gID, id string) (r *discordgo.Role, err error) {
	if r, err = s.State.Role(gID, id); err != nil && err != discordgo.ErrStateNotFound {
		return
	}

	roles, err := s.GuildRoles(gID)
	if err != nil {
		return
	}

	for _, r = range roles {
		if r.ID == id {
			return
		}
	}
	return
}

func (*Internal) User(s *discordgo.Session, id string) (u *discordgo.User, err error) {
	u, err = s.User(id)
	return
}
