package otelgin

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

const TracerKey = "otel-go-contrib-tracer"

func Span(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (spanctx context.Context, span trace.Span) {
	value := ctx.Value(TracerKey)
	tracer, ok := value.(trace.Tracer)
	if !ok {
		return ctx, nil
	}

	if c, ok := ctx.(*gin.Context); ok {
		spanctx, span = tracer.Start(c.Request.Context(), spanName, opts...)
		spanctx = context.WithValue(spanctx, TracerKey, tracer)
	} else {
		spanctx, span = tracer.Start(ctx, spanName, opts...)
	}
	return spanctx, span
}
