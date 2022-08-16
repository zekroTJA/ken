[VERSION]

## Changes

### Message Component Implementation [#8]

You can now attach message components to messages and register interaction handlers
using the [ComponentHandler](https://pkg.go.dev/github.com/zekrotja/ken#ComponentHandler).

It provides a low level implementation to add components to messages and handlers directly
to the handler registry as well as a more sophisticated
[ComponentBuilder](https://pkg.go.dev/github.com/zekrotja/ken#ComponentBuilder) to simplify
the attachment and handling of message components.

#### Example

```go
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
```

You can find the complete example in [examples/components](examples/components).

### Guild Scopes [#11]

You can now scope commands to specific guilds by implementing the
[GuildScopedCommand](https://pkg.go.dev/github.com/zekrotja/ken#GuildScopedCommand) interface
in you command. Then, the command will only be registered for the guild returned by
the `Guild` method of the command.


#### Example

```go
type Guild1Command struct{}

var (
    _ ken.SlashCommand       = (*Guild1Command)(nil)
    _ ken.GuildScopedCommand = (*Guild1Command)(nil)
)

func (c *Guild1Command) Name() string {
    return "guild1"
}

func (c *Guild1Command) Guild() string {
    return "362162947738566657"
}
```

You can find the complete example in [examples/guildscoped](examples/guildscoped).

## Update

```
go get -v -u github.com/zekrotja/ken@[VERSION]
```