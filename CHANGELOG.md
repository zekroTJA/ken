[VERSION]

## Sub Command Support [#14]

Finally, you can now register and use sub command groups!

Simply add a sub command option to your command in the `Options` implementation of your command.

*Example:*
```go
func (c *SubsCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
			Name:        "group",
			Description: "Some sub command gorup",
			Options: []*discordgo.ApplicationCommandOption{
				// ...
			},
		},
	}
}
```

Now, you can add a [`SubCommandGroup`](https://pkg.go.dev/github.com/zekrotja/ken#SubCommandGroup) handler to your `HandleSubCommands` method to register sub commands in the group.

*Example:*
```go
func (c *SubsCommand) Run(ctx ken.Context) (err error) {
	err = ctx.HandleSubCommands(
		ken.SubCommandGroup{"group", []ken.CommandHandler{
			ken.SubCommandHandler{"one", c.one},
			ken.SubCommandHandler{"two", c.two},
		}},
	)
	return
}
```

> The full example can be found in [examples/subcommandgroups](https://github.com/zekroTJA/ken/tree/master/examples/subcommandgroups).

## Minor Changes

- The [`AutocompleteContext`](https://pkg.go.dev/github.com/zekrotja/ken#AutocompleteContext) now implements [`ObjectMap`](https://pkg.go.dev/github.com/zekrotja/ken#ObjectMap) so you can pass dependencies down the line just like with command contexts.

- [`AutocompleteContext`](https://pkg.go.dev/github.com/zekrotja/ken#AutocompleteContext) now has a new helper function [`GetInputAny`](https://pkg.go.dev/github.com/zekrotja/ken#AutocompleteContext.GetInputAny) which takes multiple option names and returns the first match found in the event data.

- [`FoolowUpMessage`](https://pkg.go.dev/github.com/zekrotja/ken#FollowUpMessage) has now been added accordingly to [`RespondMessage`](https://pkg.go.dev/github.com/zekrotja/ken@master#Ctx.RespondMessage), which was added in `v0.19.0`.

## Update

```
go get -v -u github.com/zekrotja/ken@[VERSION]
```
