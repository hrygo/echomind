// Package logger provides enterprise-grade logging framework for EchoMind.
// It supports structured logging, context propagation, and integration with
// enterprise log platforms like Elasticsearch, Loki, and Splunk.
package logger

import (
	"context"
	"time"
)

// Level 定义日志级别
type Level int8

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// String 返回日志级别的字符串表示
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// ParseLevel 解析日志级别字符串
func ParseLevel(s string) Level {
	switch s {
	case "DEBUG", "debug":
		return DebugLevel
	case "INFO", "info":
		return InfoLevel
	case "WARN", "warn":
		return WarnLevel
	case "ERROR", "error":
		return ErrorLevel
	case "FATAL", "fatal":
		return FatalLevel
	default:
		return InfoLevel
	}
}

// UnmarshalText 实现 encoding.TextUnmarshaler 接口
func (l *Level) UnmarshalText(text []byte) error {
	*l = ParseLevel(string(text))
	return nil
}

// Field 结构化日志字段
type Field struct {
	Key   string
	Value interface{}
}

// Fields 多个字段
type Fields []Field

// Logger 核心日志接口
type Logger interface {
	// 基础日志方法
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)

	// 带上下文的日志方法
	DebugContext(ctx context.Context, msg string, fields ...Field)
	InfoContext(ctx context.Context, msg string, fields ...Field)
	WarnContext(ctx context.Context, msg string, fields ...Field)
	ErrorContext(ctx context.Context, msg string, fields ...Field)
	FatalContext(ctx context.Context, msg string, fields ...Field)

	// 配置方法
	SetLevel(level Level)
	GetLevel() Level
	With(fields ...Field) Logger
}

// Provider 企业级日志提供者接口
type Provider interface {
	// 输出日志
	Write(ctx context.Context, entry *LogEntry) error

	// 关闭和清理
	Close() error

	// 健康检查
	Ping() error
}

// LogEntry 标准化日志条目
type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     Level                  `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	Context   ContextInfo            `json:"context,omitempty"`
	Source    SourceInfo             `json:"source,omitempty"`
}

// ContextInfo 上下文信息
type ContextInfo struct {
	TraceID   string            `json:"trace_id,omitempty"`
	SpanID    string            `json:"span_id,omitempty"`
	UserID    string            `json:"user_id,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
	SessionID string            `json:"session_id,omitempty"`
	OrgID     string            `json:"org_id,omitempty"`
	Extra     map[string]string `json:"extra,omitempty"`
}

// SourceInfo 源码信息
type SourceInfo struct {
	Function string `json:"function,omitempty"`
	File     string `json:"file,omitempty"`
	Line     int    `json:"line,omitempty"`
}

// Config 日志配置
type Config struct {
	// 基础配置
	Level      Level `yaml:"level"`
	Production bool  `yaml:"production"`

	// 输出配置
	Output OutputConfig `yaml:"output"`

	// 上下文配置
	Context ContextConfig `yaml:"context"`

	// 采样配置
	Sampling SamplingConfig `yaml:"sampling"`

	// 企业级集成
	Providers []ProviderConfig `yaml:"providers"`
}

// OutputConfig 输出配置
type OutputConfig struct {
	// 文件输出
	File FileOutputConfig `yaml:"file"`

	// 控制台输出
	Console ConsoleOutputConfig `yaml:"console"`
}

// FileOutputConfig 文件输出配置
type FileOutputConfig struct {
	Enabled    bool   `yaml:"enabled"`
	Path       string `yaml:"path"`
	MaxSize    int    `yaml:"max_size"` // MB
	MaxAge     int    `yaml:"max_age"`  // days
	MaxBackups int    `yaml:"max_backups"`
	Compress   bool   `yaml:"compress"`
}

// ConsoleOutputConfig 控制台输出配置
type ConsoleOutputConfig struct {
	Enabled bool   `yaml:"enabled"`
	Format  string `yaml:"format"` // "json" | "console"
	Color   bool   `yaml:"color"`
}

// ContextConfig 上下文配置
type ContextConfig struct {
	// 自动提取的上下文字段
	AutoFields []string `yaml:"auto_fields"`

	// 全局固定字段
	GlobalFields map[string]interface{} `yaml:"global_fields"`
}

// SamplingConfig 采样配置
type SamplingConfig struct {
	Enabled bool    `yaml:"enabled"`
	Rate    int     `yaml:"rate"`   // 采样率，每秒日志数
	Burst   int     `yaml:"burst"`  // 突发采样数
	Levels  []Level `yaml:"levels"` // 需要采样的级别
}

// ProviderConfig 日志提供者配置
type ProviderConfig struct {
	Name     string                 `yaml:"name"`
	Type     string                 `yaml:"type"` // "elasticsearch", "loki", "splunk", etc.
	Enabled  bool                   `yaml:"enabled"`
	Settings map[string]interface{} `yaml:"settings"`
}
