package logger

import (
	"context"
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name:    "Default config",
			config:  DefaultConfig(),
			wantErr: false,
		},
		{
			name:    "Production config",
			config:  ProductionConfig(),
			wantErr: false,
		},
		{
			name:    "Development config",
			config:  DevelopmentConfig(),
			wantErr: false,
		},
		{
			name:    "Nil config",
			config:  nil,
			wantErr: false, // 应该使用默认配置
		},
		{
			name: "Custom config with providers",
			config: &Config{
				Level:      InfoLevel,
				Production: false,
				Providers: []ProviderConfig{
					{
						Name:    "noop",
						Type:    "noop",
						Enabled: true,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log, err := NewLogger(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLogger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && log == nil {
				t.Error("NewLogger() returned nil logger without error")
			}
		})
	}
}

func TestGlobalLogger(t *testing.T) {
	// 保存原有的默认日志器
	originalLogger := defaultLogger
	defer func() {
		defaultLogger = originalLogger
	}()

	// 初始化全局日志器
	config := DefaultConfig()
	err := Init(config)
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// 测试全局方法
	Debug("debug message")
	Info("info message")
	Warn("warn message")
	LogError("error message")

	// 测试上下文方法
	ctx := WithTraceID(context.Background(), "test-trace")
	DebugContext(ctx, "debug message with context")
	InfoContext(ctx, "info message with context")
	WarnContext(ctx, "warn message with context")
	ErrorContext(ctx, "error message with context")

	// 测试 With 方法
	loggerWithFields := With(String("key", "value"))
	if loggerWithFields == nil {
		t.Error("With() returned nil logger")
	}

	// 测试级别设置和获取
	SetLevel(DebugLevel)
	if got := GetLevel(); got != DebugLevel {
		t.Errorf("GetLevel() = %v, want %v", got, DebugLevel)
	}

	SetLevel(InfoLevel)
	if got := GetLevel(); got != InfoLevel {
		t.Errorf("GetLevel() = %v, want %v", got, InfoLevel)
	}
}

func TestLoggerInterface(t *testing.T) {
	log, err := NewLogger(DefaultConfig())
	if err != nil {
		t.Fatalf("NewLogger() failed: %v", err)
	}

	// 测试所有日志级别
	log.Debug("debug message")
	log.Info("info message")
	log.Warn("warn message")
	log.Error("error message")
	// Note: 不测试 Fatal 因为它会调用 os.Exit

	// 测试带字段的日志
	log.Info("message with fields",
		String("string_key", "string_value"),
		Int("int_key", 42),
		Bool("bool_key", true),
		Error(nil),
	)

	// 测试上下文日志
	ctx := context.Background()
	ctx = WithTraceID(ctx, "test-trace")
	ctx = WithUserID(ctx, "test-user")
	ctx = WithOrgID(ctx, "test-org")

	log.InfoContext(ctx, "message with context",
		String("component", "test"))

	// 测试 With 方法
	withLogger := log.With(String("global", "field"))
	withLogger.Info("message with inherited fields")
}

func TestConfigLoading(t *testing.T) {
	// 测试默认配置
	config := DefaultConfig()
	if config.Level != InfoLevel {
		t.Errorf("DefaultConfig() Level = %v, want %v", config.Level, InfoLevel)
	}

	if config.Production != false {
		t.Errorf("DefaultConfig() Production = %v, want %v", config.Production, false)
	}

	// 测试生产配置
	prodConfig := ProductionConfig()
	if prodConfig.Level != InfoLevel {
		t.Errorf("ProductionConfig() Level = %v, want %v", prodConfig.Level, InfoLevel)
	}

	if prodConfig.Production != true {
		t.Errorf("ProductionConfig() Production = %v, want %v", prodConfig.Production, true)
	}

	// 测试开发配置
	devConfig := DevelopmentConfig()
	if devConfig.Level != DebugLevel {
		t.Errorf("DevelopmentConfig() Level = %v, want %v", devConfig.Level, DebugLevel)
	}

	if devConfig.Production != false {
		t.Errorf("DevelopmentConfig() Production = %v, want %v", devConfig.Production, false)
	}
}

func TestBackwardCompatibility(t *testing.T) {
	log, err := NewLogger(DefaultConfig())
	if err != nil {
		t.Fatalf("NewLogger() failed: %v", err)
	}

	// 测试 zap 包装器
	zapLog := AsZapLogger(log)
	if zapLog == nil {
		t.Error("AsZapLogger() returned nil")
	}

	sugarLog := zapLog.Sugar()
	if sugarLog == nil {
		t.Error("Sugar() returned nil")
	}

	// 测试 sugared logger 方法
	sugarLog.Infow("info message", "key", "value")
	sugarLog.Debugw("debug message", "key", "value")
	sugarLog.Warnw("warn message", "key", "value")
	sugarLog.Errorw("error message", "key", "value")

	sugarLog.Infof("info message %s", "formatted")
	sugarLog.Debugf("debug message %s", "formatted")
	sugarLog.Warnf("warn message %s", "formatted")
	sugarLog.Errorf("error message %s", "formatted")

	err = sugarLog.Sync()
	if err != nil {
		t.Errorf("Sync() failed: %v", err)
	}
}

func TestHelperFunctions(t *testing.T) {
	// 测试字段创建函数
	strField := String("key", "value")
	if strField.Key != "key" || strField.Value != "value" {
		t.Error("String() failed")
	}

	intField := Int("key", 42)
	if intField.Key != "key" || intField.Value != 42 {
		t.Error("Int() failed")
	}

	boolField := Bool("key", true)
	if boolField.Key != "key" || boolField.Value != true {
		t.Error("Bool() failed")
	}

	errorField := Error(nil)
	if errorField.Key != "error" || errorField.Value != nil {
		t.Error("Error(nil) failed")
	}

	testErr := &testError{}
	errorField = Error(testErr)
	if errorField.Key != "error" || errorField.Value != "test error" {
		t.Error("Error(err) failed")
	}

	durationField := Duration("key", time.Second)
	if durationField.Key != "key" {
		t.Error("Duration() failed")
	}

	// 测试便捷方法
	ctx := WithTraceID(context.Background(), "test-trace")
	InfoContextWithFields(ctx, "message", map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	})

	ErrorContextWithError(ctx, "message", testErr, String("key", "value"))
	DebugContextWithFields(ctx, "message", map[string]interface{}{
		"debug_key": "debug_value",
	})
	WarnContextWithFields(ctx, "message", map[string]interface{}{
		"warn_key": "warn_value",
	})
}

func TestSampling(t *testing.T) {
	config := DefaultConfig()
	config.Sampling = SamplingConfig{
		Enabled: true,
		Rate:    10, // 每秒最多 10 条日志
		Burst:   5,  // 突发最多 5 条
		Levels:  []Level{DebugLevel, InfoLevel},
	}

	log, err := NewLogger(config)
	if err != nil {
		t.Fatalf("NewLogger() failed: %v", err)
	}

	// 快速记录大量日志，测试采样功能
	for i := 0; i < 100; i++ {
		log.Debug("debug message", Int("counter", i))
		log.Info("info message", Int("counter", i))
	}
}

func BenchmarkNewLogger(b *testing.B) {
	config := DefaultConfig()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = NewLogger(config)
	}
}

func BenchmarkLogInfo(b *testing.B) {
	log, _ := NewLogger(DefaultConfig())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info("test message", String("key", "value"))
	}
}

func BenchmarkLogInfoContext(b *testing.B) {
	log, _ := NewLogger(DefaultConfig())
	ctx := WithTraceID(context.Background(), "test-trace")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.InfoContext(ctx, "test message", String("key", "value"))
	}
}

func BenchmarkLogWithFields(b *testing.B) {
	log, _ := NewLogger(DefaultConfig())
	ctx := WithTraceID(context.Background(), "test-trace")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.InfoContext(ctx, "test message",
			String("string_key", "value"),
			Int("int_key", 42),
			Bool("bool_key", true),
			String("component", "benchmark"),
		)
	}
}

// 测试用的错误类型
type testError struct{}

func (e *testError) Error() string {
	return "test error"
}
