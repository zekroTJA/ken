package ken

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken/util"
)

type MessageComponent struct {
	discordgo.MessageComponent
}

func (t MessageComponent) GetValue() string {
	val, _ := util.GetFieldValue(t.MessageComponent, "Value")
	return val
}
