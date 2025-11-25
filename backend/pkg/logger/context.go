package logger

import (
	"context"
)

// 上下文键类型
type contextKey string

const (
	TraceIDKey   contextKey = "trace_id"
	UserIDKey    contextKey = "user_id"
	OrgIDKey     contextKey = "org_id"
	SessionIDKey contextKey = "session_id"
	RequestIDKey contextKey = "request_id"
)

// WithTraceID 在上下文中设置 Trace ID
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// GetTraceID 从上下文中获取 Trace ID
func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
		return traceID
	}
	return ""
}

// WithUserID 在上下文中设置 User ID
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// GetUserID 从上下文中获取 User ID
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}

// WithOrgID 在上下文中设置 Org ID
func WithOrgID(ctx context.Context, orgID string) context.Context {
	return context.WithValue(ctx, OrgIDKey, orgID)
}

// GetOrgID 从上下文中获取 Org ID
func GetOrgID(ctx context.Context) string {
	if orgID, ok := ctx.Value(OrgIDKey).(string); ok {
		return orgID
	}
	return ""
}

// WithSessionID 在上下文中设置 Session ID
func WithSessionID(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, SessionIDKey, sessionID)
}

// GetSessionID 从上下文中获取 Session ID
func GetSessionID(ctx context.Context) string {
	if sessionID, ok := ctx.Value(SessionIDKey).(string); ok {
		return sessionID
	}
	return ""
}

// WithRequestID 在上下文中设置 Request ID
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// GetRequestID 从上下文中获取 Request ID
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}

// WithContext 在上下文中设置多个值
func WithContext(ctx context.Context, traceID, userID, orgID, sessionID, requestID string) context.Context {
	ctx = WithTraceID(ctx, traceID)
	ctx = WithUserID(ctx, userID)
	ctx = WithOrgID(ctx, orgID)
	ctx = WithSessionID(ctx, sessionID)
	ctx = WithRequestID(ctx, requestID)
	return ctx
}
