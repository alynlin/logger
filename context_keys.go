package logger

import "context"

type ctxKey string

const (
	traceIDKey ctxKey = "trace_id"
	userIDKey  ctxKey = "user_id"
)

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func FromContext(ctx context.Context) (traceID, userID string) {
	if v := ctx.Value(traceIDKey); v != nil {
		traceID, _ = v.(string)
	}
	if v := ctx.Value(userIDKey); v != nil {
		userID, _ = v.(string)
	}
	return
}
