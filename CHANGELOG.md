[VERSION]

> **Warning**  
> This update contains breaking changes! 

In order to make the code of ken more type-safe, I am now using [safepool](https://github.com/zekrotja/safepool) 
instead of `sync.Pool` which uses generic type parameters for ensuring type safety. Therefore, the minimum
required module version has been bumped to `go1.18`. You might need to upgrade your project accordingly to
be able to use ken.

## Autocomplete Support [#18]

[Command option autocomplete](https://discord.com/developers/docs/interactions/application-commands#autocomplete)
support has now been added to ken.

Simply enable autocomplete on your command option by setting the `Autocomplete` property to `true`.

*Example:*
```go
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
```

Now, you simply need to implement the [`AutocompleteCommand`](https://pkg.go.dev/github.com/zekrotja/ken#Command) 
on your command.

*Example:*
```go
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
```

> The full example can be found in [examples/autocomplete](https://github.com/zekroTJA/ken/tree/master/examples/autocomplete).

To properly handle errors occuring during autocomplete handling, a new command handler hook `OnEventError` has been added to the [`Options`](https://pkg.go.dev/github.com/zekrotja/ken#Options). It will be called every time a non-command related user event error occurs.

## `RespondMessage`

A new respond method has been added to the `ContentResponder` called [`RespondMessage`](https://pkg.go.dev/github.com/zekrotja/ken#Ctx.RespondMessage). It simply takes a message as parameter and responds with a simple message content containing the passed message.

## Update

```
go get -v -u github.com/zekrotja/ken@[VERSION]
```
