[VERSION]

- The `User` method of the `Context` interface is now also available for the 
  `ComponentContext` interface.
- You can now set a `Condition` on the `ComponentBuilder` which will be
  checked before the registered component handlers are executed.  
  *See [component example](examples/components/commands/test.go) for further details.*

## Update

```
go get -v -u github.com/zekrotja/ken@[VERSION]
```
