package policy

import (
	"context"
	"sync"
)

// PolicyWithExclusiveExecution is a policy that ensures that only one action is executed at a time.
type PolicyWithExclusiveExecution struct {
	Policy Policy
	sync.Mutex
}

func NewPolicyWithExclusiveExecution(policy Policy) *PolicyWithExclusiveExecution {
	return &PolicyWithExclusiveExecution{Policy: policy}
}

func (that *PolicyWithExclusiveExecution) Execute(ctx context.Context, action Action) error {
	that.Lock()
	defer that.Unlock()

	return that.Policy.Execute(ctx, action)
}
