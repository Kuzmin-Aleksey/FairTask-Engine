package contextx

import (
	"context"
)

type TraceId string

func (id TraceId) String() string {
	return string(id)
}

type contextKeyTraceId struct{}

func WithTraceId(ctx context.Context, tid TraceId) context.Context {
	return context.WithValue(ctx, contextKeyTraceId{}, tid)
}

func GetTraceId(ctx context.Context) TraceId {
	v, _ := ctx.Value(contextKeyTraceId{}).(TraceId)
	return v
}
