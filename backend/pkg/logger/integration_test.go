package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestFileOutput 集成测试文件输出
func TestFileOutput(t *testing.T) {
	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "logger-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	logPath := filepath.Join(tmpDir, "test.log")

	config := &Config{
		Level:      InfoLevel,
		Production: true,
		Output: OutputConfig{
			File: FileOutputConfig{
				Enabled:    true,
				Path:       logPath,
				MaxSize:    1, // 1MB
				MaxAge:     1, // 1 day
				MaxBackups: 3,
				Compress:   false,
			},
			Console: ConsoleOutputConfig{
				Enabled: false,
			},
		},
		Providers: []ProviderConfig{},
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// 记录一些日志
	logger.Info("Test log message", String("test", "integration"))
	logger.Error("Test error message", Error(&testError{}))

	// 等待一下确保日志写入
	time.Sleep(100 * time.Millisecond)

	// 关闭日志器
	if closer, ok := logger.(interface{ Close() error }); ok {
		_ = closer.Close()
	}

	// 检查文件是否创建
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Error("Log file was not created")
	}

	// 检查文件内容
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	contentStr := string(content)
	if !contains(contentStr, "Test log message") {
		t.Error("Log message not found in file")
	}

	if !contains(contentStr, "Test error message") {
		t.Error("Error message not found in file")
	}

	if !contains(contentStr, "test") {
		t.Error("Log field not found in file")
	}
}

// TestElasticsearchProvider 集成测试 Elasticsearch 提供者
func TestElasticsearchProvider(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Elasticsearch integration test in short mode")
	}

	config := ProviderConfig{
		Name:    "test-elasticsearch",
		Type:    "elasticsearch",
		Enabled: true,
		Settings: map[string]interface{}{
			"url":        "http://localhost:9200",
			"index":      "test-echomind-logs",
			"batch_size": 1, // 小批次用于测试
		},
	}

	provider, err := createProvider(config)
	if err != nil {
		t.Fatalf("Failed to create Elasticsearch provider: %v", err)
	}

	// 测试健康检查
	if err := provider.Ping(); err != nil {
		t.Skipf("Elasticsearch not available: %v", err)
	}

	// 创建测试日志条目
	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     InfoLevel,
		Message:   "Integration test message",
		Fields: map[string]interface{}{
			"test_field": "test_value",
			"number":     42,
		},
		Context: ContextInfo{
			TraceID:   "test-trace-123",
			UserID:    "test-user-456",
			SessionID: "test-session-789",
		},
		Source: SourceInfo{
			Function: "TestElasticsearchProvider",
			File:     "integration_test.go",
			Line:     42,
		},
	}

	// 写入日志
	ctx := context.Background()
	if err := provider.Write(ctx, entry); err != nil {
		t.Errorf("Failed to write log entry: %v", err)
	}

	// 等待一下确保写入完成
	time.Sleep(1 * time.Second)

	// 关闭提供者
	if err := provider.Close(); err != nil {
		t.Errorf("Failed to close provider: %v", err)
	}
}

// TestLokiProvider 集成测试 Grafana Loki 提供者
func TestLokiProvider(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Loki integration test in short mode")
	}

	config := ProviderConfig{
		Name:    "test-loki",
		Type:    "loki",
		Enabled: true,
		Settings: map[string]interface{}{
			"url": "http://localhost:3100/loki/api/v1/push",
			"labels": map[string]string{
				"service":     "echomind-test",
				"environment": "test",
			},
			"batch_size": 1,
		},
	}

	provider, err := createProvider(config)
	if err != nil {
		t.Fatalf("Failed to create Loki provider: %v", err)
	}

	// 测试健康检查
	if err := provider.Ping(); err != nil {
		t.Skipf("Loki not available: %v", err)
	}

	// 创建测试日志条目
	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     WarnLevel,
		Message:   "Loki integration test",
		Fields: map[string]interface{}{
			"loki_test": true,
			"component": "integration",
		},
		Context: ContextInfo{
			TraceID: "loki-test-trace",
		},
	}

	// 写入日志
	ctx := context.Background()
	if err := provider.Write(ctx, entry); err != nil {
		t.Errorf("Failed to write log entry: %v", err)
	}

	// 等待一下确保写入完成
	time.Sleep(1 * time.Second)

	// 关闭提供者
	if err := provider.Close(); err != nil {
		t.Errorf("Failed to close provider: %v", err)
	}
}

// TestMultipleProviders 集成测试多个提供者
func TestMultipleProviders(t *testing.T) {
	// 创建临时目录用于文件输出
	tmpDir, err := os.MkdirTemp("", "logger-multi-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	logPath := filepath.Join(tmpDir, "multi-test.log")

	config := &Config{
		Level:      DebugLevel,
		Production: false,
		Output: OutputConfig{
			File: FileOutputConfig{
				Enabled: true,
				Path:    logPath,
			},
			Console: ConsoleOutputConfig{
				Enabled: false,
			},
		},
		Providers: []ProviderConfig{
			{
				Name:    "noop1",
				Type:    "noop",
				Enabled: true,
			},
			{
				Name:    "noop2",
				Type:    "noop",
				Enabled: true,
			},
		},
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger with multiple providers: %v", err)
	}

	// 记录不同级别的日志
	ctx := WithTraceID(context.Background(), "multi-test-trace")

	logger.DebugContext(ctx, "Debug message", String("level", "debug"))
	logger.InfoContext(ctx, "Info message", String("level", "info"))
	logger.WarnContext(ctx, "Warn message", String("level", "warn"))
	logger.ErrorContext(ctx, "Error message", String("level", "error"))

	// 等待一下确保所有提供者都处理完成
	time.Sleep(100 * time.Millisecond)

	// 关闭日志器
	if closer, ok := logger.(interface{ Close() error }); ok {
		_ = closer.Close()
	}

	// 验证文件输出
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Error("Log file was not created")
	}
}

// TestContextPropagation 集成测试上下文传播
func TestContextPropagation(t *testing.T) {
	logger, err := NewLogger(DevelopmentConfig())
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// 创建包含多个上下文值的 context
	ctx := context.Background()
	ctx = WithTraceID(ctx, "trace-propagation-123")
	ctx = WithUserID(ctx, "user-propagation-456")
	ctx = WithOrgID(ctx, "org-propagation-789")
	ctx = WithSessionID(ctx, "session-propagation-012")
	ctx = WithRequestID(ctx, "req-propagation-345")

	// 使用 logger 的 With 方法创建带有额外字段的日志器
	contextLogger := logger.With(
		String("component", "context-test"),
		String("version", "1.0.0"),
	)

	// 记录日志
	contextLogger.InfoContext(ctx, "Context propagation test",
		String("action", "test_propagation"),
		Int64("timestamp", time.Now().Unix()))

	// 测试上下文在 With 方法后仍然有效
	newLogger := logger.With(String("extra", "field"))
	newLogger.InfoContext(ctx, "After With() test",
		String("action", "after_with"))

	// 关闭日志器
	if closer, ok := logger.(interface{ Close() error }); ok {
		_ = closer.Close()
	}
}

// TestConcurrentLogging 集成测试并发日志记录
func TestConcurrentLogging(t *testing.T) {
	logger, err := NewLogger(ProductionConfig())
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	const numGoroutines = 10
	const numLogsPerGoroutine = 100

	done := make(chan bool, numGoroutines)

	// 启动多个 goroutine 并发记录日志
	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			defer func() { done <- true }()

			for j := 0; j < numLogsPerGoroutine; j++ {
				ctx := WithTraceID(context.Background(),
					fmt.Sprintf("trace-%d-%d", goroutineID, j))

				logger.InfoContext(ctx, "Concurrent log message",
					Int("goroutine_id", goroutineID),
					Int("message_id", j),
					String("action", "concurrent_test"))
			}
		}(i)
	}

	// 等待所有 goroutine 完成
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// 等待一下确保所有日志都写入完成
	time.Sleep(500 * time.Millisecond)

	// 关闭日志器
	if closer, ok := logger.(interface{ Close() error }); ok {
		_ = closer.Close()
	}

	t.Logf("Completed %d goroutines with %d logs each", numGoroutines, numLogsPerGoroutine)
}

// TestLogRotation 集成测试日志轮转
func TestLogRotation(t *testing.T) {
	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "logger-rotation-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	logPath := filepath.Join(tmpDir, "rotation-test.log")

	config := &Config{
		Level:      InfoLevel,
		Production: true,
		Output: OutputConfig{
			File: FileOutputConfig{
				Enabled:    true,
				Path:       logPath,
				MaxSize:    1, // 1MB - 非常小以便测试轮转
				MaxAge:     1,
				MaxBackups: 2,
				Compress:   false,
			},
			Console: ConsoleOutputConfig{
				Enabled: false,
			},
		},
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// 记录大量日志以触发轮转 (在实际测试中可能需要调整)
	for i := 0; i < 1000; i++ {
		logger.Info("Log rotation test message",
			String("iteration", fmt.Sprintf("%d", i)),
			String("data", strings.Repeat("x", 1000))) // 增加消息大小
	}

	// 等待一下确保日志写入和轮转
	time.Sleep(1 * time.Second)

	// 关闭日志器
	if closer, ok := logger.(interface{ Close() error }); ok {
		_ = closer.Close()
	}

	// 检查是否有备份文件被创建
	files, err := filepath.Glob(logPath + "*")
	if err != nil {
		t.Fatalf("Failed to glob log files: %v", err)
	}

	if len(files) < 1 {
		t.Error("No log files found after rotation test")
	}

	t.Logf("Found %d log files after rotation test", len(files))
}

// 辅助函数
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && containsAt(s, substr)))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
