package policy

import "context"

// DefaultPolicy is a policy that executes the action as is.
type DefaultPolicy struct{}

func NewDefault() *DefaultPolicy {
	return &DefaultPolicy{}
}

func (c *DefaultPolicy) Execute(ctx context.Context, action Action) error {
	return action.Execute(ctx)
}

var dummyPolicy = NewDefault()
