// Package ken provides an object-oriented and
// highly modular slash command handler for
// discordgo.
package ken

import (
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken/state"
)

// Options holds configurations for Ken.
type Options struct {
	// State specifies the state manager to be used.
	// When not specified, the default discordgo state
	// manager is used.
	State state.State
	// DependencyProvider can be used to inject dependencies
	// to be used in a commands or middlewares Ctx by
	// a string key.
	DependencyProvider ObjectProvider
	// OnSystemError is called when a recoverable
	// system error occurs inside Ken's lifecycle.
	OnSystemError func(context string, err error, args ...interface{})
	// OnCommandError is called when an error occurs
	// during middleware or command execution.
	OnCommandError func(err error, ctx *Ctx)
}

// Ken is the handler to register, manage and
// life-cycle commands as well as middlewares.
type Ken struct {
	s   *discordgo.Session
	opt *Options

	cmdsLock sync.RWMutex
	cmds     map[string]Command
	cmdIDs   []string
	ctxPool  sync.Pool

	mwBefore []MiddlewareBefore
	mwAfter  []MiddlewareAfter
}

var defaultOptions = Options{
	State: state.NewInternal(),
	OnSystemError: func(ctx string, err error, args ...interface{}) {
		log.Printf("[KEN] {%s} - %s\n", ctx, err.Error())
	},
	OnCommandError: func(err error, ctx *Ctx) {
		log.Printf("[KEN] {command error} - %s : %s\n", ctx.Command.Name(), err.Error())
	},
}

// New initializes a new instance of Ken with
// the passed discordgo Session s and optional
// Options.
//
// If no options are passed, default parameters
// will be applied.
func New(s *discordgo.Session, options ...Options) (k *Ken) {
	k = &Ken{
		s:      s,
		cmds:   make(map[string]Command),
		cmdIDs: make([]string, 0),
		ctxPool: sync.Pool{
			New: func() interface{} {
				return newCtx()
			},
		},
		mwBefore: make([]MiddlewareBefore, 0),
		mwAfter:  make([]MiddlewareAfter, 0),
	}

	k.opt = &defaultOptions
	if len(options) > 0 {
		o := options[0]
		if o.State != nil {
			k.opt.State = o.State
		}
		if o.OnSystemError != nil {
			k.opt.OnSystemError = o.OnSystemError
		}
		if o.OnCommandError != nil {
			k.opt.OnCommandError = o.OnCommandError
		}
	}

	k.s.AddHandler(k.onReady)
	k.s.AddHandler(k.onInteractionCreate)

	return
}

// RegisterCommands registers the passed commands to
// the command register.
//
// Keep in mind that commands are registered by Name,
// so there can be only one single command per name.
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

// RegisterMiddlewares allows to register passed
// commands to the middleware callstack.
//
// Therefore, you can register MiddlewareBefore,
// which is called before the command handler is
// executed, or MiddlewareAfter, which is called
// directly after the command handler has been
// called. Of course, you can also implement both
// interfaces in the same middleware.
//
// The middleware call order is determined by the
// order of middleware registraion in each area
// ('before' or 'after').
func (k *Ken) RegisterMiddlewares(mws ...interface{}) (err error) {
	for _, mw := range mws {
		if err = k.registerMiddleware(mw); err != nil {
			break
		}
	}
	return
}

// Unregister should be called to cleanly unregister
// all registered slash commands from the discord
// backend.
func (k *Ken) Unregister() (err error) {
	self, err := k.opt.State.SelfUser(k.s)
	if err != nil {
		return
	}
	for _, id := range k.cmdIDs {
		if err = k.s.ApplicationCommandDelete(self.ID, "", id); err != nil {
			k.opt.OnSystemError("command unregister", err)
		}
	}
	return
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

func (k *Ken) registerMiddleware(mw interface{}) (err error) {
	var (
		okBefore, okAfter bool
		mwBefore          MiddlewareBefore
		mwAfter           MiddlewareAfter
	)
	if mwBefore, okBefore = mw.(MiddlewareBefore); okBefore {
		k.mwBefore = append(k.mwBefore, mwBefore)
	}
	if mwAfter, okAfter = mw.(MiddlewareAfter); okAfter {
		k.mwAfter = append(k.mwAfter, mwAfter)
	}
	if !okBefore && !okAfter {
		err = ErrInvalidMiddleware
	}
	return
}

func (k *Ken) onReady(s *discordgo.Session, e *discordgo.Ready) {
	k.cmdsLock.RLock()
	defer k.cmdsLock.RUnlock()

	for _, cmd := range k.cmds {
		ccmd, err := s.ApplicationCommandCreate(e.User.ID, "", toApplicationCommand(cmd))
		if err != nil {
			k.opt.OnSystemError("command registration", err)
		} else {
			k.cmdIDs = append(k.cmdIDs, ccmd.ID)
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
	ctx.Purge()
	ctx.st = k.opt.State
	ctx.dp = k.opt.DependencyProvider
	ctx.Session = s
	ctx.Event = e
	ctx.Command = cmd

	for _, mw := range k.mwBefore {
		next, err := mw.Before(ctx)
		if err != nil {
			k.opt.OnCommandError(err, ctx)
		}
		if !next {
			return
		}
	}

	err := cmd.Run(ctx)
	if err != nil {
		k.opt.OnCommandError(err, ctx)
	}

	for _, mw := range k.mwAfter {
		err := mw.After(ctx, err)
		if err != nil {
			k.opt.OnCommandError(err, ctx)
		}
	}
}
