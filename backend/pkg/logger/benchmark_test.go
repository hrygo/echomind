package logger

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
)

// BenchmarkLoggerCreation 测试日志器创建性能
func BenchmarkLoggerCreation(b *testing.B) {
	config := DefaultConfig()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		log, _ := NewLogger(config)
		_ = log
	}
}

// BenchmarkLoggerCreationProduction 生产环境日志器创建性能
func BenchmarkLoggerCreationProduction(b *testing.B) {
	config := ProductionConfig()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		log, _ := NewLogger(config)
		_ = log
	}
}

// BenchmarkSimpleLogging 测试简单日志记录性能
func BenchmarkSimpleLogging(b *testing.B) {
	log, _ := NewLogger(DefaultConfig())
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		log.Info("simple message")
	}
}

// BenchmarkLoggingWithFields 测试带字段的日志记录性能
func BenchmarkLoggingWithFields(b *testing.B) {
	log, _ := NewLogger(DefaultConfig())
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		log.Info("message with fields",
			String("key1", "value1"),
			Int("key2", 42),
			Bool("key3", true),
			Float64("key4", 3.14),
		)
	}
}

// BenchmarkLoggingWithContext 测试上下文日志记录性能
func BenchmarkLoggingWithContext(b *testing.B) {
	log, _ := NewLogger(DefaultConfig())
	ctx := WithTraceID(context.Background(), "benchmark-trace")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		log.InfoContext(ctx, "context message",
			String("iteration", fmt.Sprintf("%d", i)))
	}
}

// BenchmarkMultipleContextValues 测试多个上下文值的性能
func BenchmarkMultipleContextValues(b *testing.B) {
	log, _ := NewLogger(DefaultConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		ctx = WithTraceID(ctx, fmt.Sprintf("trace-%d", i))
		ctx = WithUserID(ctx, fmt.Sprintf("user-%d", i))
		ctx = WithOrgID(ctx, fmt.Sprintf("org-%d", i))
		ctx = WithSessionID(ctx, fmt.Sprintf("session-%d", i))
		ctx = WithRequestID(ctx, fmt.Sprintf("req-%d", i))

		log.InfoContext(ctx, "multi-context message")
	}
}

// BenchmarkWithMethod 测试 With 方法的性能
func BenchmarkWithMethod(b *testing.B) {
	log, _ := NewLogger(DefaultConfig())
	withLogger := log.With(
		String("global1", "value1"),
		String("global2", "value2"),
		String("global3", "value3"),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		withLogger.Info("message with inherited fields",
			String("local", fmt.Sprintf("value-%d", i)))
	}
}

// BenchmarkNestedWith 测试嵌套 With 方法的性能
func BenchmarkNestedWith(b *testing.B) {
	log, _ := NewLogger(DefaultConfig())
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l1 := log.With(String("level1", "value1"))
		l2 := l1.With(String("level2", "value2"))
		l3 := l2.With(String("level3", "value3"))

		l3.Info("nested with message")
	}
}

// BenchmarkDifferentLogLevels 测试不同日志级别的性能
func BenchmarkDifferentLogLevels(b *testing.B) {
	log, _ := NewLogger(DefaultConfig())

	levels := []Level{DebugLevel, InfoLevel, WarnLevel, ErrorLevel}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		level := levels[i%len(levels)]
		switch level {
		case DebugLevel:
			log.Debug("debug message")
		case InfoLevel:
			log.Info("info message")
		case WarnLevel:
			log.Warn("warn message")
		case ErrorLevel:
			log.Error("error message")
		}
	}
}

// BenchmarkProviderWriting 测试提供者写入性能
func BenchmarkProviderWriting(b *testing.B) {
	config := &Config{
		Level:      InfoLevel,
		Production: false,
		Output: OutputConfig{
			Console: ConsoleOutputConfig{Enabled: false},
		},
		Providers: []ProviderConfig{
			{
				Name:    "noop",
				Type:    "noop",
				Enabled: true,
			},
		},
	}

	log, _ := NewLogger(config)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		log.Info("provider test message",
			String("iteration", fmt.Sprintf("%d", i)),
			Int("number", i),
		)
	}
}

// BenchmarkMultipleProviders 测试多个提供者的性能
func BenchmarkMultipleProviders(b *testing.B) {
	providers := make([]ProviderConfig, 5)
	for i := 0; i < 5; i++ {
		providers[i] = ProviderConfig{
			Name:    fmt.Sprintf("noop-%d", i),
			Type:    "noop",
			Enabled: true,
		}
	}

	config := &Config{
		Level:      InfoLevel,
		Production: false,
		Output: OutputConfig{
			Console: ConsoleOutputConfig{Enabled: false},
		},
		Providers: providers,
	}

	log, _ := NewLogger(config)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		log.Info("multiple providers test",
			String("iteration", fmt.Sprintf("%d", i)))
	}
}

// BenchmarkSampling 测试采样性能
func BenchmarkSampling(b *testing.B) {
	config := &Config{
		Level:      InfoLevel,
		Production: false,
		Sampling: SamplingConfig{
			Enabled: true,
			Rate:    1000, // 每秒 1000 条
			Burst:   100,  // 突发 100 条
			Levels:  []Level{DebugLevel, InfoLevel},
		},
		Providers: []ProviderConfig{
			{Name: "noop", Type: "noop", Enabled: true},
		},
	}

	log, _ := NewLogger(config)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		log.Info("sampling test message")
	}
}

// BenchmarkLargeFieldValues 测试大字段值的性能
func BenchmarkLargeFieldValues(b *testing.B) {
	log, _ := NewLogger(DefaultConfig())
	largeString := strings.Repeat("x", 1000) // 1KB 字符串

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info("large field test",
			String("large_field", largeString),
			String("normal_field", "normal_value"),
		)
	}
}

// BenchmarkJSONSerialization 测试 JSON 序列化性能
func BenchmarkJSONSerialization(b *testing.B) {
	log, _ := NewLogger(DefaultConfig())

	// 模拟复杂的日志字段
	complexFields := []Field{
		String("string_field", "test_string"),
		Int("int_field", 42),
		Bool("bool_field", true),
		Float64("float_field", 3.14159),
		Time("time_field", time.Now()),
		Any("map_field", map[string]interface{}{
			"nested_string": "nested_value",
			"nested_int":    123,
			"nested_bool":   false,
		}),
		Error(&testError{}),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info("complex fields test", complexFields...)
	}
}

// BenchmarkConcurrentLogging 测试并发日志记录性能
func BenchmarkConcurrentLogging(b *testing.B) {
	log, _ := NewLogger(ProductionConfig())

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			ctx := WithTraceID(context.Background(), fmt.Sprintf("trace-%d", i))
			log.InfoContext(ctx, "concurrent test",
				String("goroutine", "benchmark"),
				Int("iteration", i))
			i++
		}
	})
}

// BenchmarkMemoryUsage 测试内存使用情况
func BenchmarkMemoryUsage(b *testing.B) {
	log, _ := NewLogger(DefaultConfig())

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		withLogger := log.With(
			String("session", fmt.Sprintf("session-%d", i%100)),
			String("user", fmt.Sprintf("user-%d", i%50)),
		)

		withLogger.Info("memory usage test",
			String("action", "benchmark"),
			Int("timestamp", int(time.Now().Unix())),
		)
	}
}

// BenchmarkFieldCreation 测试字段创建性能
func BenchmarkFieldCreation(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = String("string_key", "string_value")
		_ = Int("int_key", 42)
		_ = Bool("bool_key", true)
		_ = Float64("float_key", 3.14)
		_ = Duration("duration_key", time.Second)
		_ = Time("time_key", time.Now())
		_ = Error(&testError{})
	}
}

// BenchmarkContextOperations 测试上下文操作性能
func BenchmarkContextOperations(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx = WithTraceID(ctx, fmt.Sprintf("trace-%d", i))
		_ = GetTraceID(ctx)
	}
}

// BenchmarkLegacyCompatibility 测试向后兼容性能
func BenchmarkLegacyCompatibility(b *testing.B) {
	log, _ := NewLogger(DefaultConfig())
	sugar := AsZapSugaredLogger(log)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sugar.Infow("legacy compatibility test",
			"key1", "value1",
			"key2", i,
			"key3", true)
	}
}

// 运行所有基准测试的辅助函数
func RunAllBenchmarks(b *testing.B) {
	b.Run("LoggerCreation", BenchmarkLoggerCreation)
	b.Run("SimpleLogging", BenchmarkSimpleLogging)
	b.Run("LoggingWithFields", BenchmarkLoggingWithFields)
	b.Run("LoggingWithContext", BenchmarkLoggingWithContext)
	b.Run("WithMethod", BenchmarkWithMethod)
	b.Run("ProviderWriting", BenchmarkProviderWriting)
	b.Run("ConcurrentLogging", BenchmarkConcurrentLogging)
	b.Run("MemoryUsage", BenchmarkMemoryUsage)
}

// 性能比较测试：新框架 vs 原生 zap
func BenchmarkComparisonWithZap(b *testing.B) {
	b.Run("NewFramework", func(b *testing.B) {
		log, _ := NewLogger(DefaultConfig())
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			log.Info("test message", String("key", "value"), Int("number", i))
		}
	})

	b.Run("DirectZap", func(b *testing.B) {
		zapLogger, _ := buildZapLogger(DefaultConfig())
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			zapLogger.Info("test message",
				zap.String("key", "value"),
				zap.Int("number", i))
		}
	})
}

// 压力测试：大量日志记录
func BenchmarkStressTest(b *testing.B) {
	log, _ := NewLogger(ProductionConfig())

	// 预热
	for i := 0; i < 1000; i++ {
		log.Info("warmup", Int("i", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := WithTraceID(context.Background(), fmt.Sprintf("stress-%d", i))
		log.InfoContext(ctx, "stress test message",
			String("category", "stress"),
			Int("iteration", i),
			Float64("timestamp", float64(time.Now().UnixNano())),
			Bool("async", true),
			String("component", "benchmark"),
		)
	}
}
