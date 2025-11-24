package logger

import (
	"context"
	"runtime/debug"
	"time"
)

// String 创建字符串字段
func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

// Strings 创建字符串数组字段
func Strings(key string, value []string) Field {
	return Field{Key: key, Value: value}
}

// Int 创建整数字段
func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

// Int64 创建 64 位整数字段
func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

// Float64 创建浮点数字段
func Float64(key string, value float64) Field {
	return Field{Key: key, Value: value}
}

// Bool 创建布尔字段
func Bool(key string, value bool) Field {
	return Field{Key: key, Value: value}
}

// Any 创建任意类型字段
func Any(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// Duration 创建时间间隔字段
func Duration(key string, value time.Duration) Field {
	return Field{Key: key, Value: value.String()}
}

// Time 创建时间字段
func Time(key string, value time.Time) Field {
	return Field{Key: key, Value: value}
}

// Error 创建错误字段
func Error(err error) Field {
	if err == nil {
		return Field{Key: "error", Value: nil}
	}
	return Field{Key: "error", Value: err.Error()}
}

// Stack 创建堆栈字段
func Stack() Field {
	return Field{Key: "stack", Value: debug.Stack()}
}

// Err 创建错误字段（别名）
func Err(err error) Field {
	return Error(err)
}

// Str 是 String 的别名
var Str = String

// FieldsFromMap 从 map 创建字段
func FieldsFromMap(m map[string]interface{}) Fields {
	fields := make(Fields, 0, len(m))
	for k, v := range m {
		fields = append(fields, Field{Key: k, Value: v})
	}
	return fields
}

// FieldsFromStruct 从结构体创建字段
// 使用反射将结构体转换为字段
func FieldsFromStruct(obj interface{}) Fields {
	// 这里可以使用反射库实现
	// 为了性能，建议手动构建字段
	return Fields{}
}

// 操作便捷方法

// InfoContextWithFields 带上下文字段的信息日志
func InfoContextWithFields(ctx context.Context, msg string, fields map[string]interface{}) {
	logFields := make(Fields, 0, len(fields))
	for k, v := range fields {
		logFields = append(logFields, Field{Key: k, Value: v})
	}
	GetDefaultLogger().InfoContext(ctx, msg, logFields...)
}

// ErrorContextWithError 带上下文和错误的错误日志
func ErrorContextWithError(ctx context.Context, msg string, err error, fields ...Field) {
	allFields := append(fields, Error(err))
	GetDefaultLogger().ErrorContext(ctx, msg, allFields...)
}

// DebugContextWithFields 带上下文字段的调试日志
func DebugContextWithFields(ctx context.Context, msg string, fields map[string]interface{}) {
	logFields := make(Fields, 0, len(fields))
	for k, v := range fields {
		logFields = append(logFields, Field{Key: k, Value: v})
	}
	GetDefaultLogger().DebugContext(ctx, msg, logFields...)
}

// WarnContextWithFields 带上下文字段的警告日志
func WarnContextWithFields(ctx context.Context, msg string, fields map[string]interface{}) {
	logFields := make(Fields, 0, len(fields))
	for k, v := range fields {
		logFields = append(logFields, Field{Key: k, Value: v})
	}
	GetDefaultLogger().WarnContext(ctx, msg, logFields...)
}