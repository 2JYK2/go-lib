package gorm

import (
	"github.com/2JYK2/go-lib/opentelemetry-go/otelgin/gorm/tracing"
	"gorm.io/gorm"
)

func GormTraceWithoutMetric() gorm.Plugin {
	return tracing.NewPlugin(tracing.WithoutMetrics())
}
