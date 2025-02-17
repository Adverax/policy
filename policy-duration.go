package policy

import (
	"context"
	"github.com/adverax/log"
	"time"
)

type Named interface {
	Name() string
}

// WithDuration is a policy that executes the action and logs the duration.
type WithDuration struct {
	Policy
	logger   log.Logger
	limit    time.Duration
	wantInfo bool
}

func NewWithDuration(
	policy Policy,
	logger log.Logger,
	limit time.Duration,
	wantInfo bool,
) *WithDuration {
	return &WithDuration{
		Policy:   policy,
		logger:   logger,
		limit:    limit,
		wantInfo: wantInfo,
	}
}

func (that *WithDuration) Execute(ctx context.Context, action Action) error {
	if a, ok := action.(Named); ok {
		started := time.Now()

		err := that.Policy.Execute(ctx, action)
		if err != nil {
			return err
		}

		duration := time.Since(started)
		if that.limit == 0 || duration < that.limit {
			if that.wantInfo {
				that.logger.WithField(log.FieldKeyDuration, duration).
					Infof(ctx, "finished executing action %s", a.Name())
			}
		} else {
			that.logger.WithField(log.FieldKeyDuration, duration).
				Warningf(ctx, "finished executing action %s", a.Name())
		}

		return nil
	}

	err := that.Policy.Execute(ctx, action)
	if err != nil {
		return err
	}

	return nil
}
