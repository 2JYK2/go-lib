package azlog

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
)

type AzLogger interface {
	InfoContext(ctx context.Context) *zerolog.Event
	DebugContext(ctx context.Context) *zerolog.Event
	WarnContext(ctx context.Context) *zerolog.Event
	ErrorContext(ctx context.Context) *zerolog.Event
}

type Logger struct {
	zerolog.Logger
}

const traceKey = "trace_id"
const spanKey = "span_id"

func New(cfg *Config) *Logger {
	writers := []io.Writer{newRollingWriter(cfg.LogPath, cfg)}
	return &Logger{zerolog.New(io.MultiWriter(writers...)).Level(zerolog.Level(cfg.Level)).With().Timestamp().Logger()}
}

func newRollingWriter(path string, conf *Config) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:  path + "/" + conf.Filename,
		MaxSize:   conf.MaxSize,
		MaxAge:    conf.MaxAge,   // days
		Compress:  conf.Compress, // disabled by default
		LocalTime: conf.LocalTime,
	}
}

func traceIdFromContext(ctx context.Context) (string, string) {
	if ct, ok := ctx.(*gin.Context); ok {
		ctx = ct.Request.Context()
	}
	traceCtx := trace.SpanFromContext(ctx).SpanContext()
	return traceCtx.TraceID().String(), traceCtx.SpanID().String()
}

func (l *Logger) InfoContext(ctx context.Context) *zerolog.Event {
	traceId, spanId := traceIdFromContext(ctx)
	return l.Info().Str(traceKey, traceId).Str(spanKey, spanId)
}

func (l *Logger) DebugContext(ctx context.Context) *zerolog.Event {
	traceId, spanId := traceIdFromContext(ctx)
	return l.Debug().Str(traceKey, traceId).Str(spanKey, spanId)
}

func (l *Logger) WarnContext(ctx context.Context) *zerolog.Event {
	traceId, spanId := traceIdFromContext(ctx)
	return l.Warn().Str(traceKey, traceId).Str(spanKey, spanId)
}

func (l *Logger) ErrorContext(ctx context.Context) *zerolog.Event {
	traceId, spanId := traceIdFromContext(ctx)
	return l.Error().Str(traceKey, traceId).Str(spanKey, spanId)
}
