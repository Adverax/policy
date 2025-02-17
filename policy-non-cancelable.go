package policy

import (
	"context"
)

// PolicyNonCancelable is a policy that makes the action non-cancelable.
type PolicyNonCancelable struct {
	Policy
}

func (that *PolicyNonCancelable) Execute(ctx context.Context, action Action) error {
	ctx = &nonCancelableContext{Context: ctx}
	return that.Policy.Execute(ctx, action)
}

func NewPolicyNonCancelable(
	policy Policy,
) *PolicyNonCancelable {
	if policy == nil {
		policy = dummyPolicy
	}

	return &PolicyNonCancelable{
		Policy: policy,
	}
}

type nonCancelableContext struct {
	context.Context
}

func (that *nonCancelableContext) Done() <-chan struct{} {
	return dummyDone
}

var dummyDone = make(chan struct{})
