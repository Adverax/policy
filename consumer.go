package policy

import (
	"context"
)

type Handler[T any] interface {
	Handle(ctx context.Context, entity T) error
}

type HandlerFunc[T any] func(ctx context.Context, entity T) error

func (fn HandlerFunc[T]) Handle(ctx context.Context, entity T) error {
	return fn(ctx, entity)
}

type Consumer[T any] interface {
	Consume(ctx context.Context, entity T)
}

type ConsumerFunc[T any] func(context.Context, T)

func (fn ConsumerFunc[T]) Consume(ctx context.Context, entity T) {
	fn(ctx, entity)
}

type BaseConsumer[T any] struct {
	executor Executor
	handler  Handler[T]
}

func NewConsumer[T any](executor Executor, handler Handler[T]) *BaseConsumer[T] {
	return &BaseConsumer[T]{
		executor: executor,
		handler:  handler,
	}
}

func (that *BaseConsumer[T]) Consume(ctx context.Context, entity T) {
	that.executor.Execute(
		ctx,
		ActionFunc(func(ctx context.Context) error {
			return that.handler.Handle(ctx, entity)
		}),
	)
}
