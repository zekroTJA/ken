package state

import "github.com/bwmarrin/discordgo"

// State defines an implementation of the
// state cache.
type State interface {
	// SelfUser returns the user objects of the
	// authenticated user.
	SelfUser(s *discordgo.Session) (*discordgo.User, error)

	// Channel returns a channel object by its ID, whether
	// from cache or fetched from the API when not stored
	// in the state chache.
	Channel(s *discordgo.Session, id string) (*discordgo.Channel, error)

	// Guild returns a guild object by its ID, whether
	// from cache or fetched from the API when not stored
	// in the state chache.
	Guild(s *discordgo.Session, id string) (*discordgo.Guild, error)
}
