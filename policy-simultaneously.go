package policy

import (
	"context"
)

type Control interface {
	Enter()
	Leave()
}

// WithAsyncExecution is a policy that executes the action in a separate goroutine.
type WithAsyncExecution struct {
	Executor
	control Control
}

func NewWithAsyncExecution(
	executor Executor,
	control Control,
) *WithAsyncExecution {
	if control == nil {
		control = dummyWG
	}
	if executor == nil {
		executor = NewDefaultExecutor()
	}
	return &WithAsyncExecution{
		control:  control,
		Executor: executor,
	}
}

func (that *WithAsyncExecution) Execute(ctx context.Context, action Action) error {
	that.control.Enter()
	go func() {
		defer that.control.Leave()

		that.Executor.Execute(ctx, action)
	}()

	return nil
}

// WithPoolExecution is a policy that executes the action in a separate goroutine.
type WithPoolExecution struct {
	Executor
	pool    chan struct{}
	control Control
}

func NewWithPoolExecution(
	executor Executor,
	control Control,
	size int,
) *WithPoolExecution {
	if control == nil {
		control = dummyWG
	}
	if executor == nil {
		executor = NewDefaultExecutor()
	}

	pool := make(chan struct{}, size)
	for i := 0; i < size; i++ {
		pool <- struct{}{}
	}

	return &WithPoolExecution{
		pool:     pool,
		control:  control,
		Executor: executor,
	}
}

func (that *WithPoolExecution) Execute(ctx context.Context, action Action) error {
	that.control.Enter()

	select {
	case <-that.pool:
	case <-ctx.Done():
		return ctx.Err()
	}

	go func() {
		defer func() {
			that.pool <- struct{}{}
			that.control.Leave()
		}()

		that.Executor.Execute(ctx, action)
	}()

	return nil
}

func NewWithSimultaneously(
	executor Executor,
	control Control,
	size int,
) Policy {
	if size == 0 {
		return NewWithAsyncExecution(executor, control)
	}

	return NewWithPoolExecution(executor, control, size)
}

type dummyWaitGroup struct {
}

func (that *dummyWaitGroup) Enter() {
}

func (that *dummyWaitGroup) Leave() {
}

func (that *dummyWaitGroup) Wait() {
}

func (that *dummyWaitGroup) WaitWithContext(ctx context.Context) error {
	return nil
}

var dummyWG = new(dummyWaitGroup)
