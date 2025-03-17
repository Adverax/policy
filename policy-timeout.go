package policy

import (
	"context"
	"time"
)

// WithTimeout is a policy that sets a timeout on the context before executing the action.
type WithTimeout struct {
	policy  Policy
	timeout time.Duration
}

func (that *WithTimeout) Execute(ctx context.Context, action Action) error {
	ctx2, cancel := context.WithTimeout(ctx, that.timeout)
	defer cancel()

	return that.policy.Execute(ctx2, action)
}

func NewWithTimeout(policy Policy, timeout time.Duration) *WithTimeout {
	if policy == nil {
		policy = dummyPolicy
	}

	return &WithTimeout{
		policy:  policy,
		timeout: timeout,
	}
}
