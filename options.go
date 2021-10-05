package ken

import (
	"github.com/bwmarrin/discordgo"
)

// CommandOptions provides additional functionailities to
// an array of ApplicationCommandInteractionDataOptions.
type CommandOptions []*discordgo.ApplicationCommandInteractionDataOption

// Get safely returns an options from command options
// by index.
func (co CommandOptions) Get(i int) *discordgo.ApplicationCommandInteractionDataOption {
	if i < 0 {
		i = 0
	}
	if i >= len(co) {
		i = len(co) - 1
	}
	return co[i]
}

// Options returns wrapped underlying options
// of a sub command by ID.
func (co CommandOptions) Options(i int) CommandOptions {
	return co.Get(i).Options
}

// GetByNameOptional returns an option by name. If the option with the
// name does not exist, the returned value for ok is false.
//
// This should be used for non-required options.
func (co CommandOptions) GetByNameOptional(name string) (opt *discordgo.ApplicationCommandInteractionDataOption, ok bool) {
	for _, c := range co {
		if c.Name == name {
			ok = true
			opt = c
			break
		}
	}
	return
}

// GetByName returns an option by name.
//
// This should only be used on required options.
func (co CommandOptions) GetByName(name string) (opt *discordgo.ApplicationCommandInteractionDataOption) {
	opt, _ = co.GetByNameOptional(name)
	return
}
