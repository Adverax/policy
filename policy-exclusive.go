package policy

import (
	"context"
	"sync"
)

// WithExclusiveExecution is a policy that ensures that only one action is executed at a time.
type WithExclusiveExecution struct {
	Policy Policy
	sync.Mutex
}

func NewWithExclusiveExecution(policy Policy) *WithExclusiveExecution {
	return &WithExclusiveExecution{Policy: policy}
}

func (that *WithExclusiveExecution) Execute(ctx context.Context, action Action) error {
	that.Lock()
	defer that.Unlock()

	return that.Policy.Execute(ctx, action)
}
