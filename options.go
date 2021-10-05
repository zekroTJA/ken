package ken

import (
	"github.com/bwmarrin/discordgo"
)

// CommandOptions provides additional functionailities to
// an array of ApplicationCommandInteractionDataOptions.
type CommandOptions []*discordgo.ApplicationCommandInteractionDataOption

// Get safely returns an options from command options
// by index.
func (co CommandOptions) Get(i int) *CommandOption {
	if i < 0 {
		i = 0
	}
	if i >= len(co) {
		i = len(co) - 1
	}
	return &CommandOption{co[i]}
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
func (co CommandOptions) GetByNameOptional(name string) (opt *CommandOption, ok bool) {
	for _, c := range co {
		if c.Name == name {
			ok = true
			opt = &CommandOption{c}
			break
		}
	}
	return
}

// GetByName returns an option by name.
//
// This should only be used on required options.
func (co CommandOptions) GetByName(name string) (opt *CommandOption) {
	opt, _ = co.GetByNameOptional(name)
	return
}

// CommandOption wraps a ApplicationCommandInteractionDataOption
// to provide additional functionalities and method overrides.
type CommandOption struct {
	*discordgo.ApplicationCommandInteractionDataOption
}

// ChannelValue is a utility function for casting option value to channel object.
//
// The object is taken from the specified state instance.
func (o *CommandOption) ChannelValue(ctx *Ctx) *discordgo.Channel {
	if o.Type != discordgo.ApplicationCommandOptionChannel {
		panic("ChannelValue called on data option of type " + o.Type.String())
	}
	chanID := o.Value.(string)

	if ctx == nil {
		return &discordgo.Channel{ID: chanID}
	}

	ch, err := ctx.k.opt.State.Channel(ctx.Session, chanID)
	if err != nil {
		return &discordgo.Channel{ID: chanID}
	}

	return ch
}

// RoleValue is a utility function for casting option value to role object.
//
// The object is taken from the specified state instance.
func (o *CommandOption) RoleValue(ctx *Ctx) *discordgo.Role {
	if o.Type != discordgo.ApplicationCommandOptionRole {
		panic("RoleValue called on data option of type " + o.Type.String())
	}
	roleID := o.Value.(string)

	if ctx == nil {
		return &discordgo.Role{ID: roleID}
	}

	role, err := ctx.k.opt.State.Role(ctx.Session, ctx.Event.GuildID, roleID)
	if err != nil {
		return &discordgo.Role{ID: roleID}
	}

	return role
}

// UserValue is a utility function for casting option value to user object.
//
// The object is taken from the specified state instance.
func (o *CommandOption) UserValue(ctx *Ctx) *discordgo.User {
	if o.Type != discordgo.ApplicationCommandOptionUser {
		panic("UserValue called on data option of type " + o.Type.String())
	}
	userID := o.Value.(string)

	if ctx == nil {
		return &discordgo.User{ID: userID}
	}

	user, err := ctx.k.opt.State.User(ctx.Session, userID)
	if err != nil {
		return &discordgo.User{ID: userID}
	}

	return user
}
