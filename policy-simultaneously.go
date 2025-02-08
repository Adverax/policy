package policy

import (
	"context"
)

type Control interface {
	Enter()
	Leave()
}

type PolicyWithAsyncExecution struct {
	Executor
	control Control
}

func newPolicyWithAsyncExecution(
	executor Executor,
	control Control,
) *PolicyWithAsyncExecution {
	if control == nil {
		control = dummyWG
	}
	if executor == nil {
		executor = NewDefaultExecutor()
	}
	return &PolicyWithAsyncExecution{
		control:  control,
		Executor: executor,
	}
}

func (that *PolicyWithAsyncExecution) Execute(ctx context.Context, action Action) error {
	that.control.Enter()
	go func() {
		defer that.control.Leave()

		that.Executor.Execute(ctx, action)
	}()

	return nil
}

type PolicyWithPoolExecution struct {
	Executor
	pool    chan struct{}
	control Control
}

func newPolicyWithPoolExecution(
	executor Executor,
	control Control,
	size int,
) *PolicyWithPoolExecution {
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

	return &PolicyWithPoolExecution{
		pool:     pool,
		control:  control,
		Executor: executor,
	}
}

func (that *PolicyWithPoolExecution) Execute(ctx context.Context, action Action) error {
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

func NewSimultaneouslyPolicy(
	executor Executor,
	control Control,
	size int,
) Policy {
	if size == 0 {
		return newPolicyWithAsyncExecution(executor, control)
	}

	return newPolicyWithPoolExecution(executor, control, size)
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
