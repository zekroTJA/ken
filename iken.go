package ken

import "github.com/bwmarrin/discordgo"

type IKen interface {
	Components() *ComponentHandler
	GetCommandInfo(keyTransformer ...KeyTransformerFunc) (cis CommandInfoList)
	RegisterCommands(cmds ...Command) (err error)
	RegisterMiddlewares(mws ...interface{}) (err error)
	Session() *discordgo.Session
	Unregister() (err error)
}
