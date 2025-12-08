package telemetry

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func TestTelemetryInitialization(t *testing.T) {
	tests := []struct {
		name    string
		config  *TelemetryConfig
		wantErr bool
	}{
		{
			name: "valid console exporter",
			config: &TelemetryConfig{
				ServiceName:    "test-service",
				ServiceVersion: "v1.0.0",
				Environment:    "test",
				ExporterType:   "console",
				SamplingType:   "always_on",
				SamplingRatio:  1.0,
			},
			wantErr: false,
		},
		{
			name: "valid file exporter",
			config: &TelemetryConfig{
				ServiceName:     "test-service",
				ExporterType:    "file",
				TracesFilePath:  "/tmp/test-traces.jsonl",
				MetricsFilePath: "/tmp/test-metrics.jsonl",
			},
			wantErr: false,
		},
		{
			name: "invalid exporter type",
			config: &TelemetryConfig{
				ServiceName:  "test-service",
				ExporterType: "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tel, err := InitTelemetry(context.Background(), tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, tel)
			assert.NotNil(t, tel.TracerProvider)
			assert.NotNil(t, tel.MeterProvider)

			// Cleanup
			err = tel.Shutdown(context.Background())
			assert.NoError(t, err)
		})
	}
}

func TestTracerCreation(t *testing.T) {
	// Initialize telemetry
	cfg := &TelemetryConfig{
		ServiceName:  "test-service",
		ExporterType: "console",
	}
	tel, err := InitTelemetry(context.Background(), cfg)
	require.NoError(t, err)
	defer tel.Shutdown(context.Background())

	// Create tracer
	tracer := otel.Tracer("test.tracer")
	require.NotNil(t, tracer)

	// Create span
	_, span := tracer.Start(context.Background(), "test-operation")
	assert.NotNil(t, span)

	// Add attributes
	span.SetAttributes(
		attribute.String("test.key", "test-value"),
		attribute.Int("test.count", 42),
	)

	// End span
	span.End()

	// Verify span context
	spanCtx := span.SpanContext()
	assert.True(t, spanCtx.IsValid())
}

func TestMetricsCreation(t *testing.T) {
	// Initialize telemetry
	cfg := &TelemetryConfig{
		ServiceName:  "test-service",
		ExporterType: "console",
	}
	tel, err := InitTelemetry(context.Background(), cfg)
	require.NoError(t, err)
	defer tel.Shutdown(context.Background())

	// Create search metrics
	metrics, err := NewSearchMetrics(context.Background())
	require.NoError(t, err)
	require.NotNil(t, metrics)

	ctx := context.Background()

	// Test counter
	metrics.IncrementSearchRequests(ctx)
	metrics.IncrementSearchRequests(ctx)

	// Test histogram
	metrics.RecordSearchLatency(ctx, 100*time.Millisecond)
	metrics.RecordEmbeddingLatency(ctx, 50*time.Millisecond)

	// Test up/down counter
	metrics.IncrementActiveSearches(ctx)
	metrics.IncrementActiveSearches(ctx)
	metrics.DecrementActiveSearches(ctx)

	// Test cache metrics
	metrics.IncrementCacheHits(ctx)
	metrics.IncrementCacheMisses(ctx)
}

func TestSpanStatus(t *testing.T) {
	cfg := &TelemetryConfig{
		ServiceName:  "test-service",
		ExporterType: "console",
	}
	tel, err := InitTelemetry(context.Background(), cfg)
	require.NoError(t, err)
	defer tel.Shutdown(context.Background())

	tracer := otel.Tracer("test.tracer")

	t.Run("success status", func(t *testing.T) {
		_, span := tracer.Start(context.Background(), "success-op")
		span.SetStatus(codes.Ok, "operation successful")
		span.End()
	})

	t.Run("error status", func(t *testing.T) {
		_, span := tracer.Start(context.Background(), "error-op")
		span.RecordError(assert.AnError)
		span.SetStatus(codes.Error, "operation failed")
		span.End()
	})
}

func TestNestedSpans(t *testing.T) {
	cfg := &TelemetryConfig{
		ServiceName:  "test-service",
		ExporterType: "console",
	}
	tel, err := InitTelemetry(context.Background(), cfg)
	require.NoError(t, err)
	defer tel.Shutdown(context.Background())

	tracer := otel.Tracer("test.tracer")

	// Parent span
	ctx, parentSpan := tracer.Start(context.Background(), "parent-operation")
	parentSpan.SetAttributes(attribute.String("level", "parent"))

	// Child span 1
	ctx, child1 := tracer.Start(ctx, "child-operation-1")
	child1.SetAttributes(attribute.String("level", "child"))
	child1.End()

	// Child span 2
	_, child2 := tracer.Start(ctx, "child-operation-2")
	child2.SetAttributes(attribute.String("level", "child"))
	child2.End()

	parentSpan.End()

	// Verify parent-child relationship
	parentCtx := parentSpan.SpanContext()
	child1Ctx := child1.SpanContext()
	assert.Equal(t, parentCtx.TraceID(), child1Ctx.TraceID())
}

func TestSamplingStrategies(t *testing.T) {
	tests := []struct {
		name          string
		samplingType  string
		samplingRatio float64
	}{
		{"always_on", "always_on", 1.0},
		{"always_off", "always_off", 0.0},
		{"50_percent", "traceidratio", 0.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &TelemetryConfig{
				ServiceName:   "test-service",
				ExporterType:  "console",
				SamplingType:  tt.samplingType,
				SamplingRatio: tt.samplingRatio,
			}

			tel, err := InitTelemetry(context.Background(), cfg)
			require.NoError(t, err)
			defer tel.Shutdown(context.Background())

			tracer := otel.Tracer("test.tracer")
			_, span := tracer.Start(context.Background(), "test-op")
			span.End()
		})
	}
}

func TestTelemetryShutdown(t *testing.T) {
	cfg := &TelemetryConfig{
		ServiceName:  "test-service",
		ExporterType: "console",
	}

	tel, err := InitTelemetry(context.Background(), cfg)
	require.NoError(t, err)

	// Shutdown should not error
	err = tel.Shutdown(context.Background())
	assert.NoError(t, err)

	// Multiple shutdowns should not error
	err = tel.Shutdown(context.Background())
	assert.NoError(t, err)
}

func BenchmarkSpanCreation(b *testing.B) {
	cfg := &TelemetryConfig{
		ServiceName:  "bench-service",
		ExporterType: "console",
	}
	tel, _ := InitTelemetry(context.Background(), cfg)
	defer tel.Shutdown(context.Background())

	tracer := otel.Tracer("bench.tracer")
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, span := tracer.Start(ctx, "bench-operation")
		span.End()
	}
}

func BenchmarkMetricsRecording(b *testing.B) {
	cfg := &TelemetryConfig{
		ServiceName:  "bench-service",
		ExporterType: "console",
	}
	tel, _ := InitTelemetry(context.Background(), cfg)
	defer tel.Shutdown(context.Background())

	metrics, _ := NewSearchMetrics(context.Background())
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		metrics.IncrementSearchRequests(ctx)
		metrics.RecordSearchLatency(ctx, time.Millisecond)
	}
}
