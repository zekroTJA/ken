[VERSION]

**Warning:** This update of ken introduces a lot of API changes. Please update carefully.

## Modals [#10]

You can now open modals in the component handler function. The `ComponentContext` now has a
method `OpenModal` which can be used to open a modal in Discord on interaction with the
message component.

```go
// ComponentContext gives access to the underlying
// MessageComponentInteractionData and gives the
// ability to open a Modal afterwards.
type ComponentContext interface {
    ContextResponder

    // GetData returns the underlying
    // MessageComponentInteractionData.
    GetData() discordgo.MessageComponentInteractionData

    // OpenModal opens a new modal with the given
    // title, content and components built with the
    // passed build function. A channel is returned
    // which will receive a ModalContext when the user
    // has interacted with the modal.
    OpenModal(
        title string,
        content string,
        build func(b ComponentAssembler),
    ) (<-chan ModalContext, error)
}
```

Please take a look at the [**modals example**](examples/modals) to see further details on
how to use modals with ken.

## Breaking API Changes

A lot of breaking changes have been introduced to use more interfaces instead of struct
instances which allows better testability using mocks.

The `Run` method of the `Command` interface now is getting passed a `Context` interface instead
of a reference to an instance of `Ctx`. This also means, if you are directly accessing `Session`
or `Event` for example from the `Ctx` instance, you need to change it to accessing these via the
available getter methods (`GetSession` or `GetEvent` for example).

The `SubCommandHandler` now also passes an interface `SubCommandContext` to
the `Run` instead of a reference to an instance of `SubCommandCtx`.

The access to `CtxResponder`, `SubCommandCtx`, `ComponentCtx` and `ModalCtx` are now private
for a cleaner API.

## Update

```
go get -v -u github.com/zekrotja/ken@[VERSION]
```
