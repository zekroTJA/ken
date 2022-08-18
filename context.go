package ken

import (
	"github.com/bwmarrin/discordgo"
)

// ContextResponder defines the implementation of an
// interaction context with functionalities to respond
// to the interaction, to set the ephemeral state and
// to retrieve the nested session and event.
type ContextResponder interface {
	Respond(r *discordgo.InteractionResponse) (err error)
	RespondEmbed(emb *discordgo.MessageEmbed) (err error)
	RespondError(content, title string) (err error)
	FollowUp(wait bool, data *discordgo.WebhookParams) (fum *FollowUpMessage)
	FollowUpEmbed(emb *discordgo.MessageEmbed) (fum *FollowUpMessage)
	FollowUpError(content, title string) (fum *FollowUpMessage)
	Defer() (err error)
	GetEphemeral() bool
	SetEphemeral(v bool)
	GetSession() *discordgo.Session
	GetEvent() *discordgo.InteractionCreate
}

// Context defines the implementation of an interaction
// command context passed to the command handler.
type Context interface {
	ContextResponder

	Get(key string) (v interface{})
	Channel() (*discordgo.Channel, error)
	Guild() (*discordgo.Guild, error)
	User() (u *discordgo.User)
	Options() CommandOptions
	SlashCommand() (cmd SlashCommand, ok bool)
	UserCommand() (cmd UserCommand, ok bool)
	MessageCommand() (cmd MessageCommand, ok bool)
	HandleSubCommands(handler ...SubCommandHandler) (err error)
	GetKen() *Ken
	GetCommand() Command
}

// CtxResponder provides functionailities to respond
// to an interaction.
type CtxResponder struct {
	responded bool

	// Ken keeps a reference to the main Ken instance.
	Ken *Ken
	// Session holds the discordgo session instance.
	Session *discordgo.Session
	// Event provides the InteractionCreate event
	// instance.
	Event *discordgo.InteractionCreate
	// Ephemeral can be set to true which will
	// send all subsequent command responses
	// only to the user which invoked the command.
	Ephemeral bool
}

var _ ContextResponder = (*CtxResponder)(nil)

// Respond to an interaction event with the given
// interaction response payload.
//
// When an interaction has already been responded to,
// the response will be edited instead on execution.
func (c *CtxResponder) Respond(r *discordgo.InteractionResponse) (err error) {
	if r.Data == nil {
		r.Data = new(discordgo.InteractionResponseData)
	}
	r.Data.Flags = c.messageFlags(r.Data.Flags)
	if c.responded {
		if r == nil || r.Data == nil {
			return
		}
		_, err = c.Session.InteractionResponseEdit(c.Event.Interaction, &discordgo.WebhookEdit{
			Content:         &r.Data.Content,
			Embeds:          &r.Data.Embeds,
			Components:      &r.Data.Components,
			Files:           r.Data.Files,
			AllowedMentions: r.Data.AllowedMentions,
		})
	} else {
		err = c.Session.InteractionRespond(c.Event.Interaction, r)
		c.responded = err == nil
		if err != nil {
			_ = err
		}
	}
	return
}

// RespondEmbed is shorthand for Respond with an
// embed payload as passed.
func (c *CtxResponder) RespondEmbed(emb *discordgo.MessageEmbed) (err error) {
	if emb.Color <= 0 {
		emb.Color = c.Ken.opt.EmbedColors.Default
	}
	return c.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				emb,
			},
		},
	})
}

// RespondError is shorthand for RespondEmbed with an
// error embed as message with the passed content and
// title.
func (c *CtxResponder) RespondError(content, title string) (err error) {
	return c.RespondEmbed(&discordgo.MessageEmbed{
		Description: content,
		Title:       title,
		Color:       c.Ken.opt.EmbedColors.Error,
	})
}

// FollowUp creates a follow up message to the
// interaction event and returns a FollowUpMessage
// object containing the created message as well as
// an error instance, if an error occurred.
//
// This way it allows to be chained in one call with
// subsequent FollowUpMessage method calls.
func (c *CtxResponder) FollowUp(wait bool, data *discordgo.WebhookParams) (fum *FollowUpMessage) {
	data.Flags = c.messageFlags(data.Flags)
	fum = &FollowUpMessage{
		ken: c.Ken,
		i:   c.Event.Interaction,
	}
	fum.Message, fum.Error = c.Session.FollowupMessageCreate(c.Event.Interaction, wait, data)
	return
}

// FollowUpEmbed is shorthand for FollowUp with an
// embed payload as passed.
func (c *CtxResponder) FollowUpEmbed(emb *discordgo.MessageEmbed) (fum *FollowUpMessage) {
	if emb.Color <= 0 {
		emb.Color = c.Ken.opt.EmbedColors.Default
	}
	return c.FollowUp(true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			emb,
		},
	})
}

// FollowUpError is shorthand for FollowUpEmbed with an
// error embed as message with the passed content and
// title.
func (c *CtxResponder) FollowUpError(content, title string) (fum *FollowUpMessage) {
	return c.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: content,
		Title:       title,
		Color:       c.Ken.opt.EmbedColors.Error,
	})
}

// Defer is shorthand for Respond with an InteractionResponse
// of the type InteractionResponseDeferredChannelMessageWithSource.
//
// It should be used when the interaction response can not be
// instantly returned.
func (c *CtxResponder) Defer() (err error) {
	err = c.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	return
}

// GetEphemeral returns the current emphemeral state
// of the command invokation.
func (c *CtxResponder) GetEphemeral() bool {
	return c.Ephemeral
}

// SetEphemeral sets the emphemeral state of the command
// invokation.
//
// Ephemeral can be set to true which will
// send all subsequent command responses
// only to the user which invoked the command.
func (c *CtxResponder) SetEphemeral(v bool) {
	c.Ephemeral = v
}

// GetSession returns the current Discordgo session instance.
func (c *CtxResponder) GetSession() *discordgo.Session {
	return c.Session
}

// GetEvent returns the InteractionCreate event instance which
// invoked the interaction command.
func (c *CtxResponder) GetEvent() *discordgo.InteractionCreate {
	return c.Event
}

func (c *CtxResponder) messageFlags(p discordgo.MessageFlags) (f discordgo.MessageFlags) {
	f = p
	if c.Ephemeral {
		f |= discordgo.MessageFlagsEphemeral
	}
	return
}

// Ctx holds the invokation context of
// a command.
//
// The Ctx must not be stored or used
// after command execution.
type Ctx struct {
	ObjectMap
	CtxResponder

	// Command provides the called command instance.
	Command Command
}

var _ Context = (*Ctx)(nil)

func newCtx() *Ctx {
	return &Ctx{
		ObjectMap: make(simpleObjectMap),
	}
}

// Get either returns an instance from the internal object map -
// if existent. Otherwise, the object is looked up in the specified
// dependency provider, if available. When no object was found in
// either of both maps, nil is returned.
func (c *Ctx) Get(key string) (v interface{}) {
	if v = c.ObjectMap.Get(key); v == nil && c.Ken.opt.DependencyProvider != nil {
		v = c.Ken.opt.DependencyProvider.Get(key)
	}
	return
}

// Channel tries to fetch the channel object from the contained
// channel ID using the specified state manager.
func (c *Ctx) Channel() (*discordgo.Channel, error) {
	return c.Ken.opt.State.Channel(c.Session, c.Event.ChannelID)
}

// Channel tries to fetch the guild object from the contained
// guild ID using the specified state manager.
func (c *Ctx) Guild() (*discordgo.Guild, error) {
	return c.Ken.opt.State.Guild(c.Session, c.Event.GuildID)
}

// User returns the User object of the executor either from
// the events User object or from the events Member object.
func (c *Ctx) User() (u *discordgo.User) {
	u = c.Event.User
	if u == nil && c.Event.Member != nil {
		u = c.Event.Member.User
	}
	return
}

// Options returns the application command data options
// with additional functionality methods.
func (c *Ctx) Options() CommandOptions {
	return c.Event.ApplicationCommandData().Options
}

// SlashCommand returns the contexts Command as a
// SlashCommand interface.
func (c *Ctx) SlashCommand() (cmd SlashCommand, ok bool) {
	cmd, ok = c.Command.(SlashCommand)
	return
}

// UserCommand returns the contexts Command as a
// UserCommand interface.
func (c *Ctx) UserCommand() (cmd UserCommand, ok bool) {
	cmd, ok = c.Command.(UserCommand)
	return
}

// MessageCommand returns the contexts Command as a
// MessageCommand interface.
func (c *Ctx) MessageCommand() (cmd MessageCommand, ok bool) {
	cmd, ok = c.Command.(MessageCommand)
	return
}

// SubCommandHandler is the handler function used
// to handle sub command calls.
type SubCommandHandler struct {
	Name string
	Run  func(ctx *SubCommandCtx) error
}

// SubCommandCtx wraps the current command Ctx and
// with the called sub command name and scopes the
// command options to the options of the called
// sub command.
//
// The SubCommandCtx must not be stored or used
// after command execution.
type SubCommandCtx struct {
	*Ctx

	SubCommandName string
}

var _ Context = (*SubCommandCtx)(nil)

// Options returns the options array of the called
// sub command.
func (c *SubCommandCtx) Options() CommandOptions {
	return c.Ctx.Options().GetByName(c.SubCommandName).Options
}

// HandleSubCommands takes a list of sub command handles.
// When the command is executed, the options are scanned
// for the sib command calls by their names. If one of
// the registered sub commands has been called, the specified
// handler function is executed.
//
// If the call occured, the passed handler function is
// getting passed the scoped sub command Ctx.
//
// The SubCommandCtx passed must not be stored or used
// after command execution.
func (c *Ctx) HandleSubCommands(handler ...SubCommandHandler) (err error) {
	for _, h := range handler {
		opt := c.Options().Get(0)
		if opt.Type != discordgo.ApplicationCommandOptionSubCommand || opt.Name != h.Name {
			continue
		}

		ctx := c.Ken.subCtxPool.Get().(*SubCommandCtx)
		ctx.Ctx = c
		ctx.SubCommandName = h.Name
		err = h.Run(ctx)
		c.Ken.subCtxPool.Put(ctx)
		break
	}
	return
}

// GetKen returns the root instance of Ken.
func (c *Ctx) GetKen() *Ken {
	return c.Ken
}

// GetCommand returns the command instance called.
func (c *Ctx) GetCommand() Command {
	return c.Command
}

type ComponentContext interface {
	ContextResponder

	GetData() discordgo.MessageComponentInteractionData
}

type ComponentCtx struct {
	CtxResponder

	Data discordgo.MessageComponentInteractionData
}

var _ ComponentContext = (*ComponentCtx)(nil)

func (c *ComponentCtx) GetData() discordgo.MessageComponentInteractionData {
	return c.Data
}
