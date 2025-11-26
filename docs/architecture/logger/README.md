# EchoMind 企业级日志框架

## 概述

EchoMind 日志框架是一个企业级的日志解决方案，提供统一的日志接口、结构化日志记录、企业日志平台集成等特性。

## 特性

- ✅ **统一接口**: 标准化的 Logger 接口，支持上下文感知日志
- ✅ **结构化日志**: 支持字段式日志记录，便于查询和分析
- ✅ **多级输出**: 同时支持控制台、文件输出
- ✅ **日志轮转**: 自动日志轮转和压缩
- ✅ **上下文传递**: 自动传递 TraceID、UserID、OrgID 等上下文信息
- ✅ **企业集成**: 支持 Elasticsearch、Loki、Splunk 等企业日志平台
- ✅ **采样控制**: 防止日志洪水，支持采样配置
- ✅ **向后兼容**: 兼容现有 zap 接口

## 快速开始

### 1. 初始化日志框架

```go
import "github.com/hrygo/echomind/pkg/logger"

// 使用默认配置
if err := logger.Init(logger.DevelopmentConfig()); err != nil {
    log.Fatal("Failed to initialize logger:", err)
}

// 或从环境变量加载配置
config := logger.LoadConfigFromEnv()
if err := logger.Init(config); err != nil {
    log.Fatal("Failed to initialize logger:", err)
}
```

### 2. 基础日志记录

```go
import "github.com/hrygo/echomind/pkg/logger"

// 简单日志
logger.Info("Server started")
logger.Error("Database connection failed", logger.Error(err))

// 带字段的结构化日志
logger.Info("User logged in",
    logger.String("user_id", "12345"),
    logger.String("email", "user@example.com"),
    logger.Duration("response_time", time.Second))
```

### 3. 上下文感知日志

```go
// 在上下文中添加信息
ctx = logger.WithTraceID(ctx, "trace-123")
ctx = logger.WithUserID(ctx, "user-456")
ctx = logger.WithOrgID(ctx, "org-789")

// 自动提取上下文信息
logger.InfoContext(ctx, "Processing request",
    logger.String("endpoint", "/api/emails"))
```

## 配置

### 环境变量

- `LOG_LEVEL`: 日志级别 (DEBUG, INFO, WARN, ERROR, FATAL)
- `LOG_PRODUCTION`: 是否生产模式 (true/false)
- `LOG_FILE_PATH`: 日志文件路径
- `LOG_CONSOLE_FORMAT`: 控制台格式 (console/json)
- `LOG_CONSOLE_COLOR`: 是否彩色输出 (true/false)

### 配置文件示例

```yaml
level: INFO
production: false

output:
  file:
    enabled: true
    path: "logs/app.log"
    max_size: 100
    max_age: 7
  console:
    enabled: true
    format: "console"
    color: true

context:
  auto_fields: ["trace_id", "user_id", "org_id"]
  global_fields:
    service: "echomind"
    version: "1.1.0"

providers:
  - name: "elasticsearch"
    type: "elasticsearch"
    enabled: true
    settings:
      url: "http://localhost:9200"
      index: "echomind-logs"
```

## 企业日志平台集成

### Elasticsearch

```go
// 配置 Elasticsearch 提供者
config := logger.Config{
    Providers: []logger.ProviderConfig{
        {
            Name:    "elasticsearch",
            Type:    "elasticsearch",
            Enabled: true,
            Settings: map[string]interface{}{
                "url":        "http://localhost:9200",
                "index":      "echomind-logs",
                "batch_size": 100,
            },
        },
    },
}
```

### Grafana Loki

```go
// 配置 Loki 提供者
config := logger.Config{
    Providers: []logger.ProviderConfig{
        {
            Name:    "loki",
            Type:    "loki",
            Enabled: true,
            Settings: map[string]interface{}{
                "url": "http://localhost:3100/loki/api/v1/push",
                "labels": map[string]string{
                    "service":     "echomind",
                    "environment": "production",
                },
            },
        },
    },
}
```

### Splunk

```go
// 配置 Splunk HEC 提供者
config := logger.Config{
    Providers: []logger.ProviderConfig{
        {
            Name:    "splunk",
            Type:    "splunk",
            Enabled: true,
            Settings: map[string]interface{}{
                "url":    "https://splunk.example.com:8088/services/collector/event",
                "token":  "your-hec-token",
                "index":  "echomind",
                "source": "backend",
            },
        },
    },
}
```

## 最佳实践

### 1. 使用结构化字段

```go
// ✅ 推荐
logger.Info("User action completed",
    logger.String("user_id", userID.String()),
    logger.String("action", "email_sent"),
    logger.Int("email_count", 5),
    logger.Duration("duration", time.Since(start)))

// ❌ 不推荐
logger.Info(fmt.Sprintf("User %s completed action %s with %d emails in %v",
    userID, action, count, duration))
```

### 2. 传递上下文

```go
// 在处理器中设置上下文
func HandleRequest(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    ctx = logger.WithTraceID(ctx, getTraceID(r))
    ctx = logger.WithUserID(ctx, getUserID(r))

    // 传递给业务逻辑
    processBusinessLogic(ctx, w, r)
}
```

### 3. 组件化日志

```go
// 为不同组件添加标识
logger.Info("Email analysis started",
    logger.String("component", "email_analyzer"),
    logger.String("email_id", emailID))

logger.Info("Contact updated",
    logger.String("component", "contact_manager"),
    logger.String("contact_id", contactID))
```

### 4. 错误处理

```go
// 使用 Error 字段
logger.ErrorContext(ctx, "Failed to process email",
    logger.Error(err),
    logger.String("email_id", emailID),
    logger.String("component", "email_processor"))
```

## 迁移指南

### 从 zap 迁移

```go
// 旧代码
logger.Infow("User logged in", "user_id", userID, "email", email)

// 新代码
logger.Info("User logged in",
    logger.String("user_id", userID),
    logger.String("email", email))

// 或使用上下文
logger.InfoContext(ctx, "User logged in",
    logger.String("email", email))
```

### 从 SugaredLogger 迁移

```go
// 旧代码
sugar.Infow("Task completed", "task_id", taskID, "duration", duration.Seconds())

// 新代码
logger.Info("Task completed",
    logger.String("task_id", taskID),
    logger.Float64("duration_seconds", duration.Seconds()))
```

## 性能优化

1. **启用采样**: 生产环境启用日志采样
2. **合理设置级别**: 避免在生产环境记录 DEBUG 日志
3. **批量输出**: 使用批量发送到企业日志平台
4. **异步写入**: 所有企业日志提供者使用异步写入

## 扩展

### 自定义提供者

```go
type CustomProvider struct {
    // 自定义字段
}

func (p *CustomProvider) Write(ctx context.Context, entry *logger.LogEntry) error {
    // 实现自定义写入逻辑
    return nil
}

func (p *CustomProvider) Close() error {
    // 实现清理逻辑
    return nil
}

func (p *CustomProvider) Ping() error {
    // 实现健康检查
    return nil
}
```

## 监控和告警

1. **日志量监控**: 监控日志输出量和错误率
2. **存储空间**: 监控日志文件和存储空间使用
3. **平台健康**: 定期检查企业日志平台连接状态

## 故障排查

1. **配置验证**: 确认配置文件和环境变量正确
2. **权限检查**: 确认日志目录和文件权限
3. **网络连接**: 检查企业日志平台网络连接
4. **采样配置**: 确认采样配置不会导致日志丢失