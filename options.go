package ken

import (
	"strings"

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
func (o *CommandOption) ChannelValue(ctx Context) *discordgo.Channel {
	if o.Type != discordgo.ApplicationCommandOptionChannel {
		panic("ChannelValue called on data option of type " + o.Type.String())
	}
	chanID := o.Value.(string)

	if ctx == nil {
		return &discordgo.Channel{ID: chanID}
	}

	ch, err := ctx.GetKen().opt.State.Channel(ctx.GetSession(), chanID)
	if err != nil {
		return &discordgo.Channel{ID: chanID}
	}

	return ch
}

// RoleValue is a utility function for casting option value to role object.
//
// The object is taken from the specified state instance.
func (o *CommandOption) RoleValue(ctx Context) *discordgo.Role {
	if o.Type != discordgo.ApplicationCommandOptionRole {
		panic("RoleValue called on data option of type " + o.Type.String())
	}
	roleID := o.Value.(string)

	if ctx == nil {
		return &discordgo.Role{ID: roleID}
	}

	role, err := ctx.GetKen().opt.State.Role(ctx.GetSession(), ctx.GetEvent().GuildID, roleID)
	if err != nil {
		return &discordgo.Role{ID: roleID}
	}

	return role
}

// UserValue is a utility function for casting option value to user object.
//
// The object is taken from the specified state instance.
func (o *CommandOption) UserValue(ctx Context) *discordgo.User {
	if o.Type != discordgo.ApplicationCommandOptionUser {
		panic("UserValue called on data option of type " + o.Type.String())
	}
	userID := o.Value.(string)

	if ctx == nil {
		return &discordgo.User{ID: userID}
	}

	user, err := ctx.GetKen().opt.State.User(ctx.GetSession(), userID)
	if err != nil {
		return &discordgo.User{ID: userID}
	}

	return user
}

// StringValue is a utility function for casting option value to string.
//
// Because you can not pass multiline string entries to slash commands,
// this converts `\n` in a message to an actual line break.
func (o *CommandOption) StringValue() (v string) {
	v = o.ApplicationCommandInteractionDataOption.StringValue()
	v = strings.ReplaceAll(v, "\\n", "\n")
	return
}
