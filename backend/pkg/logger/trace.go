package logger

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// TraceContextFields 从 context 中提取 trace 信息作为日志字段
func TraceContextFields(ctx context.Context) []Field {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return nil
	}

	spanCtx := span.SpanContext()
	fields := make([]Field, 0, 2)

	if spanCtx.HasTraceID() {
		fields = append(fields, String("trace_id", spanCtx.TraceID().String()))
	}

	if spanCtx.HasSpanID() {
		fields = append(fields, String("span_id", spanCtx.SpanID().String()))
	}

	return fields
}

// WithTraceContext 创建带 trace 信息的日志实例
func WithTraceContext(ctx context.Context) Logger {
	fields := TraceContextFields(ctx)
	if len(fields) == 0 {
		return GetDefaultLogger()
	}
	return GetDefaultLogger().With(fields...)
}

// DebugCtx 使用 trace context 记录 debug 日志
func DebugCtx(ctx context.Context, msg string, fields ...Field) {
	allFields := append(TraceContextFields(ctx), fields...)
	GetDefaultLogger().Debug(msg, allFields...)
}

// InfoCtx 使用 trace context 记录 info 日志
func InfoCtx(ctx context.Context, msg string, fields ...Field) {
	allFields := append(TraceContextFields(ctx), fields...)
	GetDefaultLogger().Info(msg, allFields...)
}

// WarnCtx 使用 trace context 记录 warn 日志
func WarnCtx(ctx context.Context, msg string, fields ...Field) {
	allFields := append(TraceContextFields(ctx), fields...)
	GetDefaultLogger().Warn(msg, allFields...)
}

// ErrorCtx 使用 trace context 记录 error 日志
func ErrorCtx(ctx context.Context, msg string, fields ...Field) {
	allFields := append(TraceContextFields(ctx), fields...)
	GetDefaultLogger().Error(msg, allFields...)
}

// FatalCtx 使用 trace context 记录 fatal 日志
func FatalCtx(ctx context.Context, msg string, fields ...Field) {
	allFields := append(TraceContextFields(ctx), fields...)
	GetDefaultLogger().Fatal(msg, allFields...)
}
