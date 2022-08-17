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
type ComponentHandlerFunc func(ctx ComponentContext)

// ComponentHandler keeps a registry of component handler
// callbacks to be executed when a given component has
// been interacted with.
type ComponentHandler struct {
	ken            *Ken
	unregisterFunc func()

	mtx      sync.RWMutex
	handlers map[string]ComponentHandlerFunc

	ctxPool sync.Pool
}

// NewComponentHandler returns a new instance of
// ComponentHandler using the given instance of
// Ken.
func NewComponentHandler(ken *Ken) *ComponentHandler {
	var t ComponentHandler

	t.ken = ken
	t.handlers = make(map[string]ComponentHandlerFunc)
	t.unregisterFunc = t.ken.s.AddHandler(t.handle)
	t.ctxPool = sync.Pool{
		New: func() interface{} {
			return &ComponentCtx{}
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

func (t *ComponentHandler) handle(_ *discordgo.Session, e *discordgo.InteractionCreate) {
	if e.Type != discordgo.InteractionMessageComponent {
		return
	}

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
