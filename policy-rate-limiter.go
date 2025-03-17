package policy

import (
	"context"
	"time"
)

type RateLimiter interface {
	Wait(ctx context.Context) error
}

type rateLimiter struct {
	rate   int // requests per second
	tokens chan struct{}
}

func NewRateLimiter(rate, burst int) RateLimiter {
	limiter := &rateLimiter{
		rate:   rate,
		tokens: make(chan struct{}, burst),
	}

	go func() {
		ticker := time.NewTicker(time.Second / time.Duration(rate))
		defer ticker.Stop()
		for range ticker.C {
			select {
			case limiter.tokens <- struct{}{}:
			default: // Do not add more tokens if the channel is full.
			}
		}
	}()

	return limiter
}

func (that *rateLimiter) Wait(ctx context.Context) error {
	select {
	case <-that.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// WithRateLimiter is a policy that ensures that rate limit is not exceeded.
type WithRateLimiter struct {
	policy  Policy
	limiter RateLimiter
}

func NewWithRateLimiter(policy Policy, limiter RateLimiter) *WithRateLimiter {
	return &WithRateLimiter{
		policy:  policy,
		limiter: limiter,
	}
}

func (that *WithRateLimiter) Execute(ctx context.Context, action Action) error {
	err := that.limiter.Wait(ctx)
	if err != nil {
		return err
	}

	return that.policy.Execute(ctx, action)
}
