package policy

import (
	"context"
	"fmt"
	"github.com/adverax/log"
)

// WithRecovery is a policy that recovers from panics.
type WithRecovery struct {
	policy Policy
	logger log.Logger
}

func (that *WithRecovery) Execute(ctx context.Context, action Action) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered: %v", r)
			if that.logger != nil {
				that.logger.Fatal(ctx, err.Error())
			}
		}
	}()

	return that.policy.Execute(ctx, action)
}

func NewWithRecovery(policy Policy, logger log.Logger) *WithRecovery {
	if policy == nil {
		policy = dummyPolicy
	}

	return &WithRecovery{
		policy: policy,
		logger: logger,
	}
}
