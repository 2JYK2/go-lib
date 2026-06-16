package trace

import (
	"context"
)

const reqIpKey = "ReqIp"

// WithReqIP set reqip context
func WithReqIP(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, reqIpKey, traceID)
}

// GetReqIP context get reqip
func GetReqIP(ctx context.Context) string {
	if traceID, ok := ctx.Value(reqIpKey).(string); ok {
		return traceID
	}
	return ""
}
