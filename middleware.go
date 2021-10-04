package ken

// MiddlewareBefore specifies a middleware which is
// called before the execution of a command.
type MiddlewareBefore interface {
	// Before is called before a command is executed.
	// It is getting passed the same context as which
	// will be passed to the command. So you are able
	// to attach or alter data of the context.
	//
	// Ctx contains an ObjectMap which can be used to
	// pass data to the command.
	//
	// The method returns a bool which specifies if
	// the subsequent command should be executed. If
	// it is set to false, the execution will be
	// canceled.
	//
	// The error return value is either bubbled up to
	// the OnCommandError, when next is set to false.
	// Otherwise, the error is passed to OnCommandError
	// but the execution continues.
	Before(ctx *Ctx) (next bool, err error)
}

// MiddlewareAfter specifies a middleware which is
// called after the execution of a command.
type MiddlewareAfter interface {
	// After is called after a command has been executed.
	//
	// It is getting passed the Ctx which was also passed
	// to the command Run handler. Also, the method is
	// getting passed potential errors which were returned
	// from the executed command to be custom handled.
	//
	// The error returned is finally passed to the
	// OnCommandError handler.
	After(ctx *Ctx, cmdError error) (err error)
}

// Middleware combines MiddlewareBefore and
// MiddlewareAfter.
type Middleware interface {
	MiddlewareBefore
	MiddlewareAfter
}
