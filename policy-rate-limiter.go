package policy

import (
	"context"
)

// WithRateLimiter is a policy that ensures that rate limit is not exceeded.
type WithRateLimiter struct {
	Policy Policy
}

func NewWithRateLimiter(policy Policy) *WithRateLimiter {
	return &WithRateLimiter{Policy: policy}
}

func (that *WithRateLimiter) Execute(ctx context.Context, action Action) error {

	return that.Policy.Execute(ctx, action)
}
