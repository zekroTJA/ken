package ken

// MessageCommand defines a callable message command.
type MessageCommand interface {
	Command

	// Run is called on command invokation getting
	// passed the invocation context.
	//
	// When something goes wrong during command
	// execution, you can return an error which is
	// then handled by Ken's OnCommandError handler.
	Run(ctx *Ctx) (err error)
}
