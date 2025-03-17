package policy

import (
	"context"
)

type Tracer interface {
	NewTrace(ctx context.Context, info string) context.Context
}

// WithTraceId is a policy that adds trace id to the context.
type WithTraceId struct {
	policy Policy
	tracer Tracer
	info   string
}

func (that *WithTraceId) Execute(ctx context.Context, action Action) error {
	ctx = that.tracer.NewTrace(ctx, that.info)
	return that.policy.Execute(ctx, action)
}

func NewWithTraceId(
	policy Policy,
	tracer Tracer,
	info string,
) *WithTraceId {
	if policy == nil {
		policy = dummyPolicy
	}

	return &WithTraceId{
		policy: policy,
		tracer: tracer,
		info:   info,
	}
}
