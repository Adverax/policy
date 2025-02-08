package policy

import (
	"context"
)

type Action interface {
	Execute(ctx context.Context) error
}

type ActionFunc func(ctx context.Context) error

func (fn ActionFunc) Execute(ctx context.Context) error {
	return fn(ctx)
}

type Policy interface {
	Execute(ctx context.Context, action Action) error
}

type Handler[T any] interface {
	Handle(ctx context.Context, entity T) error
}

type HandlerFunc[T any] func(ctx context.Context, entity T) error

func (fn HandlerFunc[T]) Handle(ctx context.Context, entity T) error {
	return fn(ctx, entity)
}
