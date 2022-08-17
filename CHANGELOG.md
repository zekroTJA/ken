[VERSION]

## Changes

- Fix a bug where the wrapped handler would not be registered in the build wrapper.
- The `ComponentHandlerFunc` now returns a boolean to indicate success of execution.
- The `Ken` instance now has a method `Session` which returns the wrapped Discordgo `Session` 
  instance.
- The `Build` method of the `ComponentBuilder` now also returns a function to remove components
  from the message as well as unregister the handlers.
- The `Unregister` method of the `ComponentHandler` now can take one or more `customId`s to 
  be unregistered at once.

## Update

```
go get -v -u github.com/zekrotja/ken@[VERSION]
```