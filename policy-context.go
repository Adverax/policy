package policy

import (
	"context"
)

// PolicyWithContextValue is a policy that sets a value in the context before executing the action.
type PolicyWithContextValue struct {
	Policy
	key   interface{}
	value interface{}
}

func (that *PolicyWithContextValue) Execute(ctx context.Context, action Action) error {
	ctx = context.WithValue(ctx, that.key, that.value)
	return that.Policy.Execute(ctx, action)
}

func NewPolicyWithContextValue(
	policy Policy,
	key interface{},
	value interface{},
) *PolicyWithContextValue {
	if policy == nil {
		policy = dummyPolicy
	}

	return &PolicyWithContextValue{
		Policy: policy,
		key:    key,
		value:  value,
	}
}
