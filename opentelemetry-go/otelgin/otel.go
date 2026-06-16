package otelgin

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	logsdk "go.opentelemetry.io/otel/sdk/log"
	metersdk "go.opentelemetry.io/otel/sdk/metric"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"
	"time"
)

//var (
//	insecure     = os.Getenv("INSECURE_MODE")
//	collectorURL = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
//)

func newOtelConfig(opts ...Option) otelConfig {
	var c = otelConfig{}
	for _, opt := range opts {
		c = opt.apply(c)
	}
	return c
}

func getContext(ctx context.Context) context.Context {
	if c, ok := ctx.(*gin.Context); ok {
		return c.Request.Context()
	}
	return ctx
}

func SetupOtelSDK(options ...Option) (func(context.Context) error, error) {
	cfg := newOtelConfig(options...)
	var shutdownFuncs []func(context.Context) error
	shutdown := func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}
	// set up propagator
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// trace provider
	traceProvider, err := newTraceProvider(cfg)
	if err != nil {
		return nil, err
	}
	shutdownFuncs = append(shutdownFuncs, traceProvider.Shutdown)
	otel.SetTracerProvider(traceProvider)

	// meter provider
	meterProvider, err := newMeterProvider(cfg)
	if err != nil {
		return nil, err
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	// logger provider
	loggerProvider, err := newLogProvider(cfg)
	if err != nil {
		return nil, err
	}
	shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
	global.SetLoggerProvider(loggerProvider)
	return shutdown, nil
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTraceProvider(cfg otelConfig) (*tracesdk.TracerProvider, error) {
	var secureOption otlptracegrpc.Option

	if !cfg.insecure {
		secureOption = otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	} else {
		secureOption = otlptracegrpc.WithInsecure()
	}
	exporter, err := otlptracegrpc.New(
		context.Background(),
		secureOption,
		otlptracegrpc.WithEndpoint(cfg.collectorURL),
	)
	if err != nil {
		return nil, err
	}
	opts := []tracesdk.TracerProviderOption{
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		tracesdk.WithBatcher(exporter),
	}
	opts = append(opts, cfg.tfs...)
	traceProvider := tracesdk.NewTracerProvider(
		opts...,
	)
	return traceProvider, nil
}

func newMeterProvider(cfg otelConfig) (*metersdk.MeterProvider, error) {
	var secureOption otlpmetricgrpc.Option
	if !cfg.insecure {
		secureOption = otlpmetricgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	} else {
		secureOption = otlpmetricgrpc.WithInsecure()
	}
	exporter, err := otlpmetricgrpc.New(
		context.Background(),
		secureOption,
		otlpmetricgrpc.WithEndpoint(cfg.collectorURL),
	)
	if err != nil {
		return nil, err
	}
	opts := []metersdk.Option{
		metersdk.WithReader(metersdk.NewPeriodicReader(exporter, metersdk.WithInterval(10*time.Second))),
	}
	opts = append(opts, cfg.mfs...)
	meterProvider := metersdk.NewMeterProvider(
		opts...,
	)
	return meterProvider, nil
}

func newLogProvider(cfg otelConfig) (*logsdk.LoggerProvider, error) {
	var secureOption otlploggrpc.Option
	if !cfg.insecure {
		secureOption = otlploggrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	} else {
		secureOption = otlploggrpc.WithInsecure()
	}
	exporter, err := otlploggrpc.New(
		context.Background(),
		otlploggrpc.WithEndpoint(cfg.collectorURL),
		secureOption,
	)
	if err != nil {
		return nil, err
	}
	opts := []logsdk.LoggerProviderOption{
		logsdk.WithProcessor(logsdk.NewBatchProcessor(exporter)),
	}
	opts = append(opts, cfg.lfs...)
	logProvider := logsdk.NewLoggerProvider(
		opts...,
	)
	return logProvider, nil
}
