package ken

import (
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken/state"
)

type Ken struct {
	s   *discordgo.Session
	opt *Options

	cmdsLock sync.RWMutex
	cmds     map[string]Command
	ctxPool  sync.Pool
}

type Options struct {
	State   state.State
	OnError func(area string, err error, args ...interface{})
}

var defaultOptions = Options{
	State: state.NewInternal(),
	OnError: func(ctx string, err error, args ...interface{}) {
		log.Printf("[KEN] {%s} - %s\n", ctx, err.Error())
	},
}

func New(s *discordgo.Session, options ...Options) (k *Ken) {
	k = &Ken{
		s:    s,
		cmds: make(map[string]Command),
		ctxPool: sync.Pool{
			New: func() interface{} {
				return newCtx()
			},
		},
	}

	k.opt = &defaultOptions
	if len(options) > 0 {
		o := options[0]
		if o.State != nil {
			k.opt.State = o.State
		}
		if o.OnError != nil {
			k.opt.OnError = o.OnError
		}
	}

	k.s.AddHandler(k.onReady)
	k.s.AddHandler(k.onInteractionCreate)

	return
}

func (k *Ken) RegisterCommands(cmds ...Command) (err error) {
	k.cmdsLock.Lock()
	defer k.cmdsLock.Unlock()

	for _, c := range cmds {
		err = k.registerCommand(c)

		if err != nil {
			return
		}
	}

	return
}

func (k *Ken) RegisterMiddlewares(mws ...Middleware) {
	for _, mw := range mws {
		k.registerMiddleware(mw)
	}
}

func (k *Ken) registerCommand(cmd Command) (err error) {
	if cmd.Name() == "" {
		err = ErrEmptyCommandName
		return
	}
	if _, ok := k.cmds[cmd.Name()]; ok {
		err = ErrCommandAlreadyRegistered
		return
	}

	k.cmds[cmd.Name()] = cmd

	return
}

func (k *Ken) registerMiddleware(mw Middleware) {

}

func (k *Ken) onReady(s *discordgo.Session, e *discordgo.Ready) {
	k.cmdsLock.RLock()
	defer k.cmdsLock.RUnlock()

	for _, cmd := range k.cmds {
		if _, err := s.ApplicationCommandCreate(e.User.ID, "", toApplicationCommand(cmd)); err != nil {
			k.opt.OnError("command registration", err)
		}
	}
}

func (k *Ken) onInteractionCreate(s *discordgo.Session, e *discordgo.InteractionCreate) {
	k.cmdsLock.RLock()
	cmd := k.cmds[e.ApplicationCommandData().Name]
	k.cmdsLock.RUnlock()

	if cmd == nil {
		return
	}

	ctx := k.ctxPool.Get().(*Ctx)
	defer k.ctxPool.Put(ctx)
	ctx.st = k.opt.State
	ctx.Session = s
	ctx.Event = e
	ctx.Command = cmd

	if err := cmd.Run(ctx); err != nil {
		k.opt.OnError("command error", err, ctx)
	}
}
