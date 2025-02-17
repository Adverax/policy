package policy

import (
	"context"
	"github.com/adverax/log"
)

type Loggable interface {
	Named
	Log(ctx context.Context, logger log.Logger)
}

// WithLogging is a policy that logs the action before executing it.
type WithLogging struct {
	Policy
	logger log.Logger
}

func NewWithLogging(policy Policy, logger log.Logger) *WithLogging {
	return &WithLogging{
		Policy: policy,
		logger: logger,
	}
}

func (that *WithLogging) Execute(ctx context.Context, action Action) error {
	if a, ok := action.(Loggable); ok {
		a.Log(ctx, that.logger)

		err := that.Policy.Execute(ctx, action)
		if err != nil {
			return err
		}

		return nil
	}

	return that.Policy.Execute(ctx, action)
}
