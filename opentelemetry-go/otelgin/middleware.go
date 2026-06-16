package otelgin

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const traceIDKey = "traceID"

func ResponseHeaderWithTraceInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// add traceID to response header
		// when context does not initialize span, tranceID is 00000000000000000000000000000000
		span := trace.SpanFromContext(ctx.Request.Context())
		traceID := span.SpanContext().TraceID().String()
		ctx.Header(traceIDKey, traceID)
		// add attrbute
		span.SetAttributes(
			attribute.String("server.ip", LocalIp),
		)
		// add traceParent to response header
		pp := otel.GetTextMapPropagator()
		carrier := propagation.MapCarrier{}
		pp.Inject(ctx.Request.Context(), carrier)

		for k, v := range carrier {
			ctx.Header(k, v)
		}
		ctx.Next()
		if len(ctx.Errors) != 0 {
			span.SetStatus(codes.Error, ctx.Errors.String())
		}
	}
}
