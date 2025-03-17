package policy

import (
	"context"
)

type Control interface {
	Enter()
	Leave()
}

// WithAsync is a policy that executes the action in a separate goroutine.
type WithAsync struct {
	Executor
	control Control
}

func NewWithAsync(
	executor Executor,
	control Control,
) *WithAsync {
	if control == nil {
		control = dummyWG
	}
	if executor == nil {
		executor = NewDefaultExecutor()
	}
	return &WithAsync{
		control:  control,
		Executor: executor,
	}
}

func (that *WithAsync) Execute(ctx context.Context, action Action) error {
	that.control.Enter()
	go func() {
		defer that.control.Leave()

		that.Executor.Execute(ctx, action)
	}()

	return nil
}

// WithPool is a policy that executes the action in a separate goroutine.
type WithPool struct {
	Executor
	pool    chan struct{}
	control Control
}

func NewWithPool(
	executor Executor,
	control Control,
	size int,
) *WithPool {
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

	return &WithPool{
		pool:     pool,
		control:  control,
		Executor: executor,
	}
}

func (that *WithPool) Execute(ctx context.Context, action Action) error {
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
		return NewWithAsync(executor, control)
	}

	return NewWithPool(executor, control, size)
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
