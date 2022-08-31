package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

type ModalCommand struct{}

var (
	_ ken.SlashCommand = (*ModalCommand)(nil)
	_ ken.DmCapable    = (*ModalCommand)(nil)
)

func (c *ModalCommand) Name() string {
	return "modal"
}

func (c *ModalCommand) Description() string {
	return "Modal Test Command"
}

func (c *ModalCommand) Version() string {
	return "1.0.0"
}

func (c *ModalCommand) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *ModalCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}

func (c *ModalCommand) IsDmCapable() bool {
	return false
}

func (c *ModalCommand) Run(ctx ken.Context) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	fum := ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: "How are you?",
	})
	if fum.HasError() {
		return fum.Error
	}

	_, err = fum.AddComponents().
		AddActionsRow(func(b ken.ComponentAssembler) {
			b.Add(discordgo.Button{
				CustomID: "open-modal",
				Label:    "Write it!",
				Style:    discordgo.PrimaryButton,
			}, func(ctx ken.ComponentContext) bool {
				cCtx, err := ctx.OpenModal("Hello world", "Lorem ipsum ...", func(b ken.ComponentAssembler) {
					b.AddActionsRow(func(b ken.ComponentAssembler) {
						b.Add(discordgo.TextInput{
							CustomID:  "text-input",
							Label:     "How are you?",
							Style:     discordgo.TextInputShort,
							Required:  true,
							MaxLength: 1000,
						}, nil)
					})
				})

				if err != nil {
					fmt.Println("Error:", err)
					return false
				}

				embCtx := <-cCtx

				resp := embCtx.GetComponentByID("text-input").GetValue()
				embCtx.RespondEmbed(&discordgo.MessageEmbed{
					Description: fmt.Sprintf(`"%s" - ok, thats cool`, resp),
				})
				return true
			})
		}, true).
		Build()

	return err
}
