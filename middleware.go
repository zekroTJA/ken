package ken

type MiddlewareBefore interface {
	Before(ctx *Ctx) (next bool, err error)
}

type MiddlewareAfter interface {
	After(ctx *Ctx, cmdError error) (err error)
}

type Middleware interface {
	MiddlewareBefore
	MiddlewareAfter
}
