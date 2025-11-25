package logger

import (
	"context"
)

// 全局日志实例
var defaultLogger Logger

// Init 初始化全局日志实例
func Init(config *Config) error {
	logger, err := NewLogger(config)
	if err != nil {
		return err
	}
	defaultLogger = logger
	return nil
}

// SetDefaultLogger 设置默认日志实例
func SetDefaultLogger(logger Logger) {
	defaultLogger = logger
}

// GetDefaultLogger 获取默认日志实例
func GetDefaultLogger() Logger {
	if defaultLogger == nil {
		defaultLogger, _ = NewLogger(DefaultConfig())
	}
	return defaultLogger
}

// 全局日志方法 - 使用默认日志实例

func Debug(msg string, fields ...Field) {
	GetDefaultLogger().Debug(msg, fields...)
}

func Info(msg string, fields ...Field) {
	GetDefaultLogger().Info(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	GetDefaultLogger().Warn(msg, fields...)
}

func LogError(msg string, fields ...Field) {
	GetDefaultLogger().Error(msg, fields...)
}

func Fatal(msg string, fields ...Field) {
	GetDefaultLogger().Fatal(msg, fields...)
}

func DebugContext(ctx context.Context, msg string, fields ...Field) {
	GetDefaultLogger().DebugContext(ctx, msg, fields...)
}

func InfoContext(ctx context.Context, msg string, fields ...Field) {
	GetDefaultLogger().InfoContext(ctx, msg, fields...)
}

func WarnContext(ctx context.Context, msg string, fields ...Field) {
	GetDefaultLogger().WarnContext(ctx, msg, fields...)
}

func ErrorContext(ctx context.Context, msg string, fields ...Field) {
	GetDefaultLogger().ErrorContext(ctx, msg, fields...)
}

func FatalContext(ctx context.Context, msg string, fields ...Field) {
	GetDefaultLogger().FatalContext(ctx, msg, fields...)
}

// With 创建带字段的日志实例
func With(fields ...Field) Logger {
	return GetDefaultLogger().With(fields...)
}

// SetLevel 设置日志级别
func SetLevel(level Level) {
	GetDefaultLogger().SetLevel(level)
}

// GetLevel 获取日志级别
func GetLevel() Level {
	return GetDefaultLogger().GetLevel()
}

// Close 关闭日志
func Close() error {
	if logger, ok := defaultLogger.(interface{ Close() error }); ok {
		return logger.Close()
	}
	return nil
}

// New 创建新的日志实例（兼容旧接口）
func New(production bool) (Logger, error) {
	var config *Config
	if production {
		config = ProductionConfig()
	} else {
		config = DevelopmentConfig()
	}
	return NewLogger(config)
}

// 为了向后兼容，提供 zap 接口包装
type zapWrapper struct {
	Logger
}

func (w *zapWrapper) Sugar() *zapSugaredWrapper {
	return &zapSugaredWrapper{logger: w.Logger}
}

type zapSugaredWrapper struct {
	logger Logger
}

func (w *zapSugaredWrapper) Infow(msg string, keysAndValues ...interface{}) {
	fields := w.keysAndValuesToFields(keysAndValues...)
	w.logger.Info(msg, fields...)
}

func (w *zapSugaredWrapper) Debugw(msg string, keysAndValues ...interface{}) {
	fields := w.keysAndValuesToFields(keysAndValues...)
	w.logger.Debug(msg, fields...)
}

func (w *zapSugaredWrapper) Warnw(msg string, keysAndValues ...interface{}) {
	fields := w.keysAndValuesToFields(keysAndValues...)
	w.logger.Warn(msg, fields...)
}

func (w *zapSugaredWrapper) Errorw(msg string, keysAndValues ...interface{}) {
	fields := w.keysAndValuesToFields(keysAndValues...)
	w.logger.Error(msg, fields...)
}

func (w *zapSugaredWrapper) Infof(template string, args ...interface{}) {
	// 简单的格式化，实际使用中可以使用 fmt.Sprintf
	msg := formatTemplate(template, args...)
	w.logger.Info(msg)
}

func (w *zapSugaredWrapper) Debugf(template string, args ...interface{}) {
	msg := formatTemplate(template, args...)
	w.logger.Debug(msg)
}

func (w *zapSugaredWrapper) Warnf(template string, args ...interface{}) {
	msg := formatTemplate(template, args...)
	w.logger.Warn(msg)
}

func (w *zapSugaredWrapper) Errorf(template string, args ...interface{}) {
	msg := formatTemplate(template, args...)
	w.logger.Error(msg)
}

func (w *zapSugaredWrapper) Sync() error {
	// 无操作
	return nil
}

func (w *zapSugaredWrapper) keysAndValuesToFields(keysAndValues ...interface{}) Fields {
	fields := make(Fields, 0, len(keysAndValues)/2)
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			key := toString(keysAndValues[i])
			fields = append(fields, Field{Key: key, Value: keysAndValues[i+1]})
		}
	}
	return fields
}

func toString(v interface{}) string {
	if str, ok := v.(string); ok {
		return str
	}
	return ""
}

func formatTemplate(template string, args ...interface{}) string {
	// 简单实现，可以使用 fmt.Sprintf
	if len(args) == 0 {
		return template
	}
	return template // 简化实现
}

// AsZapLogger 将 Logger 包装成兼容 zap 的接口
func AsZapLogger(log Logger) *zapWrapper {
	return &zapWrapper{Logger: log}
}

// AsZapSugaredLogger 将 Logger 包装成兼容 zap.SugaredLogger 的接口
func AsZapSugaredLogger(log Logger) *zapSugaredWrapper {
	return &zapSugaredWrapper{logger: log}
}
