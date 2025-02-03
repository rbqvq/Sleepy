package context

import "context"

type causeContext struct {
	context.Context
}

func (ctx *causeContext) Err() error {
	return context.Cause(ctx.Context)
}

func WithCancelCause(parent context.Context) (ctx context.Context, cancel context.CancelCauseFunc) {
	ctx, cancel = context.WithCancelCause(parent)
	return &causeContext{ctx}, cancel
}
