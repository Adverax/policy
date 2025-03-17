package policy

import (
	"context"
	"github.com/adverax/log"
)

type Loggable interface {
	LogEnter(ctx context.Context, logger log.Logger)
	LogLeave(ctx context.Context, logger log.Logger)
	LogError(ctx context.Context, logger log.Logger, err error)
}

// WithLogging is a policy that logs the action before executing it.
type WithLogging struct {
	policy Policy
	logger log.Logger
}

func NewWithLogging(policy Policy, logger log.Logger) *WithLogging {
	return &WithLogging{
		policy: policy,
		logger: logger,
	}
}

func (that *WithLogging) Execute(ctx context.Context, action Action) error {
	if a, ok := action.(Loggable); ok {
		a.LogEnter(ctx, that.logger)

		err := that.policy.Execute(ctx, action)
		if err != nil {
			a.LogError(ctx, that.logger, err)
			return err
		}

		a.LogLeave(ctx, that.logger)
		return nil
	}

	return that.policy.Execute(ctx, action)
}
