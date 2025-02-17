package policy

import (
	"context"
	"github.com/adverax/log"
)

type Loggable interface {
	Named
	Log(ctx context.Context, logger log.Logger)
}

// PolicyWithLogging is a policy that logs the action before executing it.
type PolicyWithLogging struct {
	Policy
	logger log.Logger
}

func NewPolicyWithLogging(policy Policy, logger log.Logger) *PolicyWithLogging {
	return &PolicyWithLogging{
		Policy: policy,
		logger: logger,
	}
}

func (that *PolicyWithLogging) Execute(ctx context.Context, action Action) error {
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
