package ken

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

// ComponentHandleFunc is the handler function for
// message component interactions. It is getting
// passed a ComponentContext which contians the
// interaction event data and can be used to respond
// to the interaction.
//
// A boolean is returned to indicate the success of
// the execution of the handler.
type ComponentHandlerFunc func(ctx ComponentContext) bool

type ModalHandlerFunc func(ctx ModalContext) bool

// ComponentHandler keeps a registry of component handler
// callbacks to be executed when a given component has
// been interacted with.
type ComponentHandler struct {
	ken            *Ken
	unregisterFunc func()

	mtx           sync.RWMutex
	handlers      map[string]ComponentHandlerFunc
	modalHandlers map[string]ModalHandlerFunc

	ctxPool      sync.Pool
	modalCtxPool sync.Pool
}

// NewComponentHandler returns a new instance of
// ComponentHandler using the given instance of
// Ken.
func NewComponentHandler(ken *Ken) *ComponentHandler {
	var t ComponentHandler

	t.ken = ken
	t.handlers = make(map[string]ComponentHandlerFunc)
	t.modalHandlers = make(map[string]ModalHandlerFunc)
	t.unregisterFunc = t.ken.s.AddHandler(t.handle)
	t.ctxPool = sync.Pool{
		New: func() interface{} {
			return &ComponentCtx{}
		},
	}
	t.modalCtxPool = sync.Pool{
		New: func() interface{} {
			return &ModalCtx{}
		},
	}

	return &t
}

// Add returns a new ComponentBuilder which attaches the
// added components to the given message by messageId and
// channelId on build.
func (t *ComponentHandler) Add(messageId, channelId string) *ComponentBuilder {
	return newBuilder(t, messageId, channelId)
}

// Register registers a raw ComponentHandlerFunc which is
// fired when the component with the specified customId has
// been interacted with.
//
// Te returned function unregisters the specified handler
// from the registry but does not remove the added message
// components from the message.
//
// Registering a handler twice on the smae customId
// overwrites the previously registered handler function.
func (t *ComponentHandler) Register(customId string, handler ComponentHandlerFunc) func() {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.handlers[customId] = handler

	return func() {
		t.Unregister(customId)
	}
}

// Unregister removes one or more handlers from the
// registry set to the given customId(s) of the
// message component(s).
func (t *ComponentHandler) Unregister(customId ...string) {
	if len(customId) == 0 {
		return
	}
	t.mtx.Lock()
	defer t.mtx.Unlock()
	for _, id := range customId {
		delete(t.handlers, id)
	}
}

// AppendToMessage edits the given message by messageId and
// channelId which adds the passed message components to
// the message.
func (t *ComponentHandler) AppendToMessage(
	messageId string,
	channelId string,
	components []discordgo.MessageComponent,
) error {
	_, err := t.ken.s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID:         messageId,
		Channel:    channelId,
		Components: components,
	})
	return err
}

// UnregisterDiscordHandler removes the Discord event handler
// function from the internal DiscordGo Session.
func (t *ComponentHandler) UnregisterDiscordHandler() {
	t.unregisterFunc()
}

func (t *ComponentHandler) registerModalHandler(customId string, handler ModalHandlerFunc) func() {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.modalHandlers[customId] = func(ctx ModalContext) bool {
		ok := handler(ctx)
		if ok {
			t.unregisterModalhandler(customId)
		}
		return ok
	}

	return func() {
		t.unregisterModalhandler(customId)
	}
}

func (t *ComponentHandler) unregisterModalhandler(customId ...string) {
	if len(customId) == 0 {
		return
	}
	t.mtx.Lock()
	defer t.mtx.Unlock()
	for _, id := range customId {
		delete(t.modalHandlers, id)
	}
}

func (t *ComponentHandler) handle(_ *discordgo.Session, e *discordgo.InteractionCreate) {
	switch e.Type {
	case discordgo.InteractionMessageComponent:
		t.handleMessageComponent(e)
	case discordgo.InteractionModalSubmit:
		t.handleModalSubmit(e)
	}
}

func (t *ComponentHandler) handleMessageComponent(e *discordgo.InteractionCreate) {
	data := e.MessageComponentData()

	t.mtx.RLock()
	handler, ok := t.handlers[data.CustomID]
	t.mtx.RUnlock()

	if !ok {
		return
	}

	ctx := t.ctxPool.Get().(*ComponentCtx)
	ctx.Data = data
	ctx.Ephemeral = false
	ctx.Event = e
	ctx.Session = t.ken.s
	ctx.Ken = t.ken
	ctx.responded = false

	defer func() {
		t.ctxPool.Put(ctx)
	}()

	handler(ctx)
}

func (t *ComponentHandler) handleModalSubmit(e *discordgo.InteractionCreate) {
	data := e.ModalSubmitData()

	t.mtx.RLock()
	handler, ok := t.modalHandlers[data.CustomID]
	t.mtx.RUnlock()

	if !ok {
		return
	}

	ctx := t.modalCtxPool.Get().(*ModalCtx)
	ctx.Data = data
	ctx.Ephemeral = false
	ctx.Event = e
	ctx.Session = t.ken.s
	ctx.Ken = t.ken
	ctx.responded = false

	defer func() {
		t.modalCtxPool.Put(ctx)
	}()

	handler(ctx)
}
