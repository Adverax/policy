package policy

import (
	"context"
)

type ErrorHandler interface {
	HandleError(ctx context.Context, err error)
}

type Executor interface {
	Execute(ctx context.Context, action Action)
}

type BaseExecutor struct {
	policy Policy
	errors ErrorHandler
}

func NewBaseExecutor(
	policy Policy,
	errors ErrorHandler,
) *BaseExecutor {
	if errors == nil {
		errors = defErrorHandler
	}

	return &BaseExecutor{
		policy: policy,
		errors: errors,
	}
}

func (that *BaseExecutor) Execute(
	ctx context.Context,
	action Action,
) {
	err := that.policy.Execute(ctx, action)
	if err != nil {
		that.errors.HandleError(ctx, err)
	}
}

func NewDefaultExecutor() Executor {
	return NewBaseExecutor(
		NewDefault(),
		defErrorHandler,
	)
}

type defaultErrorHandler struct{}

func (that *defaultErrorHandler) HandleError(context.Context, error) {}

var defErrorHandler = &defaultErrorHandler{}
