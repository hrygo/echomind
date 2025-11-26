package telemetry

import (
	"context"
	"fmt"
	"io"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// TelemetryConfig defines the OpenTelemetry configuration
type TelemetryConfig struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	ExporterType   string // "console", "file", "otlp"

	// File exporter configuration
	TracesFilePath  string
	MetricsFilePath string

	// OTLP exporter configuration
	OTLPEndpoint string
	OTLPInsecure bool

	// Sampling configuration
	SamplingType  string  // "always_on", "always_off", "traceidratio"
	SamplingRatio float64 // 0.0 - 1.0
}

// Telemetry manages OpenTelemetry providers
type Telemetry struct {
	TracerProvider *trace.TracerProvider
	MeterProvider  *metric.MeterProvider
	config         *TelemetryConfig

	// File handles for cleanup
	traceFile  io.WriteCloser
	metricFile io.WriteCloser

	// Track shutdown state
	shutdown bool
}

// InitTelemetry initializes OpenTelemetry with the given configuration
func InitTelemetry(ctx context.Context, cfg *TelemetryConfig) (*Telemetry, error) {
	if cfg == nil {
		return nil, fmt.Errorf("telemetry config cannot be nil")
	}

	// Set defaults
	if cfg.ServiceName == "" {
		cfg.ServiceName = "echomind-backend"
	}
	if cfg.ServiceVersion == "" {
		cfg.ServiceVersion = "v1.2.0"
	}
	if cfg.Environment == "" {
		cfg.Environment = "development"
	}
	if cfg.ExporterType == "" {
		cfg.ExporterType = "console"
	}
	if cfg.SamplingType == "" {
		cfg.SamplingType = "always_on"
	}
	if cfg.SamplingRatio == 0 {
		cfg.SamplingRatio = 1.0
	}

	// Create resource with service identification
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			semconv.DeploymentEnvironment(cfg.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	telemetry := &Telemetry{
		config: cfg,
	}

	// Initialize Trace Provider
	if err := telemetry.initTraceProvider(ctx, res); err != nil {
		return nil, fmt.Errorf("failed to initialize trace provider: %w", err)
	}

	// Initialize Metric Provider
	if err := telemetry.initMeterProvider(ctx, res); err != nil {
		telemetry.Shutdown(ctx) // Cleanup trace provider
		return nil, fmt.Errorf("failed to initialize meter provider: %w", err)
	}

	return telemetry, nil
}

// initTraceProvider initializes the trace provider with appropriate exporter
func (t *Telemetry) initTraceProvider(ctx context.Context, res *resource.Resource) error {
	var exporter trace.SpanExporter
	var err error

	switch t.config.ExporterType {
	case "console":
		exporter, err = stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
		)
	case "file":
		if t.config.TracesFilePath == "" {
			t.config.TracesFilePath = "./logs/traces.jsonl"
		}
		t.traceFile, err = os.OpenFile(t.config.TracesFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("failed to open trace file: %w", err)
		}
		exporter, err = stdouttrace.New(
			stdouttrace.WithWriter(t.traceFile),
		)
	case "otlp":
		opts := []otlptracehttp.Option{
			otlptracehttp.WithEndpoint(t.config.OTLPEndpoint),
		}
		if t.config.OTLPInsecure {
			opts = append(opts, otlptracehttp.WithInsecure())
		}
		exporter, err = otlptracehttp.New(ctx, opts...)
	default:
		return fmt.Errorf("unsupported exporter type: %s", t.config.ExporterType)
	}

	if err != nil {
		return fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Configure sampler
	var sampler trace.Sampler
	switch t.config.SamplingType {
	case "always_on":
		sampler = trace.AlwaysSample()
	case "always_off":
		sampler = trace.NeverSample()
	case "traceidratio":
		sampler = trace.TraceIDRatioBased(t.config.SamplingRatio)
	default:
		sampler = trace.AlwaysSample()
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
		trace.WithSampler(sampler),
	)

	otel.SetTracerProvider(tp)
	t.TracerProvider = tp

	return nil
}

// initMeterProvider initializes the meter provider with appropriate exporter
func (t *Telemetry) initMeterProvider(ctx context.Context, res *resource.Resource) error {
	var exporter metric.Exporter
	var err error

	switch t.config.ExporterType {
	case "console":
		exporter, err = stdoutmetric.New(
			stdoutmetric.WithPrettyPrint(),
		)
	case "file":
		if t.config.MetricsFilePath == "" {
			t.config.MetricsFilePath = "./logs/metrics.jsonl"
		}
		t.metricFile, err = os.OpenFile(t.config.MetricsFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("failed to open metrics file: %w", err)
		}
		exporter, err = stdoutmetric.New(
			stdoutmetric.WithWriter(t.metricFile),
		)
	case "otlp":
		opts := []otlpmetrichttp.Option{
			otlpmetrichttp.WithEndpoint(t.config.OTLPEndpoint),
		}
		if t.config.OTLPInsecure {
			opts = append(opts, otlpmetrichttp.WithInsecure())
		}
		exporter, err = otlpmetrichttp.New(ctx, opts...)
	default:
		return fmt.Errorf("unsupported exporter type: %s", t.config.ExporterType)
	}

	if err != nil {
		return fmt.Errorf("failed to create metric exporter: %w", err)
	}

	mp := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exporter)),
		metric.WithResource(res),
	)

	otel.SetMeterProvider(mp)
	t.MeterProvider = mp

	return nil
}

// Shutdown gracefully shuts down the telemetry providers
func (t *Telemetry) Shutdown(ctx context.Context) error {
	// Allow multiple shutdowns without error
	if t.shutdown {
		return nil
	}
	t.shutdown = true

	var errs []error

	if t.TracerProvider != nil {
		if err := t.TracerProvider.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("trace provider shutdown: %w", err))
		}
	}

	if t.MeterProvider != nil {
		if err := t.MeterProvider.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("meter provider shutdown: %w", err))
		}
	}

	// Close file handles
	if t.traceFile != nil {
		if err := t.traceFile.Close(); err != nil {
			errs = append(errs, fmt.Errorf("trace file close: %w", err))
		}
	}
	if t.metricFile != nil {
		if err := t.metricFile.Close(); err != nil {
			errs = append(errs, fmt.Errorf("metric file close: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("telemetry shutdown errors: %v", errs)
	}

	return nil
}
