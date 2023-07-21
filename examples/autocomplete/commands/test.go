package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

var programmingLanguages = [][]string{
	{"Go", "go"},
	{"Rust", "rust"},
	{"TypeScript", "typescript"},
	{"JavaScript", "javascript"},
	{"Java", "java"},
	{"C++", "cpp"},
	{"C", "c"},
	{"Ruby", "ruby"},
	{"Dart", "dart"},
	{"C#", "csharp"},
}

type TestCommand struct{}

var (
	_ ken.SlashCommand        = (*TestCommand)(nil)
	_ ken.DmCapable           = (*TestCommand)(nil)
	_ ken.AutocompleteCommand = (*TestCommand)(nil)
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
			Type:         discordgo.ApplicationCommandOptionString,
			Name:         "language",
			Required:     true,
			Description:  "Choose a programming language.",
			Autocomplete: true,
		},
	}
}

func (c *TestCommand) Guild() string {
	return "526196711962705925"
}

func (c *TestCommand) IsDmCapable() bool {
	return true
}

func (c *TestCommand) Autocomplete(ctx *ken.AutocompleteContext) ([]*discordgo.ApplicationCommandOptionChoice, error) {
	input, ok := ctx.GetInput("language")

	if !ok {
		return nil, nil
	}

	choises := make([]*discordgo.ApplicationCommandOptionChoice, 0, len(programmingLanguages))
	input = strings.ToLower(input)

	for _, lang := range programmingLanguages {
		if strings.HasPrefix(lang[1], input) {
			choises = append(choises, &discordgo.ApplicationCommandOptionChoice{
				Name:  lang[0],
				Value: lang[0],
			})
		}
	}

	return choises, nil
}

func (c *TestCommand) Run(ctx ken.Context) (err error) {
	lang := ctx.Options().GetByName("language").StringValue()

	return ctx.RespondMessage(
		fmt.Sprintf("%s is an awesome language!", lang))
}
