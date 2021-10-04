<img width="100%" src=".github/media/ken.png" />

# ken &nbsp; [![](https://img.shields.io/badge/docs-pkg.do.dev-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/zekrotja/ken?tab=doc) [![Go Report Card](https://goreportcard.com/badge/github.com/zekrotja/ken)](https://goreportcard.com/report/github.com/zekrotja/ken)

剣 *(`ken` - japanese for `Sword`)* - A cutting edge *(haha)*, prototype, object-oriented and highly modular slash command handler for discordgo.

This work-in-progress slash command handler is designed to be used in future versions of [shinpuru](https://github.com/zekroTJA/shinpuru).

If you are interested how to use this package, see the [basic example](examples/basic). In [this]([basic example](examples/middlewares)) example, you can see how to use middlewares.

## Example Usage

```go
package main

// imports ...

type TestCommand struct{}

var _ ken.Command = (*TestCommand)(nil)

func (c *TestCommand) Name() string {
	return "ping"
}

func (c *TestCommand) Description() string {
	return "Basic Ping Command"
}

func (c *TestCommand) Version() string {
	return "1.0.0"
}

func (c *TestCommand) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *TestCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}

func (c *TestCommand) Run(ctx *ken.Ctx) (err error) {
	err = ctx.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong!",
		},
	})
	return
}

func main() {
	token := os.Getenv("TOKEN")

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	k := ken.New(session)
	k.RegisterCommands(
		new(commands.TestCommand),
	)

	defer k.Unregister()

	err = session.Open()
	if err != nil {
		panic(err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
```

You can also find a "real world" implementation in my Discord bot [shinpuru](https://github.com/zekrotja/shinpuru), where ken is used as main slash command framework.

---

© 2021 Ringo Hoffmann (zekro Development).  
Covered by the MIT Licence.
