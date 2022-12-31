[VERSION]

> **Warning**  
> This update contains breaking changes!

In order to support the attachment of message components to a flollow-up message on creation [see #13],
the methods `ContextResponder#FollowUp`, `ContextResponder#FollowUpEmbed` and `ContextResponder#FollowUpError`
now return a `*FollowUpMessageBuilder`, which can be used to attach components and handlers before the
follow-up message is sent.

**Example**

```go
b := ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: "Hello",
})

b.AddComponents(func(cb *ken.ComponentBuilder) {
    cb.Add(
        discordgo.Button{
			      CustomID: "button-1",
				    Label:    "Absolutely fantastic!",
	      }, 
        func(ctx ken.ComponentContext) bool {
				    ctx.RespondEmbed(&discordgo.MessageEmbed{
					      Description: fmt.Sprintf("Responded to %s", ctx.GetData().CustomID),
				    })
				    return true
		    }, true,
    )  
})

fum := b.Send()
```

A full example can be found in [examples/components/commands/test.go](examples/components/commands/test.go).

## Update

```
go get -v -u github.com/zekrotja/ken@[VERSION]
```
