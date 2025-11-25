package logger

import (
	"context"
	"testing"
)

func TestContextFunctions(t *testing.T) {
	ctx := context.Background()

	// 测试 WithTraceID 和 GetTraceID
	traceID := "trace-123456"
	ctx = WithTraceID(ctx, traceID)
	if got := GetTraceID(ctx); got != traceID {
		t.Errorf("GetTraceID() = %v, want %v", got, traceID)
	}

	// 测试 WithUserID 和 GetUserID
	userID := "user-789012"
	ctx = WithUserID(ctx, userID)
	if got := GetUserID(ctx); got != userID {
		t.Errorf("GetUserID() = %v, want %v", got, userID)
	}

	// 测试 WithOrgID 和 GetOrgID
	orgID := "org-345678"
	ctx = WithOrgID(ctx, orgID)
	if got := GetOrgID(ctx); got != orgID {
		t.Errorf("GetOrgID() = %v, want %v", got, orgID)
	}

	// 测试 WithSessionID 和 GetSessionID
	sessionID := "session-901234"
	ctx = WithSessionID(ctx, sessionID)
	if got := GetSessionID(ctx); got != sessionID {
		t.Errorf("GetSessionID() = %v, want %v", got, sessionID)
	}

	// 测试 WithRequestID 和 GetRequestID
	requestID := "req-567890"
	ctx = WithRequestID(ctx, requestID)
	if got := GetRequestID(ctx); got != requestID {
		t.Errorf("GetRequestID() = %v, want %v", got, requestID)
	}
}

func TestWithContext(t *testing.T) {
	ctx := context.Background()

	// 使用 WithContext 一次性设置多个值
	ctx = WithContext(ctx, "trace-123", "user-456", "org-789", "session-012", "req-345")

	tests := []struct {
		name     string
		getFunc  func(context.Context) string
		expected string
	}{
		{"TraceID", GetTraceID, "trace-123"},
		{"UserID", GetUserID, "user-456"},
		{"OrgID", GetOrgID, "org-789"},
		{"SessionID", GetSessionID, "session-012"},
		{"RequestID", GetRequestID, "req-345"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.getFunc(ctx); got != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, got, tt.expected)
			}
		})
	}
}

func TestEmptyContextValues(t *testing.T) {
	ctx := context.Background()

	// 测试从空上下文中获取值
	if got := GetTraceID(ctx); got != "" {
		t.Errorf("GetTraceID() from empty context = %v, want empty string", got)
	}

	if got := GetUserID(ctx); got != "" {
		t.Errorf("GetUserID() from empty context = %v, want empty string", got)
	}

	if got := GetOrgID(ctx); got != "" {
		t.Errorf("GetOrgID() from empty context = %v, want empty string", got)
	}

	if got := GetSessionID(ctx); got != "" {
		t.Errorf("GetSessionID() from empty context = %v, want empty string", got)
	}

	if got := GetRequestID(ctx); got != "" {
		t.Errorf("GetRequestID() from empty context = %v, want empty string", got)
	}
}

func TestContextOverride(t *testing.T) {
	ctx := context.Background()

	// 设置初始值
	ctx = WithTraceID(ctx, "original-trace")
	ctx = WithUserID(ctx, "original-user")

	// 验证初始值
	if got := GetTraceID(ctx); got != "original-trace" {
		t.Errorf("Initial trace ID = %v, want original-trace", got)
	}

	// 覆盖值
	ctx = WithTraceID(ctx, "new-trace")
	if got := GetTraceID(ctx); got != "new-trace" {
		t.Errorf("Overridden trace ID = %v, want new-trace", got)
	}

	// 验证其他值未受影响
	if got := GetUserID(ctx); got != "original-user" {
		t.Errorf("User ID should remain unchanged = %v, want original-user", got)
	}
}

func TestContextKeyType(t *testing.T) {
	// 验证 context key 的唯一性
	ctx := context.Background()

	ctx1 := WithTraceID(ctx, "trace-1")
	ctx2 := WithUserID(ctx, "user-1")

	// 确保不同的键不会互相覆盖
	if got := GetTraceID(ctx1); got != "trace-1" {
		t.Errorf("Trace ID from ctx1 = %v, want trace-1", got)
	}

	if got := GetTraceID(ctx2); got != "" {
		t.Errorf("Trace ID from ctx2 (should be empty) = %v, want empty", got)
	}

	if got := GetUserID(ctx1); got != "" {
		t.Errorf("User ID from ctx1 (should be empty) = %v, want empty", got)
	}

	if got := GetUserID(ctx2); got != "user-1" {
		t.Errorf("User ID from ctx2 = %v, want user-1", got)
	}
}

func BenchmarkWithTraceID(b *testing.B) {
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WithTraceID(ctx, "trace-123456")
	}
}

func BenchmarkGetTraceID(b *testing.B) {
	ctx := WithTraceID(context.Background(), "trace-123456")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetTraceID(ctx)
	}
}

func BenchmarkWithContext(b *testing.B) {
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WithContext(ctx, "trace-123", "user-456", "org-789", "session-012", "req-345")
	}
}
