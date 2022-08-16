package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

type TestCommand struct{}

var (
	_ ken.SlashCommand = (*TestCommand)(nil)
	_ ken.DmCapable    = (*TestCommand)(nil)
)

func (c *TestCommand) Name() string {
	return "test"
}

func (c *TestCommand) Description() string {
	return "Basic Test Command"
}

func (c *TestCommand) Version() string {
	return "1.0.0"
}

func (c *TestCommand) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *TestCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionBoolean,
			Name:        "clear-full-action-row",
			Description: "Clear the full action row instead of single button on click.",
		},
	}
}

func (c *TestCommand) IsDmCapable() bool {
	return true
}

func (c *TestCommand) Run(ctx *ken.Ctx) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	var clearAll bool
	if v, ok := ctx.Options().GetByNameOptional("clear-full-action-row"); ok {
		clearAll = v.BoolValue()
	}

	fum := ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: "How are you?",
	})
	if fum.HasError() {
		return fum.Error
	}

	err = fum.AddComponents().
		AddActionsRow(func(b ken.ComponentAssembler) {
			b.Add(discordgo.Button{
				CustomID: "button-1",
				Label:    "Absolutely fantastic!",
			}, func(ctx ken.ComponentContext) {
				ctx.RespondEmbed(&discordgo.MessageEmbed{
					Description: fmt.Sprintf("Responded to %s", ctx.GetData().CustomID),
				})
			}, !clearAll)
			b.Add(discordgo.Button{
				CustomID: "button-2",
				Style:    discordgo.DangerButton,
				Label:    "Not so well",
			}, func(ctx ken.ComponentContext) {
				ctx.RespondEmbed(&discordgo.MessageEmbed{
					Description: fmt.Sprintf("Responded to %s", ctx.GetData().CustomID),
				})
			}, !clearAll)
		}, clearAll).
		Build()

	return err
}
