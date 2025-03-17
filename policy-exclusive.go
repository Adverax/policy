package policy

import (
	"context"
	"sync"
)

// WithExclusive is a policy that ensures that only one action is executed at a time.
type WithExclusive struct {
	Policy Policy
	sync.Mutex
}

func NewWithExclusive(policy Policy) *WithExclusive {
	return &WithExclusive{Policy: policy}
}

func (that *WithExclusive) Execute(ctx context.Context, action Action) error {
	that.Lock()
	defer that.Unlock()

	return that.Policy.Execute(ctx, action)
}
