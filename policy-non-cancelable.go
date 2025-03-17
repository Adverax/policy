package policy

import (
	"context"
)

// NonCancelable is a policy that makes the action non-cancelable.
type NonCancelable struct {
	policy Policy
}

func (that *NonCancelable) Execute(ctx context.Context, action Action) error {
	ctx = &nonCancelableContext{Context: ctx}
	return that.policy.Execute(ctx, action)
}

func NewNonCancelable(
	policy Policy,
) *NonCancelable {
	if policy == nil {
		policy = dummyPolicy
	}

	return &NonCancelable{
		policy: policy,
	}
}

type nonCancelableContext struct {
	context.Context
}

func (that *nonCancelableContext) Done() <-chan struct{} {
	return dummyDone
}

var dummyDone = make(chan struct{})
