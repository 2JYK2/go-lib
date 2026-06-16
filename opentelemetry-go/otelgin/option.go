package otelgin

import (
	"context"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

type otelConfig struct {
	tfs          []trace.TracerProviderOption
	lfs          []log.LoggerProviderOption
	mfs          []metric.Option
	insecure     bool
	collectorURL string
}

type Option interface {
	apply(otelConfig) otelConfig
}

type otelOptionFunc func(otelConfig) otelConfig

func (fn otelOptionFunc) apply(cfg otelConfig) otelConfig {
	return fn(cfg)
}

func WithInsecure() Option {
	return otelOptionFunc(func(config otelConfig) otelConfig {
		config.insecure = true
		return config
	})
}

func WithExporterOtlpEndpoint(url string) Option {
	return otelOptionFunc(func(config otelConfig) otelConfig {
		config.collectorURL = url
		return config
	})
}

func WithTracerOption(option trace.TracerProviderOption) Option {
	return otelOptionFunc(func(config otelConfig) otelConfig {
		config.tfs = append(config.tfs, option)
		return config
	})
}

func WithLoggerOption(option log.LoggerProviderOption) Option {
	return otelOptionFunc(func(config otelConfig) otelConfig {
		config.lfs = append(config.lfs, option)
		return config
	})
}

func WithMetricOption(option metric.Option) Option {
	return otelOptionFunc(func(config otelConfig) otelConfig {
		config.mfs = append(config.mfs, option)
		return config
	})
}

func FormatResource(serviceName, env string) (*resource.Resource, error) {
	return resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.DeploymentEnvironmentKey.String(env),
		),
		resource.WithFromEnv(),
		resource.WithHost(),
		resource.WithTelemetrySDK(),
	)
}
