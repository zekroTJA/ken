package cmdhelp

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

// HelpProvider defines a command which provides
// help content.
type HelpProvider interface {
	Help(ctx *ken.SubCommandCtx) (emb *discordgo.MessageEmbed, err error)
}
