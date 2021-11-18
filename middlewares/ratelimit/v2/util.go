package ratelimit

import "github.com/zekrotja/ken"

func Skip(ctx *ken.Ctx) {
	ctx.Set(skipKey, true)
}
