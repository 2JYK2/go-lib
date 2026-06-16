package trace

import (
	"context"
	"github.com/2JYK2/go-lib/common/trace/internal/routine"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"sync"
)

type Value struct {
	TraceId string
	Data    map[string]zap.Field
}

var traceLocal = routine.NewLocalStorage()

func Initialize(fields ...zap.Field) {
	tv := newTraceValue()
	for _, field := range fields {
		if field.Key != "" {
			tv.Data[field.Key] = field
		}
	}
}

func setTraceValue(traceId string) *Value {
	if trace := traceLocal.Get(); trace != nil {
		tv := trace.(*Value)
		tv.TraceId = traceId
		for k := range tv.Data {
			delete(tv.Data, k)
		}
		return tv
	} else {
		tv := &Value{
			TraceId: traceId,
			Data:    make(map[string]zap.Field),
		}
		traceLocal.Set(tv)
		return tv
	}
}

func InitializeByTrace(traceId string, fields ...zap.Field) {
	tv := setTraceValue(traceId)
	for _, field := range fields {
		if field.Key != "" {
			tv.Data[field.Key] = field
		}
	}
}

func Add(fields ...zap.Field) {
	tv := GetTraceValue()
	if tv != nil {
		for _, field := range fields {
			if field.Key != "" {
				tv.Data[field.Key] = field
			}
		}
	}
}

func GetTraceValue() *Value {
	if trace := traceLocal.Get(); trace != nil {
		return trace.(*Value)
	}
	return nil
}

func GetTraceValueByKey(key string) zap.Field {
	if trace := traceLocal.Get(); trace != nil {
		if value, ok := trace.(*Value).Data[key]; ok {
			return value
		}

	}
	return zap.Field{}
}

var lock sync.Mutex

// Deep Copy
func CopyAndSetTraceValue(tv *Value) {
	if tv == nil {
		return
	}
	lock.Lock()
	defer lock.Unlock()
	deepCopy := &Value{
		TraceId: tv.TraceId,
		Data:    make(map[string]zap.Field),
	}
	for k, v := range tv.Data {
		deepCopy.Data[k] = v
	}
	traceLocal.Set(deepCopy)
}

func newTraceValue() *Value {
	if trace := traceLocal.Get(); trace != nil {
		tv := trace.(*Value)
		tv.TraceId = uuid.New().String()
		for k := range tv.Data {
			delete(tv.Data, k)
		}
		return tv
	} else {
		tv := &Value{
			TraceId: uuid.New().String(),
			Data:    make(map[string]zap.Field),
		}
		traceLocal.Set(tv)
		return tv
	}
}

func Delete() {
	traceLocal.Del()
}

type contextKey string

const selfTraceIDKey contextKey = "SelfTraceID"

// WithTraceID set traceID context
func WithOtherTraceID(ctx context.Context, traceID string) context.Context {
	if traceID == "" {
		return ctx
	}

	return context.WithValue(ctx, selfTraceIDKey, traceID)
}

// TraceID context get traceID
func GetOtherTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if traceID, ok := ctx.Value(selfTraceIDKey).(string); ok {
		return traceID
	}
	return ""
}
