package ken

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken/util"
)

type MessageComponent struct {
	discordgo.MessageComponent
}

func (t MessageComponent) GetValue() string {
	if t.MessageComponent == nil {
		return ""
	}

	val, _ := util.GetFieldValue(t.MessageComponent, "Value")
	return val
}

func (t MessageComponent) IsEmpty() bool {
	return t.MessageComponent == nil
}
