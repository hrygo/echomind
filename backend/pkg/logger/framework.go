package logger

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// zapLoggerImpl 基于 zap 的日志实现
type zapLoggerImpl struct {
	config   *Config
	zap      *zap.Logger
	sugar    *zap.SugaredLogger
	providers []Provider

	// 采样控制
	sampler  *sampler
	mu       sync.RWMutex
}

// NewLogger 创建新的日志实例
func NewLogger(config *Config) (Logger, error) {
	if config == nil {
		config = DefaultConfig()
	}

	zapLogger, err := buildZapLogger(config)
	if err != nil {
		return nil, fmt.Errorf("failed to build zap logger: %w", err)
	}

	providers, err := buildProviders(config.Providers)
	if err != nil {
		return nil, fmt.Errorf("failed to build providers: %w", err)
	}

	loggerImpl := &zapLoggerImpl{
		config:    config,
		zap:       zapLogger,
		sugar:     zapLogger.Sugar(),
		providers: providers,
	}

	// 如果启用采样
	if config.Sampling.Enabled {
		loggerImpl.sampler = newSampler(config.Sampling)
	}

	return loggerImpl, nil
}

// buildZapLogger 构建 zap logger
func buildZapLogger(config *Config) (*zap.Logger, error) {
	var cores []zapcore.Core

	// 控制台输出
	if config.Output.Console.Enabled {
		consoleCore, err := buildConsoleCore(config)
		if err != nil {
			return nil, err
		}
		cores = append(cores, consoleCore)
	}

	// 文件输出
	if config.Output.File.Enabled {
		fileCore, err := buildFileCore(config)
		if err != nil {
			return nil, err
		}
		cores = append(cores, fileCore)
	}

	if len(cores) == 0 {
		// 默认控制台输出
		consoleCore, err := buildConsoleCore(config)
		if err != nil {
			return nil, err
		}
		cores = append(cores, consoleCore)
	}

	core := zapcore.NewTee(cores...)
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)), nil
}

// buildConsoleCore 构建控制台核心
func buildConsoleCore(config *Config) (zapcore.Core, error) {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	if config.Production {
		if config.Output.Console.Format == "json" {
			return zapcore.NewCore(
				zapcore.NewJSONEncoder(encoderConfig),
				zapcore.AddSync(os.Stdout),
				zap.NewAtomicLevelAt(convertLevel(config.Level)),
			), nil
		}
	}

	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zap.NewAtomicLevelAt(convertLevel(config.Level)),
	), nil
}

// buildFileCore 构建文件核心
func buildFileCore(config *Config) (zapcore.Core, error) {
	// 使用 lumberjack 进行日志轮转
	w := &lumberjack.Logger{
		Filename:   config.Output.File.Path,
		MaxSize:    config.Output.File.MaxSize,
		MaxBackups: config.Output.File.MaxBackups,
		MaxAge:     config.Output.File.MaxAge,
		Compress:   config.Output.File.Compress,
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(w),
		zap.NewAtomicLevelAt(convertLevel(config.Level)),
	), nil
}

// buildProviders 构建日志提供者
func buildProviders(providerConfigs []ProviderConfig) ([]Provider, error) {
	var providers []Provider

	for _, pc := range providerConfigs {
		if !pc.Enabled {
			continue
		}

		provider, err := createProvider(pc)
		if err != nil {
			return nil, fmt.Errorf("failed to create provider %s: %w", pc.Name, err)
		}

		providers = append(providers, provider)
	}

	return providers, nil
}

// createProvider 创建日志提供者
func createProvider(config ProviderConfig) (Provider, error) {
	// 这里可以根据 Type 创建不同的提供者
	// 例如: elasticsearch, loki, splunk 等
	switch config.Type {
	case "elasticsearch":
		return NewElasticsearchProvider(config.Settings)
	case "loki":
		return NewLokiProvider(config.Settings)
	case "splunk":
		return NewSplunkProvider(config.Settings)
	default:
		return NewNoopProvider(), nil
	}
}

// convertLevel 转换日志级别
func convertLevel(level Level) zapcore.Level {
	switch level {
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case FatalLevel:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// 实现 Logger 接口

func (l *zapLoggerImpl) Debug(msg string, fields ...Field) {
	l.log(DebugLevel, msg, fields...)
}

func (l *zapLoggerImpl) Info(msg string, fields ...Field) {
	l.log(InfoLevel, msg, fields...)
}

func (l *zapLoggerImpl) Warn(msg string, fields ...Field) {
	l.log(WarnLevel, msg, fields...)
}

func (l *zapLoggerImpl) Error(msg string, fields ...Field) {
	l.log(ErrorLevel, msg, fields...)
}

func (l *zapLoggerImpl) Fatal(msg string, fields ...Field) {
	l.log(FatalLevel, msg, fields...)
}

func (l *zapLoggerImpl) DebugContext(ctx context.Context, msg string, fields ...Field) {
	l.logWithContext(ctx, DebugLevel, msg, fields...)
}

func (l *zapLoggerImpl) InfoContext(ctx context.Context, msg string, fields ...Field) {
	l.logWithContext(ctx, InfoLevel, msg, fields...)
}

func (l *zapLoggerImpl) WarnContext(ctx context.Context, msg string, fields ...Field) {
	l.logWithContext(ctx, WarnLevel, msg, fields...)
}

func (l *zapLoggerImpl) ErrorContext(ctx context.Context, msg string, fields ...Field) {
	l.logWithContext(ctx, ErrorLevel, msg, fields...)
}

func (l *zapLoggerImpl) FatalContext(ctx context.Context, msg string, fields ...Field) {
	l.logWithContext(ctx, FatalLevel, msg, fields...)
}

func (l *zapLoggerImpl) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.config.Level = level
	// 这里需要重建 zap logger 以应用新级别
}

func (l *zapLoggerImpl) GetLevel() Level {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.config.Level
}

func (l *zapLoggerImpl) With(fields ...Field) Logger {
	interfaces := make([]interface{}, 0, len(fields)*2)
	for _, field := range fields {
		interfaces = append(interfaces, field.Key, field.Value)
	}

	return &zapLoggerImpl{
		config:    l.config,
		zap:       l.zap.With(convertFields(fields)...),
		sugar:     l.sugar.With(interfaces...),
		providers: l.providers,
		sampler:   l.sampler,
	}
}

// 内部方法

func (l *zapLoggerImpl) log(level Level, msg string, fields ...Field) {
	l.logWithContext(context.Background(), level, msg, fields...)
}

func (l *zapLoggerImpl) logWithContext(ctx context.Context, level Level, msg string, fields ...Field) {
	// 检查级别
	if level < l.GetLevel() {
		return
	}

	// 采样检查
	if l.sampler != nil && !l.sampler.shouldLog(level) {
		return
	}

	// 构建 zap 字段
	zapFields := l.buildFields(ctx, fields...)

	// 输出到 zap
	switch level {
	case DebugLevel:
		l.zap.Debug(msg, zapFields...)
	case InfoLevel:
		l.zap.Info(msg, zapFields...)
	case WarnLevel:
		l.zap.Warn(msg, zapFields...)
	case ErrorLevel:
		l.zap.Error(msg, zapFields...)
	case FatalLevel:
		l.zap.Fatal(msg, zapFields...)
	}

	// 输出到企业级提供者
	if len(l.providers) > 0 {
		entry := l.buildLogEntry(ctx, level, msg, fields...)
		for _, provider := range l.providers {
			_ = provider.Write(ctx, entry) // 异步写入，忽略错误
		}
	}
}

func (l *zapLoggerImpl) buildFields(ctx context.Context, fields ...Field) []zap.Field {
	var zapFields []zap.Field

	// 添加上下文字段
	if ctx != nil {
		if traceID := GetTraceID(ctx); traceID != "" {
			zapFields = append(zapFields, zap.String("trace_id", traceID))
		}
		if userID := GetUserID(ctx); userID != "" {
			zapFields = append(zapFields, zap.String("user_id", userID))
		}
		if orgID := GetOrgID(ctx); orgID != "" {
			zapFields = append(zapFields, zap.String("org_id", orgID))
		}
	}

	// 添加全局字段
	for k, v := range l.config.Context.GlobalFields {
		zapFields = append(zapFields, zap.Any(k, v))
	}

	// 添加传入字段
	zapFields = append(zapFields, convertFields(fields)...)

	return zapFields
}

func (l *zapLoggerImpl) buildLogEntry(ctx context.Context, level Level, msg string, fields ...Field) *LogEntry {
	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   msg,
		Fields:    make(map[string]interface{}),
		Context: ContextInfo{
			Extra: make(map[string]string),
		},
	}

	// 添加字段
	for _, field := range fields {
		entry.Fields[field.Key] = field.Value
	}

	// 添加上下文信息
	if ctx != nil {
		entry.Context.TraceID = GetTraceID(ctx)
		entry.Context.UserID = GetUserID(ctx)
		entry.Context.OrgID = GetOrgID(ctx)
		entry.Context.SessionID = GetSessionID(ctx)
	}

	// 添加源码信息
	if fn, file, line, ok := runtime.Caller(2); ok {
		entry.Source.Function = runtime.FuncForPC(fn).Name()
		entry.Source.File = file
		entry.Source.Line = line
	}

	return entry
}

// convertFields 转换字段
func convertFields(fields Fields) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Value)
	}
	return zapFields
}

// Close 关闭日志
func (l *zapLoggerImpl) Close() error {
	var lastErr error

	for _, provider := range l.providers {
		if err := provider.Close(); err != nil {
			lastErr = err
		}
	}

	_ = l.zap.Sync()
	return lastErr
}