package ken

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken/util"
)

// ComponentAssembler helps to build the message
// component tree.
type ComponentAssembler interface {

	// AddActionsRow adds an Action Row component to
	// the message. Use the builder passed by the
	// build function to assemble the components of
	// the Action Row.
	//
	// If you pass once as `true`, after the first
	// interaction inside the Action Row, all handlers
	// of the Action Row children are removed as well
	// as the Action Row component itself from the message.
	AddActionsRow(build func(b ComponentAssembler), once ...bool) ComponentAssembler

	// Add appends the passed message component to the
	// message with the given handler called on
	// interaction with the component.
	//
	// If you pass once as `true`, the handler is
	// removed after interaction with the component
	// as well as the component itself from the message.
	Add(
		component discordgo.MessageComponent,
		handler ComponentHandlerFunc,
		once ...bool,
	) ComponentAssembler
}

type handlerWrapper struct {
	handler   ComponentHandlerFunc
	once      bool
	onceGroup []string
}

type componentAssembler struct {
	components []discordgo.MessageComponent
	handlers   map[string]handlerWrapper
}

func newComponentAssembler() *componentAssembler {
	return &componentAssembler{
		handlers: make(map[string]handlerWrapper),
	}
}

func (t *componentAssembler) Add(
	component discordgo.MessageComponent,
	handler ComponentHandlerFunc,
	once ...bool,
) ComponentAssembler {
	t.components = append(t.components, component)

	customId := getCustomId(component)

	if customId == "" {
		return t
	}

	t.handlers[customId] = handlerWrapper{
		handler: handler,
		once:    len(once) != 0 && once[0],
	}

	return t
}

func (t *componentAssembler) AddActionsRow(build func(b ComponentAssembler), once ...bool) ComponentAssembler {
	b := newComponentAssembler()
	build(b)

	var onceGroup []string

	if len(once) != 0 && once[0] {
		for id := range b.handlers {
			onceGroup = append(onceGroup, id)
		}
	}

	for id, handler := range b.handlers {
		handler.onceGroup = onceGroup
		t.handlers[id] = handler
	}

	t.components = append(t.components, discordgo.ActionsRow{
		Components: b.components,
	})

	return t
}

// ComponentBuilder helps to build the message component
// tree, attaches the components to the given message
// and registers the interaction handlers for the given
// components.
type ComponentBuilder struct {
	ch     *ComponentHandler
	msgId  string
	chanId string

	condition ComponentHandlerFunc

	*componentAssembler
}

func newBuilderAttach(ch *ComponentHandler, msgId, chanId string) *ComponentBuilder {
	return &ComponentBuilder{
		ch:                 ch,
		msgId:              msgId,
		chanId:             chanId,
		componentAssembler: newComponentAssembler(),
	}
}

func newBuilder(ch *ComponentHandler) *ComponentBuilder {
	return &ComponentBuilder{
		ch:                 ch,
		componentAssembler: newComponentAssembler(),
	}
}

// Add appends the passed message component to the
// message with the given handler called on
// interaction with the component.
//
// If you pass once as `true`, the handler is
// removed after interaction with the component
// as well as the component itself from the message.
func (t *ComponentBuilder) Add(
	component discordgo.MessageComponent,
	handler ComponentHandlerFunc,
	once ...bool,
) *ComponentBuilder {
	t.componentAssembler.Add(component, handler, once...)
	return t
}

// AddActionsRow adds an Action Row component to
// the message. Use the builder passed by the
// build function to assemble the components of
// the Action Row.
//
// If you pass once as `true`, after the first
// interaction inside the Action Row, all handlers
// of the Action Row children are removed as well
// as the Action Row component itself from the message.
func (t *ComponentBuilder) AddActionsRow(build func(b ComponentAssembler), once ...bool) *ComponentBuilder {
	t.componentAssembler.AddActionsRow(build, once...)
	return t
}

// Condition sets a condition handler which needs to
// be met so that the component handler is activated.
func (t *ComponentBuilder) Condition(cond ComponentHandlerFunc) *ComponentBuilder {
	t.condition = cond
	return t
}

// Build attaches the registered messgae components to
// the specified message and registers the interaction
// handlers to the handler registry.
//
// It returns an unregister function which can be called
// to remove all message components appendet and and all
// interaction handler registered with this builder.
func (t *ComponentBuilder) Build() (unreg func() error, err error) {
	err = t.ch.AppendToMessage(t.msgId, t.chanId, t.components)
	if err != nil {
		return unreg, err
	}
	return t.build()
}

func (t *ComponentBuilder) build() (unreg func() error, err error) {
	t.ch.mtx.Lock()
	defer t.ch.mtx.Unlock()

	for key := range t.handlers {
		handler := t.handlers[key]

		if t.condition == nil {
			t.condition = func(ctx ComponentContext) bool { return true }
		}

		if len(handler.onceGroup) > 0 {
			t.ch.handlers[key] = func(ctx ComponentContext) bool {
				if !t.condition(ctx) || !handler.handler(ctx) {
					return false
				}

				t.components = []discordgo.MessageComponent{}
				t.ch.ken.s.ChannelMessageEditComplex(&discordgo.MessageEdit{
					ID:         t.msgId,
					Channel:    t.chanId,
					Components: &t.components,
				})
				kRems := make([]string, 0, len(handler.onceGroup))
				for _, kRem := range handler.onceGroup {
					kRems = append(kRems, kRem)
				}
				t.ch.Unregister(kRems...)
				return true
			}
		} else if handler.once {
			k := key // copy key for anonymous function
			t.ch.handlers[key] = func(ctx ComponentContext) bool {
				if !t.condition(ctx) || !handler.handler(ctx) {
					return false
				}

				t.components = removeComponentRecursive(t.components, k)
				t.ch.ken.s.ChannelMessageEditComplex(&discordgo.MessageEdit{
					ID:         t.msgId,
					Channel:    t.chanId,
					Components: &t.components,
				})

				t.ch.Unregister(k)
				return true
			}
		} else {
			t.ch.handlers[key] = func(ctx ComponentContext) bool {
				return t.condition(ctx) && handler.handler(ctx)
			}
		}
	}

	unreg = func() error {
		_, err := t.ch.ken.s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			ID:         t.msgId,
			Channel:    t.chanId,
			Components: &[]discordgo.MessageComponent{},
		})
		if err != nil {
			return err
		}
		keys := make([]string, 0, len(t.handlers))
		for key := range t.handlers {
			keys = append(keys, key)
		}
		t.ch.Unregister(keys...)
		return nil
	}

	return unreg, nil
}

func getCustomId(component discordgo.MessageComponent) string {
	val, _ := util.GetFieldValue(component, "CustomID")
	return val
}

func removeComponentRecursive(components []discordgo.MessageComponent, customKey string) []discordgo.MessageComponent {
	newComponents := make([]discordgo.MessageComponent, 0, len(components))
	for _, comp := range components {
		if comp.Type() == discordgo.ActionsRowComponent {
			ar := comp.(discordgo.ActionsRow)
			newChildren := removeComponentRecursive(ar.Components, customKey)
			if len(newChildren) > 0 {
				newComponents = append(newComponents, discordgo.ActionsRow{
					Components: newChildren,
				})
			}
			continue
		}
		if getCustomId(comp) != customKey {
			newComponents = append(newComponents, comp)
		}
	}

	return newComponents
}
