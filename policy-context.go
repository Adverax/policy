package policy

import (
	"context"
)

// WithContextValue is a policy that sets a value in the context before executing the action.
type WithContextValue struct {
	Policy
	key   interface{}
	value interface{}
}

func (that *WithContextValue) Execute(ctx context.Context, action Action) error {
	ctx = context.WithValue(ctx, that.key, that.value)
	return that.Policy.Execute(ctx, action)
}

func NewWithContextValue(
	policy Policy,
	key interface{},
	value interface{},
) *WithContextValue {
	if policy == nil {
		policy = dummyPolicy
	}

	return &WithContextValue{
		Policy: policy,
		key:    key,
		value:  value,
	}
}
