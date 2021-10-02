package ken

type MiddlewareBefore interface {
	Before(ctx *Ctx) (next bool, err error)
}

type MiddlewareAfter interface {
	After(ctx *Ctx) (err error)
}

type Middleware interface {
	MiddlewareBefore
	MiddlewareAfter
}
