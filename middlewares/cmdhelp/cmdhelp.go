package cmdhelp

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

// Middleware implements ken.MiddlewareBefore. It checks if a
// command implements HelpProvider on execution and attaches
// a help sub command handler. When the help sub command was
// called, a message is responded containing the embed returned
// by the Help implementation.
type Middleware struct {
	subCommandName string
}

var _ ken.MiddlewareBefore = (*Middleware)(nil)

// New returns a new instance of Middleware.
//
// You can also pass a custom name for the help
// sub command. This defaults to "help" otherwise.
func New(subCommandName ...string) *Middleware {
	scn := "help"
	if len(subCommandName) != 0 {
		scn = subCommandName[0]
	}
	return &Middleware{scn}
}

func (m *Middleware) Before(ctx *ken.Ctx) (next bool, err error) {
	next = true

	cmd, ok := ctx.Command.(HelpProvider)
	if !ok {
		return
	}

	cExecuted := make(chan bool, 1)
	if err = ctx.HandleSubCommands(m.subCmdHandler(cmd, cExecuted)); err != nil {
		next = false
		return
	}

	if len(cExecuted) == 1 && <-cExecuted {
		next = false
	}

	return
}

func (m *Middleware) subCmdHandler(cmd HelpProvider, executed chan<- bool) ken.SubCommandHandler {
	return ken.SubCommandHandler{
		Name: m.subCommandName,
		Run: func(ctx ken.SubCommandContext) (err error) {
			executed <- true
			emb, err := cmd.Help(ctx)
			if err != nil {
				return
			}
			err = ctx.Respond(&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{emb},
				},
			})
			return
		},
	}
}
