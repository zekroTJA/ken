package ken

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/safepool"
)

// AutocompleteCommand can be implemented by your command to enable autocomplete support
// for your command options.
type AutocompleteCommand interface {
	// Autocomplete will be called every time an autocomplete input event has veen received
	// for the registered command. It is getting passed a context which contains the event
	// data.
	//
	// Return the choises that shall be displayed or an error if something went wrong
	// during fetching the choises.
	//
	// The context object should not be used after the handler call has been completed.
	Autocomplete(ctx *AutocompleteContext) ([]*discordgo.ApplicationCommandOptionChoice, error)
}

// AutocompleteContext provides easy acces to the underlying event data.
type AutocompleteContext struct {
	session *discordgo.Session
	ken     *Ken
	event   *discordgo.InteractionCreate
}

var _ safepool.ResetState = (*AutocompleteContext)(nil)

func (t *AutocompleteContext) ResetState() {
	t.ken = nil
	t.session = nil
	t.event = nil
}

// Event returns the underlying InteractionCreate event.
func (t *AutocompleteContext) Event() *discordgo.InteractionCreate {
	return t.event
}

// GetSession returns the Discordgo Session instance.
func (t *AutocompleteContext) GetSession() *discordgo.Session {
	return t.session
}

// GetKen returns the Ken instance.
func (t *AutocompleteContext) GetKen() *Ken {
	return t.ken
}

// User returns the user object of the event caller. It may be nil if no user has been
// set to the event.
func (t *AutocompleteContext) User() (u *discordgo.User) {
	u = t.event.User
	if u == nil {
		u = t.event.Member.User
	}
	return u
}

// Member returns the user object of the event caller. It may be nil if no member has been
// set to the event.
func (t *AutocompleteContext) Member() (u *discordgo.Member) {
	return t.event.Member
}

// Channel tries to fetch the channel object from the contained
// channel ID using the specified state manager.
func (t *AutocompleteContext) Channel() (*discordgo.Channel, error) {
	return t.ken.opt.State.Channel(t.session, t.event.ChannelID)
}

// Guild tries to fetch the guild object from the contained
// guild ID using the specified state manager.
func (t *AutocompleteContext) Guild() (*discordgo.Guild, error) {
	return t.ken.opt.State.Guild(t.session, t.event.GuildID)
}

// GetData returns the ApplicationCommandInteractionData of the internal event.
func (t *AutocompleteContext) GetData() discordgo.ApplicationCommandInteractionData {
	return t.event.ApplicationCommandData()
}

// GetInput takes the name of a command option and returns the input value from
// the event for that option.
//
// If ok is false, no value could be found for the given option.
func (t *AutocompleteContext) GetInput(optionName string) (value string, ok bool) {
	data := t.GetData()
	for _, opt := range data.Options {
		if opt.Name == optionName {
			return opt.StringValue(), true
		}
	}

	return "", false
}
