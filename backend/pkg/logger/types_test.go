package logger

import (
	"testing"
	"time"
)

func TestLevel_String(t *testing.T) {
	tests := []struct {
		name     string
		level    Level
		expected string
	}{
		{"DebugLevel", DebugLevel, "DEBUG"},
		{"InfoLevel", InfoLevel, "INFO"},
		{"WarnLevel", WarnLevel, "WARN"},
		{"ErrorLevel", ErrorLevel, "ERROR"},
		{"FatalLevel", FatalLevel, "FATAL"},
		{"Unknown Level", Level(99), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.level.String(); got != tt.expected {
				t.Errorf("Level.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestField(t *testing.T) {
	tests := []struct {
		name  string
		field Field
	}{
		{"String field", String("key", "value")},
		{"Int field", Int("count", 42)},
		{"Bool field", Bool("enabled", true)},
		{"Error field", Error(nil)},
		{"Error field with value", Error(&customTestError{})},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.field.Key == "" {
				t.Error("Field Key should not be empty")
			}
		})
	}
}

func TestLogEntry(t *testing.T) {
	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     InfoLevel,
		Message:   "test message",
		Fields:    map[string]interface{}{"key": "value"},
		Context: ContextInfo{
			TraceID:   "trace-123",
			UserID:    "user-456",
			SessionID: "session-789",
		},
		Source: SourceInfo{
			Function: "TestFunction",
			File:     "test.go",
			Line:     42,
		},
	}

	if entry.Level.String() != "INFO" {
		t.Errorf("Expected INFO level, got %s", entry.Level.String())
	}

	if entry.Message != "test message" {
		t.Errorf("Expected 'test message', got %s", entry.Message)
	}

	if entry.Context.TraceID != "trace-123" {
		t.Errorf("Expected trace-123, got %s", entry.Context.TraceID)
	}
}

func TestFieldsFromMap(t *testing.T) {
	m := map[string]interface{}{
		"string_key": "string_value",
		"int_key":    42,
		"bool_key":   true,
	}

	fields := FieldsFromMap(m)

	if len(fields) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(fields))
	}

	// 验证字段存在
	fieldMap := make(map[string]interface{})
	for _, field := range fields {
		fieldMap[field.Key] = field.Value
	}

	if fieldMap["string_key"] != "string_value" {
		t.Error("String field not found or incorrect")
	}

	if fieldMap["int_key"] != 42 {
		t.Error("Int field not found or incorrect")
	}

	if fieldMap["bool_key"] != true {
		t.Error("Bool field not found or incorrect")
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		isValid bool
	}{
		{
			name:    "Default config",
			config:  DefaultConfig(),
			isValid: true,
		},
		{
			name:    "Production config",
			config:  ProductionConfig(),
			isValid: true,
		},
		{
			name:    "Development config",
			config:  DevelopmentConfig(),
			isValid: true,
		},
		{
			name: "Custom config",
			config: &Config{
				Level:      InfoLevel,
				Production: true,
				Output: OutputConfig{
					File: FileOutputConfig{
						Enabled: true,
						Path:    "/tmp/test.log",
					},
				},
				Providers: []ProviderConfig{},
			},
			isValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 基本验证 - 在实际实现中可以添加更多验证逻辑
			if tt.config == nil && tt.isValid {
				t.Error("Config should not be nil for valid test case")
			}
		})
	}
}

// 测试辅助类型
type customTestError struct{}

func (e *customTestError) Error() string {
	return "test error"
}