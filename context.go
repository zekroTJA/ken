package ken

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/xid"
)

// ContextResponder defines the implementation of an
// interaction context with functionalities to respond
// to the interaction, to set the ephemeral state and
// to retrieve the nested session and event.
type ContextResponder interface {

	// Respond to an interaction event with the given
	// interaction response payload.
	//
	// When an interaction has already been responded to,
	// the response will be edited instead on execution.
	Respond(r *discordgo.InteractionResponse) (err error)

	// RespondEmbed is shorthand for Respond with an
	// embed payload as passed.
	RespondEmbed(emb *discordgo.MessageEmbed) (err error)

	// RespondError is shorthand for RespondEmbed with an
	// error embed as message with the passed content and
	// title.
	RespondError(content, title string) (err error)

	// FollowUp creates a follow up message to the
	// interaction event and returns a FollowUpMessage
	// object containing the created message as well as
	// an error instance, if an error occurred.
	//
	// This way it allows to be chained in one call with
	// subsequent FollowUpMessage method calls.
	FollowUp(wait bool, data *discordgo.WebhookParams) (fumb *FollowUpMessageBuilder)

	// FollowUpEmbed is shorthand for FollowUp with an
	// embed payload as passed.
	FollowUpEmbed(emb *discordgo.MessageEmbed) (fumb *FollowUpMessageBuilder)

	// FollowUpError is shorthand for FollowUpEmbed with an
	// error embed as message with the passed content and
	// title.
	FollowUpError(content, title string) (fumb *FollowUpMessageBuilder)

	// Defer is shorthand for Respond with an InteractionResponse
	// of the type InteractionResponseDeferredChannelMessageWithSource.
	//
	// It should be used when the interaction response can not be
	// instantly returned.
	Defer() (err error)

	// GetEphemeral returns the current emphemeral state
	// of the command invokation.
	GetEphemeral() bool

	// SetEphemeral sets the emphemeral state of the command
	// invokation.
	//
	// Ephemeral can be set to true which will
	// send all subsequent command responses
	// only to the user which invoked the command.
	SetEphemeral(v bool)

	// GetSession returns the current Discordgo session instance.
	GetSession() *discordgo.Session

	// GetEvent returns the InteractionCreate event instance which
	// invoked the interaction command.
	GetEvent() *discordgo.InteractionCreate

	// User returns the User object of the executor either from
	// the events User object or from the events Member object.
	User() (u *discordgo.User)
}

// Context defines the implementation of an interaction
// command context passed to the command handler.
type Context interface {
	ContextResponder
	ObjectProvider

	// Channel tries to fetch the channel object from the contained
	// channel ID using the specified state manager.
	Channel() (*discordgo.Channel, error)

	// Channel tries to fetch the guild object from the contained
	// guild ID using the specified state manager.
	Guild() (*discordgo.Guild, error)

	// Options returns the application command data options
	// with additional functionality methods.
	Options() CommandOptions

	// SlashCommand returns the contexts Command as a
	// SlashCommand interface.
	SlashCommand() (cmd SlashCommand, ok bool)

	// UserCommand returns the contexts Command as a
	// UserCommand interface.
	UserCommand() (cmd UserCommand, ok bool)

	// MessageCommand returns the contexts Command as a
	// MessageCommand interface.
	MessageCommand() (cmd MessageCommand, ok bool)

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
	HandleSubCommands(handler ...SubCommandHandler) (err error)

	// GetKen returns the root instance of Ken.
	GetKen() *Ken

	// GetCommand returns the command instance called.
	GetCommand() Command
}

// ctxResponder provides functionailities to respond
// to an interaction.
type ctxResponder struct {
	responded bool
	ken       *Ken
	session   *discordgo.Session
	event     *discordgo.InteractionCreate
	ephemeral bool
}

var _ ContextResponder = (*ctxResponder)(nil)

func (c *ctxResponder) Respond(r *discordgo.InteractionResponse) (err error) {
	if r.Data == nil {
		r.Data = new(discordgo.InteractionResponseData)
	}
	r.Data.Flags = c.messageFlags(r.Data.Flags)
	if c.responded {
		if r == nil || r.Data == nil {
			return
		}
		_, err = c.GetSession().InteractionResponseEdit(c.event.Interaction, &discordgo.WebhookEdit{
			Content:         &r.Data.Content,
			Embeds:          &r.Data.Embeds,
			Components:      &r.Data.Components,
			Files:           r.Data.Files,
			AllowedMentions: r.Data.AllowedMentions,
		})
	} else {
		err = c.GetSession().InteractionRespond(c.event.Interaction, r)
		c.responded = err == nil
		if err != nil {
			_ = err
		}
	}
	return
}

func (c *ctxResponder) RespondEmbed(emb *discordgo.MessageEmbed) (err error) {
	if emb.Color <= 0 {
		emb.Color = c.ken.opt.EmbedColors.Default
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

func (c *ctxResponder) RespondError(content, title string) (err error) {
	return c.RespondEmbed(&discordgo.MessageEmbed{
		Description: content,
		Title:       title,
		Color:       c.ken.opt.EmbedColors.Error,
	})
}

func (c *ctxResponder) FollowUp(wait bool, data *discordgo.WebhookParams) (fumb *FollowUpMessageBuilder) {
	data.Flags = c.messageFlags(data.Flags)
	return &FollowUpMessageBuilder{
		ken:  c.ken,
		i:    c.event.Interaction,
		data: data,
		wait: wait,
	}
}

func (c *ctxResponder) FollowUpEmbed(emb *discordgo.MessageEmbed) (fumb *FollowUpMessageBuilder) {
	if emb.Color <= 0 {
		emb.Color = c.ken.opt.EmbedColors.Default
	}
	return c.FollowUp(true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			emb,
		},
	})
}

func (c *ctxResponder) FollowUpError(content, title string) (fumb *FollowUpMessageBuilder) {
	return c.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: content,
		Title:       title,
		Color:       c.ken.opt.EmbedColors.Error,
	})
}

func (c *ctxResponder) Defer() (err error) {
	err = c.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	return
}

func (c *ctxResponder) GetEphemeral() bool {
	return c.ephemeral
}

func (c *ctxResponder) SetEphemeral(v bool) {
	c.ephemeral = v
}

func (c *ctxResponder) GetSession() *discordgo.Session {
	return c.session
}

func (c *ctxResponder) GetEvent() *discordgo.InteractionCreate {
	return c.event
}

func (c *ctxResponder) messageFlags(p discordgo.MessageFlags) (f discordgo.MessageFlags) {
	f = p
	if c.ephemeral {
		f |= discordgo.MessageFlagsEphemeral
	}
	return
}

// User returns the User object of the executor either from
// the events User object or from the events Member object.
func (c *ctxResponder) User() (u *discordgo.User) {
	u = c.event.User
	if u == nil && c.event.Member != nil {
		u = c.event.Member.User
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
	ctxResponder

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
	if v = c.ObjectMap.Get(key); v == nil && c.ken.opt.DependencyProvider != nil {
		v = c.ken.opt.DependencyProvider.Get(key)
	}
	return
}

// Channel tries to fetch the channel object from the contained
// channel ID using the specified state manager.
func (c *Ctx) Channel() (*discordgo.Channel, error) {
	return c.ken.opt.State.Channel(c.session, c.event.ChannelID)
}

// Channel tries to fetch the guild object from the contained
// guild ID using the specified state manager.
func (c *Ctx) Guild() (*discordgo.Guild, error) {
	return c.ken.opt.State.Guild(c.session, c.event.GuildID)
}

// Options returns the application command data options
// with additional functionality methods.
func (c *Ctx) Options() CommandOptions {
	return c.event.ApplicationCommandData().Options
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
	Run  func(ctx SubCommandContext) error
}

// SubCommandContext wraps the current command
// Context and with the called sub command name
// and scopes the command options to the
// options of the called sub command.
//
// The SubCommandCtx must not be stored or used
// after command execution.
type SubCommandContext interface {
	Context

	// GetSubCommandName returns the sub command
	// name which has been invoked.
	GetSubCommandName() string
}

type subCommandCtx struct {
	*Ctx

	subCommandName string
}

var _ SubCommandContext = (*subCommandCtx)(nil)

// Options returns the options array of the called
// sub command.
func (c *subCommandCtx) Options() CommandOptions {
	return c.Ctx.Options().GetByName(c.subCommandName).Options
}

func (c *subCommandCtx) GetSubCommandName() string {
	return c.subCommandName
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

		ctx := c.ken.subCtxPool.Get().(*subCommandCtx)
		ctx.Ctx = c
		ctx.subCommandName = h.Name
		err = h.Run(ctx)
		c.ken.subCtxPool.Put(ctx)
		break
	}
	return
}

// GetKen returns the root instance of Ken.
func (c *Ctx) GetKen() *Ken {
	return c.ken
}

// GetCommand returns the command instance called.
func (c *Ctx) GetCommand() Command {
	return c.Command
}

// ComponentContext gives access to the underlying
// MessageComponentInteractionData and gives the
// ability to open a Modal afterwards.
type ComponentContext interface {
	ContextResponder

	// GetData returns the underlying
	// MessageComponentInteractionData.
	GetData() discordgo.MessageComponentInteractionData

	// OpenModal opens a new modal with the given
	// title, content and components built with the
	// passed build function. A channel is returned
	// which will receive a ModalContext when the user
	// has interacted with the modal.
	OpenModal(
		title string,
		content string,
		build func(b ComponentAssembler),
	) (<-chan ModalContext, error)
}

type componentCtx struct {
	ctxResponder

	Data discordgo.MessageComponentInteractionData
}

var _ ComponentContext = (*componentCtx)(nil)

func (c *componentCtx) GetData() discordgo.MessageComponentInteractionData {
	return c.Data
}

func (c *componentCtx) OpenModal(
	title string,
	content string,
	build func(b ComponentAssembler),
) (<-chan ModalContext, error) {
	b := newComponentAssembler()
	build(b)

	modalId := xid.New().String()
	err := c.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID:   modalId,
			Title:      title,
			Content:    content,
			Components: b.components,
		},
	})
	if err != nil {
		return nil, err
	}

	cCtx := make(chan ModalContext, 1)

	c.ken.componentHandler.registerModalHandler(modalId, func(ctx ModalContext) bool {
		cCtx <- ctx
		return true
	})

	return cCtx, nil
}

// ModalContext provides access to the underlying
// ModalSubmitInteractionData and some utility
// methods to access component data from the
// response.
type ModalContext interface {
	ContextResponder

	// GetData returns the underlying
	// ModalSubmitInteractionData.
	GetData() discordgo.ModalSubmitInteractionData

	// GetComponentByID tries to find a message component
	// by CustomID in the response data and returns it
	// wrapped into MessageComponent.
	//
	// The returned MessageComponent will contain a nil
	// value for the wrapped discordgo.MessageComponent
	// if it could not be found in the response.
	//
	// Subsequent method calls to MessageComponent will
	// not fail though to ensure the ability to chain
	// method calls.
	GetComponentByID(customId string) MessageComponent
}

type modalCtx struct {
	ctxResponder

	Data discordgo.ModalSubmitInteractionData
}

var _ ModalContext = (*modalCtx)(nil)

func (c modalCtx) GetData() discordgo.ModalSubmitInteractionData {
	return c.Data
}

func (c modalCtx) GetComponentByID(customId string) MessageComponent {
	return MessageComponent{getComponentByID(customId, c.GetData().Components)}
}

func getComponentByID(
	customId string,
	comps []discordgo.MessageComponent,
) discordgo.MessageComponent {
	for _, comp := range comps {
		if row, ok := comp.(*discordgo.ActionsRow); ok {
			found := getComponentByID(customId, row.Components)
			if found != nil {
				return found
			}
		}
		if customId == getCustomId(comp) {
			return comp
		}
	}
	return nil
}
