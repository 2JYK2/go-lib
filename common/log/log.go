package log

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"

	traceL "github.com/2JYK2/go-lib/common/trace"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func GetTracer(name string) trace.Tracer {
	return otel.Tracer(name)
}

func ErrSave(ctx context.Context, msg string, err error, field ...zap.Field) {
	field = append(field, zap.NamedError(Err, err))
	field = getLogField(ctx, field...)
	errorObj.Error(msg, field...)
}

func EmailErrSave(ctx context.Context, msg string, code int64, solution string, err error, field ...zap.Field) {
	field = append(field,
		zap.Int64(Code, code),
		zap.String(Solution, solution),
		zap.NamedError(Err, err),
		zap.Int(Alarm, EmailAlarm),
	)
	field = getLogField(ctx, field...)
	errorObj.Error(msg, field...)
}

func ShortMsgErrSave(ctx context.Context, msg string, code int64, solution string, err error, field ...zap.Field) {
	field = append(field,
		zap.Int64(Code, code),
		zap.String(Solution, solution),
		zap.NamedError(Err, err),
		zap.Int(Alarm, ShortMsgAlarm),
	)
	field = getLogField(ctx, field...)
	errorObj.Error(msg, field...)
}

func PhoneErrSave(ctx context.Context, msg string, code int64, solution string, err error, field ...zap.Field) {
	field = append(field,
		zap.Int64(Code, code),
		zap.String(Solution, solution),
		zap.NamedError(Err, err),
		zap.Int(Alarm, PhoneAlarm),
	)
	field = getLogField(ctx, field...)
	errorObj.Error(msg, field...)
}

func WarnSave(ctx context.Context, msg string, err error, field ...zap.Field) {
	field = append(field, zap.NamedError(Err, err))
	field = getLogField(ctx, field...)
	runObj.Warn(msg, field...)
}

func InfoSave(ctx context.Context, msg string, field ...zap.Field) {
	field = getLogField(ctx, field...)
	runObj.Info(msg, field...)
}

func DebugSave(ctx context.Context, msg string, field ...zap.Field) {
	field = getLogField(ctx, field...)
	runObj.Debug(msg, field...)
}

func FatalSave(ctx context.Context, msg string, err error, field ...zap.Field) {
	field = append(field, zap.NamedError(Err, err))
	field = getLogField(ctx, field...)
	errorObj.Fatal(msg, field...)
}

func StatisSave(msg string, field ...zap.Field) {
	statisObj.Info(msg, field...)
}

func GetTraceIdFromContext(ctx context.Context) (string, string) {
	if ctx == nil {
		return "", ""
	}

	if ct, ok := ctx.(*gin.Context); ok {
		ctx = ct.Request.Context()
	}
	traceCtx := trace.SpanFromContext(ctx).SpanContext()
	if !traceCtx.TraceID().IsValid() {
		traceV := traceL.GetTraceValue()
		if traceV != nil {
			return traceV.TraceId, ""
		}
	}
	traceID := traceCtx.TraceID().String()
	spanID := traceCtx.SpanID().String()

	return traceID, spanID
}

func getReqIP(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if ct, ok := ctx.(*gin.Context); ok {
		ctx = ct.Request.Context()
	}
	return traceL.GetReqIP(ctx)
}

func getLogNonLineField(ctx context.Context, field ...zap.Field) []zap.Field {
	traceId, spanId := GetTraceIdFromContext(ctx)
	field = append(field, zap.String("trace_id", traceId), zap.String("span_id", spanId))

	reqIp := getReqIP(ctx)
	if reqIp != "" {
		field = append(field, zap.String(ReqIp, reqIp))
	}

	wsTraceId := traceL.GetOtherTraceID(ctx)
	if wsTraceId != "" {
		field = append(field, zap.String("SelfTraceId", wsTraceId))
	}
	return field
}

func getLogField(ctx context.Context, field ...zap.Field) []zap.Field {
	field = getLogNonLineField(ctx, field...)

	_, file, line, _ := runtime.Caller(2)
	p, _ := os.Getwd()
	field = append(field, zap.String(Line, fmt.Sprintf("%s:%d", strings.TrimPrefix(file, p), line)))
	return field
}

func GetLogger(name string, options ...any) *zap.Logger {
	var field []zap.Field

	skip := 0
	ctx := context.Background()
	for _, option := range options {
		if v, ok := option.(context.Context); ok {
			ctx = v
			continue
		}
		if v, ok := option.(int); ok {
			skip = v
			continue
		}
	}
	if skip < 1 {
		skip = 1
	}

	_, file, line, _ := runtime.Caller(skip)
	p, _ := os.Getwd()
	field = append(field, zap.String(Line, fmt.Sprintf("%s:%d", strings.TrimPrefix(file, p), line)))
	field = getLogNonLineField(ctx, field...)

	switch name {
	case ErrorLogName:
		return errorObj.With(field...)
	case RunLogName:
		return runObj.With(field...)
	case StatisLogName:
		return statisObj.With(field...)
	default:
		return nil
	}
}

// to lixianda
func GetLoggerByName(ctx context.Context, name string) *zap.Logger {
	return GetLogger(name, ctx, 3)
}
