# EchoMind 日志框架最佳实践指南

## 📋 目录

1. [日志级别使用规范](#日志级别使用规范)
2. [字段命名和组织](#字段命名和组织)
3. [上下文管理](#上下文管理)
4. [性能优化](#性能优化)
5. [安全性考虑](#安全性考虑)
6. [监控和告警](#监控和告警)
7. [部署和配置](#部署和配置)
8. [故障排除](#故障排除)

## 🎯 日志级别使用规范

### 级别定义

```go
const (
    DebugLevel logger.Level = iota // 调试信息 - 仅用于开发环境
    InfoLevel                      // 一般信息 - 正常业务流程
    WarnLevel                      // 警告信息 - 潜在问题
    ErrorLevel                     // 错误信息 - 需要关注
    FatalLevel                     // 致命错误 - 系统无法继续
)
```

### 使用原则

#### ✅ 正确使用

```go
// DEBUG: 详细的技术信息，用于问题排查
logger.Debug("开始执行数据库查询",
    logger.String("sql", "SELECT * FROM users WHERE id = ?"),
    logger.String("table", "users"),
    logger.String("operation", "select"))

// INFO: 重要的业务事件
logger.Info("用户注册成功",
    logger.String("user_id", "usr_12345"),
    logger.String("email", "user@example.com"),
    logger.String("source", "web_signup"))

// WARN: 系统可以继续但需要注意的情况
logger.Warn("API 调用接近速率限制",
    logger.String("endpoint", "/api/v1/data"),
    logger.Int("current_rate", 95),
    logger.Int("limit", 100),
    logger.Time("reset_time", time.Now().Add(time.Hour)))

// ERROR: 错误情况但系统可以恢复
logger.Error("邮件发送失败",
    logger.Error(err),
    logger.String("template", "welcome"),
    logger.String("recipient", "user@example.com"),
    logger.Int("retry_count", 3))

// FATAL: 无法恢复的错误，系统需要退出
logger.Fatal("数据库连接初始化失败",
    logger.Error(err),
    logger.String("database", "production_db"),
    logger.String("host", "db.example.com"))
```

#### ❌ 错误使用

```go
// 不要在 INFO 中记录调试信息
logger.Info("SQL 查询执行") // 应该使用 Debug

// 不要在 ERROR 中记录一般信息
logger.Error("请求处理完成") // 应该使用 Info

// 不要过度使用 WARN
logger.Warn("用户操作成功") // 成功操作应该使用 Info
```

### 级别配置建议

```yaml
# 开发环境
level: DEBUG

# 测试环境
level: INFO

# 预生产环境
level: INFO

# 生产环境
level: WARN
```

## 🏷️ 字段命名和组织

### 命名规范

#### ✅ 推荐命名

```go
// ID 类字段
logger.String("user_id", "usr_12345")
logger.String("order_id", "ord_67890")
logger.String("session_id", "sess_abcde")
logger.String("trace_id", "trace_123456")

// 时间字段
logger.Time("created_at", time.Now())
logger.Time("updated_at", time.Now())
logger.Duration("processing_time", 150*time.Millisecond)
logger.Int("timestamp_ms", time.Now().UnixMilli())

// 状态字段
logger.String("status", "success|pending|failed|completed")
logger.String("state", "active|inactive|suspended")
logger.Bool("enabled", true)

// 计数字段
logger.Int("count", 42)
logger.Float64("amount", 19.99)
logger.String("percentage", "85.5%")

// 技术字段
logger.String("component", "email_service")
logger.String("version", "1.2.3")
logger.String("host", "server-01")
logger.Int("port", 8080)
```

#### ❌ 避免命名

```go
// 避免模糊的命名
logger.String("id", "123")          // 应该指明是什么 ID
logger.String("data", "some data")    // 应该说明数据类型
logger.String("value", "value")        // 应该说明值的含义
logger.String("result", "result")      // 应该说明结果类型

// 避免缩写和不一致
logger.String("usr_id", "123")        // 应该使用 user_id
logger.String("reqId", "123")         // 应该使用 request_id
logger.String("UID", "123")           // 应该使用 user_id
```

### 字段组织

#### 结构化字段组织

```go
// ✅ 推荐：按功能组织字段
logger.Info("订单处理完成",
    // 核心业务字段
    logger.String("order_id", "ORD-2024-001"),
    logger.String("customer_id", "CUST-12345"),
    logger.Float64("amount", 199.99),
    logger.String("currency", "USD"),

    // 状态字段
    logger.String("status", "completed"),
    logger.String("payment_status", "paid"),
    logger.String("shipping_status", "delivered"),

    // 时间字段
    logger.Time("order_time", orderTime),
    logger.Duration("processing_time", time.Since(orderTime)),

    // 技术字段
    logger.String("component", "order_processor"),
    logger.String("version", "2.1.0"),
    logger.String("host", "order-service-01")
)
```

#### 分层字段设计

```go
// 请求层字段
func logRequest(ctx context.Context, req *http.Request) {
    logger.InfoContext(ctx, "HTTP 请求",
        logger.String("method", req.Method),
        logger.String("path", req.URL.Path),
        logger.String("query", req.URL.RawQuery),
        logger.String("user_agent", req.UserAgent()),
        logger.String("remote_addr", req.RemoteAddr),
        logger.Int("content_length", req.ContentLength),
    )
}

// 业务层字段
func logBusinessEvent(ctx context.Context, event BusinessEvent) {
    logger.InfoContext(ctx, event.Name,
        logger.String("event_id", event.ID),
        logger.String("user_id", event.UserID),
        logger.String("entity_type", event.EntityType),
        logger.String("entity_id", event.EntityID),
        logger.String("action", event.Action),
        logger.Any("metadata", event.Metadata),
    )
}

// 系统层字段
func logSystemEvent(event string, details map[string]interface{}) {
    logger.Info("系统事件",
        logger.String("event", event),
        logger.String("component", "system_monitor"),
        logger.String("hostname", getHostname()),
        logger.String("version", getVersion()),
        logger.Any("details", details),
    )
}
```

## 🔄 上下文管理

### 上下文传播策略

#### 1. HTTP 请求上下文

```go
func requestMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 开始时间
        startTime := time.Now()

        // 生成请求 ID
        requestID := generateRequestID()
        r = r.WithContext(setRequestID(r.Context(), requestID))

        // 设置追踪 ID（如果从头部获取）
        if traceID := r.Header.Get("X-Trace-ID"); traceID != "" {
            r = r.WithContext(setTraceID(r.Context(), traceID))
        }

        // 获取用户信息（如果已认证）
        if userID := getUserIDFromToken(r.Header.Get("Authorization")); userID != "" {
            r = r.WithContext(setUserID(r.Context(), userID))
            r = r.WithContext(setOrgID(r.Context(), getUserOrg(userID)))
        }

        // 记录请求开始
        logger.InfoContext(r.Context(), "请求开始",
            logger.String("method", r.Method),
            logger.String("path", r.URL.Path))

        // 处理请求
        next.ServeHTTP(w, r)

        // 记录请求结束
        logger.InfoContext(r.Context(), "请求完成",
            logger.Int("status_code", w.(*responseWriter).statusCode),
            logger.Duration("duration", time.Since(startTime)))
    })
}
```

#### 2. 微服务调用上下文

```go
func callExternalService(ctx context.Context, serviceURL string) error {
    // 确保上下文包含追踪信息
    if traceID := logger.GetTraceID(ctx); traceID == "" {
        traceID = generateTraceID()
        ctx = logger.WithTraceID(ctx, traceID)
    }

    // 记录外部调用
    logger.InfoContext(ctx, "外部服务调用",
        logger.String("service_url", serviceURL),
        logger.String("service_name", extractServiceName(serviceURL)))

    startTime := time.Now()
    resp, err := http.DefaultClient.Do(buildRequest(ctx, serviceURL))
    duration := time.Since(startTime)

    if err != nil {
        logger.ErrorContext(ctx, "外部服务调用失败",
            logger.Error(err),
            logger.String("service_url", serviceURL),
            logger.Duration("duration", duration))
        return err
    }

    logger.InfoContext(ctx, "外部服务调用成功",
        logger.String("service_url", serviceURL),
        logger.Int("status_code", resp.StatusCode),
        logger.Duration("duration", duration))

    return nil
}
```

#### 3. 后台任务上下文

```go
func processBackgroundTask(ctx context.Context, taskID string) {
    // 为后台任务创建独立的上下文
    taskCtx := logger.WithTraceID(context.Background(), taskID)
    taskCtx = logger.WithContext(taskCtx,
        "", "", "", "", taskID)

    logger.InfoContext(taskCtx, "后台任务开始",
        logger.String("task_id", taskID),
        logger.String("task_type", "email_processing"))

    // 执行任务逻辑
    if err := processEmails(taskCtx); err != nil {
        logger.ErrorContext(taskCtx, "后台任务失败",
            logger.Error(err),
            logger.String("task_id", taskID))
        return
    }

    logger.InfoContext(taskCtx, "后台任务完成",
        logger.String("task_id", taskID))
}
```

### 上下文生命周期管理

```go
// ✅ 推荐：短生命周期上下文
func handleRequest(ctx context.Context, req *http.Request) {
    // 请求级别的上下文
    ctx = logger.WithTraceID(ctx, generateRequestID())

    // 业务逻辑使用上下文
    processBusinessLogic(ctx)
}

// ✅ 推荐：长生命周期上下文
func (s *Service) Start() error {
    // 服务级别的上下文
    s.ctx = logger.WithContext(context.Background(),
        "service-trace",
        "service-123",
        "org-456",
        "", // session 不适用于服务级别
        "")

    go s.runBackgroundTasks()
    return nil
}

// ❌ 避免：上下文泄漏
func badExample() {
    // 不要在全局变量中存储上下文
    globalCtx = logger.WithTraceID(context.Background(), "leak")

    // 上下文可能被意外修改
    someFunc(globalCtx)
}
```

## ⚡ 性能优化

### 1. 字段创建优化

```go
// ✅ 推荐：重用字段创建器
var (
    componentField = logger.String("component", "user_service")
    versionField   = logger.String("version", "1.2.3")
    hostField      = logger.String("host", getHostname())
)

func optimizedLogging(ctx context.Context, event string) {
    logger.InfoContext(ctx, event,
        componentField,
        versionField,
        hostField,
        logger.String("event_type", event),
        logger.Time("timestamp", time.Now()),
    )
}

// ✅ 推荐：批量字段创建
func createLogFields(event string) logger.Fields {
    return []logger.Field{
        logger.String("component", "user_service"),
        logger.String("version", "1.2.3"),
        logger.String("event", event),
        logger.Time("timestamp", time.Now()),
    }
}

func batchLogging(ctx context.Context, events []string) {
    baseFields := createLogFields("batch_operation")

    for _, event := range events {
        logger.InfoContext(ctx, "处理事件",
            append(baseFields, logger.String("specific_event", event))...)
    }
}
```

### 2. 条件日志记录

```go
// ✅ 推荐：条件性记录详细信息
func conditionalLogging(ctx context.Context, debug bool, data []byte) {
    if debug {
        logger.DebugContext(ctx, "处理详细数据",
            logger.String("data_size", fmt.Sprintf("%d bytes", len(data))),
            logger.String("data_hash", hashData(data)))
    }

    // 始终记录关键信息
    logger.InfoContext(ctx, "数据处理完成",
        logger.Int("data_size", len(data)))
}

// ✅ 推荐：使用采样避免日志洪水
func highFrequencyLogging(ctx context.Context, metric int) {
    // 使用采样减少高频日志
    if metric%100 == 0 { // 每 100 次记录一次
        logger.InfoContext(ctx, "批量处理统计",
            logger.Int("processed_count", 100),
            logger.Float64("average_metric", float64(metric)/100.0))
    }
}
```

### 3. 内存优化

```go
// ✅ 推荐：避免在热路径中创建大对象
func hotPathLogging(ctx context.Context, req *http.Request) {
    // 记录关键信息，避免记录大对象
    logger.InfoContext(ctx, "处理请求",
        logger.String("method", req.Method),
        logger.String("path", req.URL.Path),
        logger.Int("content_length", req.ContentLength))

    // 业务逻辑处理...
}

// ❌ 避免：在热路径中创建大对象
func badHotPathLogging(ctx context.Context, req *http.Request) {
    // 不要这样做：创建大对象用于日志记录
    requestDetails := map[string]interface{}{
        "headers": req.Headers,           // 可能很大
        "body":    string(req.Body),      // 可能很大
        "cookies": req.Cookies,           // 可能很多
    }

    logger.InfoContext(ctx, "处理请求",
        logger.Any("request_details", requestDetails)) // 大对象创建开销大
}
```

### 4. 异步日志记录

```go
// ✅ 推荐：使用异步批量处理器
func setupAsyncLogging() {
    config := &logger.Config{
        Providers: []logger.ProviderConfig{
            {
                Name:    "async_elasticsearch",
                Type:    "elasticsearch",
                Enabled: true,
                Settings: map[string]interface{}{
                    "url":        "http://elasticsearch:9200",
                    "index":      "echomind-logs",
                    "batch_size": 500,
                    "workers":    4,
                },
            },
        },
    }

    logger.Init(config)
}

// ✅ 推荐：后台任务使用专用日志器
func setupBackgroundLogger() {
    bgLogger := logger.With(
        logger.String("component", "background_worker"),
        logger.String("process_id", os.Getpid()),
        logger.String("hostname", getHostname()),
    )

    // 在后台任务中使用 bgLogger
    go func() {
        for {
            processBackgroundTask(bgLogger)
        }
    }()
}
```

## 🔒 安全性考虑

### 敏感信息处理

```go
// ✅ 推荐：脱敏处理
func sanitizeLogData(data map[string]interface{}) map[string]interface{} {
    sanitized := make(map[string]interface{})

    for key, value := range data {
        switch key {
        case "password", "token", "secret", "key":
            sanitized[key] = "***REDACTED***"
        case "email", "phone":
            sanitized[key] = maskPII(fmt.Sprintf("%v", value))
        case "credit_card":
            sanitized[key] = maskCreditCard(fmt.Sprintf("%v", value))
        default:
            sanitized[key] = value
        }
    }

    return sanitized
}

// PII 掩码函数
func maskPII(value string) string {
    if len(value) <= 4 {
        return strings.Repeat("*", len(value))
    }
    return value[:2] + strings.Repeat("*", len(value)-4) + value[len(value)-2:]
}

func maskCreditCard(card string) string {
    if len(card) <= 4 {
        return strings.Repeat("*", len(card))
    }
    return strings.Repeat("*", len(card)-4) + card[len(card)-4:]
}
```

### 访问控制

```go
// ✅ 推荐：基于角色的日志记录
func logByRole(ctx context.Context, event string, userRole string, data interface{}) {
    switch userRole {
    case "admin":
        // 管理员可以看到所有信息
        logger.InfoContext(ctx, event,
            logger.String("role", userRole),
            logger.Any("full_data", data))

    case "user":
        // 普通用户只能看到有限信息
        logger.InfoContext(ctx, event,
            logger.String("role", userRole),
            logger.String("summary", summarizeData(data)))

    default:
        // 其他角色只记录事件发生
        logger.InfoContext(ctx, event, logger.String("role", userRole))
    }
}

// 审计日志记录
func auditLog(ctx context.Context, action string, resource string, result string) {
    // 审计日志必须包含足够的信息用于追踪
    auditLogger := logger.With(
        logger.String("log_type", "audit"),
        logger.String("compliance", "SOX"),
    )

    auditLogger.InfoContext(ctx, "审计事件",
        logger.String("action", action),
        logger.String("resource", resource),
        logger.String("result", result),
        logger.Time("audit_timestamp", time.Now()),
        logger.String("trace_id", logger.GetTraceID(ctx)),
        logger.String("user_id", logger.GetUserID(ctx)),
        logger.String("org_id", logger.GetOrgID(ctx)))
}
```

### 加密日志

```go
// ✅ 推荐：对敏感日志字段加密
func encryptSensitiveField(value string) (string, error) {
    // 使用加密算法加密敏感字段
    encrypted, err := encryption.Encrypt(value)
    if err != nil {
        return "", err
    }
    return "ENC:" + encrypted, nil
}

func logWithEncryption(ctx context.Context, sensitiveData string) {
    // 加密敏感数据
    encryptedData, err := encryptSensitiveField(sensitiveData)
    if err != nil {
        logger.ErrorContext(ctx, "加密敏感数据失败", logger.Error(err))
        return
    }

    logger.InfoContext(ctx, "处理敏感数据",
        logger.String("encrypted_data", encryptedData),
        logger.String("data_type", "sensitive"))
}
```

## 📊 监控和告警

### 关键指标监控

```go
// 日志指标监控
func setupLogMetrics() {
    go func() {
        ticker := time.NewTicker(1 * time.Minute)
        defer ticker.Stop()

        for range ticker.C {
            monitorLogLevels()
            monitorProviderHealth()
            monitorLogVolume()
        }
    }()
}

func monitorLogLevels() {
    // 检查不同级别日志的比例
    totalLogs := getTotalLogCount()
    errorLogs := getErrorLogCount()

    if totalLogs > 0 {
        errorRate := float64(errorLogs) / float64(totalLogs) * 100

        if errorRate > 10.0 { // 错误率超过 10%
            sendAlert("日志错误率过高", map[string]interface{}{
                "error_rate": errorRate,
                "threshold": "10%",
                "error_count": errorLogs,
                "total_count": totalLogs,
            })
        }
    }
}

func monitorProviderHealth() {
    // 检查日志提供者健康状态
    for _, provider := range getLogProviders() {
        if err := provider.Ping(); err != nil {
            sendAlert("日志提供者连接失败", map[string]interface{}{
                "provider": provider.GetType(),
                "error": err.Error(),
            })
        }
    }
}

func monitorLogVolume() {
    // 监控日志量变化
    currentVolume := getCurrentLogVolume()
    previousVolume := getPreviousLogVolume()

    if previousVolume > 0 {
        changeRate := float64(currentVolume-previousVolume) / float64(previousVolume) * 100

        // 日志量异常变化
        if changeRate > 200 || changeRate < -50 {
            sendAlert("日志量异常变化", map[string]interface{}{
                "change_rate": changeRate,
                "current_volume": currentVolume,
                "previous_volume": previousVolume,
            })
        }
    }
}
```

### 告警配置

```go
// 告警规则配置
type AlertRule struct {
    Name        string
    Condition   func() bool
    Severity    string
    Description string
    Cooldown    time.Duration
}

var alertRules = []AlertRule{
    {
        Name: "高错误率",
        Condition: func() bool {
            return getErrorRate() > 5.0
        },
        Severity:    "high",
        Description: "系统错误率超过 5%",
        Cooldown:    5 * time.Minute,
    },
    {
        Name: "日志提供者离线",
        Condition: func() bool {
            return !areAllProvidersHealthy()
        },
        Severity:    "critical",
        Description: "一个或多个日志提供者不可用",
        Cooldown:    1 * time.Minute,
    },
    {
        Name: "日志缓冲区满",
        Condition: func() bool {
            return getBufferUsage() > 0.9
        },
        Severity:    "medium",
        Description: "日志缓冲区使用率超过 90%",
        Cooldown:    2 * time.Minute,
    },
}

func checkAlerts() {
    for _, rule := range alertRules {
        if rule.Condition() {
            sendAlert(rule.Name, map[string]interface{}{
                "severity":    rule.Severity,
                "description": rule.Description,
            })

            // 实施冷却期
            time.Sleep(rule.Cooldown)
        }
    }
}
```

### 自动化响应

```go
// 自动化故障响应
func autoResponseAlert(alertName string, details map[string]interface{}) {
    switch alertName {
    case "高错误率":
        // 自动启用详细日志
        logger.SetLevel(logger.DebugLevel)
        logger.WarnContext(context.Background(), "自动启用详细日志记录")

        // 5 分钟后恢复到正常级别
        time.AfterFunc(5*time.Minute, func() {
            logger.SetLevel(logger.InfoLevel)
            logger.Info("恢复到正常日志级别")
        })

    case "日志提供者离线":
        // 切换到备用提供者
        switchToBackupProviders()

        // 降低日志级别以减少影响
        logger.SetLevel(logger.ErrorLevel)
        logger.Error("切换到备用日志提供者，降低日志级别")

    case "日志缓冲区满":
        // 强制刷新缓冲区
        forceFlushLogBuffers()

        // 临时增大缓冲区
        increaseLogBufferSize()
    }
}

func switchToBackupProviders() {
    // 实现备用提供者切换逻辑
    // 1. 停用失败的提供者
    // 2. 启用文件输出作为备用
    // 3. 发送通知给运维团队
}

func forceFlushLogBuffers() {
    // 强制刷新所有日志缓冲区
    logger.Close()

    // 重新初始化日志系统
    logger.Init(logger.GetConfig())
}

func increaseLogBufferSize() {
    // 动态增加缓冲区大小
    // 实现缓冲区扩容逻辑
}
```

## 🚀 部署和配置

### 环境配置

#### 开发环境配置

```yaml
# config/logger-dev.yaml
level: DEBUG
production: false

output:
  console:
    enabled: true
    format: console
    color: true
  file:
    enabled: false

context:
  auto_fields: ["trace_id", "user_id", "session_id"]
  global_fields:
    service: echomind
    environment: development
    version: 1.2.3

sampling:
  enabled: false

providers:
  - name: "local_file"
    type: "noop"
    enabled: true
```

#### 生产环境配置

```yaml
# config/logger-prod.yaml
level: WARN
production: true

output:
  console:
    enabled: false
  file:
    enabled: true
    path: "/var/log/echomind/app.log"
    max_size: 500
    max_age: 30
    max_backups: 10
    compress: true

context:
  auto_fields: ["trace_id", "user_id", "org_id", "session_id", "request_id"]
  global_fields:
    service: echomind
    environment: production
    version: 1.2.3
    cluster: production

sampling:
  enabled: true
  rate: 1000
  burst: 100
  levels: [DEBUG, INFO]

providers:
  - name: "elasticsearch"
    type: "elasticsearch"
    enabled: true
    settings:
      url: "http://elasticsearch.internal:9200"
      index: "echomind-logs"
      batch_size: 500
      username: "elastic"
      password: "${ELASTIC_PASSWORD}"
      timeout: "10s"

  - name: "loki"
    type: "loki"
    enabled: false
    settings:
      url: "http://loki.internal:3100/loki/api/v1/push"
      labels:
        service: echomind
        environment: production
```

### Kubernetes 部署配置

```yaml
# k8s/logger-configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: logger-config
data:
  config.yaml: |
    level: INFO
    production: true
    output:
      console:
        enabled: false
      file:
        enabled: true
        path: "/var/log/app/app.log"
        max_size: 500
        max_age: 30
        max_backups: 10
    providers:
      - name: "elasticsearch"
        type: "elasticsearch"
        enabled: true
        settings:
          url: "http://elasticsearch:9200"
          index: "echomind-logs"
          batch_size: 200
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: echomind-backend
spec:
  template:
    spec:
      containers:
      - name: echomind
        env:
        - name: LOG_LEVEL
          value: "INFO"
        - name: LOG_PRODUCTION
          value: "true"
        - name: ELASTIC_PASSWORD
          valueFrom:
            secretKeyRef:
              name: elasticsearch-secret
              key: password
        volumeMounts:
        - name: log-volume
          mountPath: /var/log/app
      volumes:
      - name: log-volume
        emptyDir: {}
      - name: logger-config
        configMap:
          name: logger-config
          mountPath: /app/config/logger
```

### Docker 配置

```dockerfile
# Dockerfile
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY backend/ .
RUN go mod tidy && go build -o /bin/echomind ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /bin/echomind /bin/
COPY configs/ /app/configs/

# 创建日志目录
RUN mkdir -p /var/log/echomind

# 设置日志配置环境变量
ENV LOG_LEVEL=INFO
ENV LOG_PRODUCTION=true
ENV LOG_FILE_PATH=/var/log/echomind/app.log

CMD ["/bin/echomind"]
```

### Docker Compose 配置

```yaml
# docker-compose.yml
version: '3.8'

services:
  echomind:
    build:
      context: .
      dockerfile: Dockerfile
    environment:
      - LOG_LEVEL=INFO
      - LOG_PRODUCTION=true
      - ELASTIC_PASSWORD=${ELASTIC_PASSWORD}
    volumes:
      - ./logs:/var/log/echomind
      - ./configs:/app/configs
    depends_on:
      - elasticsearch
      - redis

  elasticsearch:
    image: elasticsearch:8.11.0
    environment:
      - discovery.type=single-node
      - "ES_JAVA_OPTS=-Xms512m -Xmx512m"
    ports:
      - "9200:9200"
    volumes:
      - elasticsearch_data:/usr/share/elasticsearch/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

volumes:
  elasticsearch_data:
  redis_data:
```

## 🔧 故障排除

### 常见问题和解决方案

#### 1. 日志级别问题

```go
// 检查当前日志级别
func checkLogLevel() {
    currentLevel := logger.GetLevel()
    fmt.Printf("当前日志级别: %s\n", currentLevel.String())

    // 测试不同级别的日志
    logger.Debug("这是 DEBUG 级别日志")
    logger.Info("这是 INFO 级别日志")
    logger.Warn("这是 WARN 级别日志")
    logger.LogError("这是 ERROR 级别日志")
}

// 动态调整日志级别
func adjustLogLevel(newLevel string) {
    var level logger.Level
    switch newLevel {
    case "DEBUG":
        level = logger.DebugLevel
    case "INFO":
        level = logger.InfoLevel
    case "WARN":
        level = logger.WarnLevel
    case "ERROR":
        level = logger.ErrorLevel
    default:
        fmt.Printf("无效的日志级别: %s\n", newLevel)
        return
    }

    logger.SetLevel(level)
    fmt.Printf("日志级别已调整为: %s\n", level.String())
}
```

#### 2. 提供者连接问题

```go
// 诊断提供者连接
func diagnoseProviders() {
    config := logger.GetConfig()

    for _, providerConfig := range config.Providers {
        if !providerConfig.Enabled {
            fmt.Printf("提供者 %s 已禁用\n", providerConfig.Name)
            continue
        }

        fmt.Printf("检查提供者: %s (类型: %s)\n", providerConfig.Name, providerConfig.Type)

        provider, err := logger.createProvider(providerConfig)
        if err != nil {
            fmt.Printf("创建提供者失败: %v\n", err)
            continue
        }

        // 测试连接
        if err := provider.Ping(); err != nil {
            fmt.Printf("提供者连接失败: %v\n", err)

            // 提供解决建议
            switch providerConfig.Type {
            case "elasticsearch":
                fmt.Println("建议检查:")
                fmt.Println("  - Elasticsearch 服务是否运行")
                fmt.Println("  - 网络连接是否正常")
                fmt.Println("  - 认证信息是否正确")
            case "loki":
                fmt.Println("建议检查:")
                fmt.Println("  - Loki 服务是否运行")
                fmt.Println("  - 端口配置是否正确")
            case "splunk":
                fmt.Println("建议检查:")
                fmt.Println("  - Splunk HEC 端口是否开放")
                fmt.Println("  - HEC Token 是否有效")
            }
        } else {
            fmt.Printf("提供者连接正常\n")
        }

        // 关闭提供者
        provider.Close()
    }
}
```

#### 3. 性能问题诊断

```go
// 性能诊断工具
func diagnosePerformance() {
    // 运行性能基准测试
    fmt.Println("运行性能基准测试...")

    // 基础日志记录性能
    start := time.Now()
    for i := 0; i < 10000; i++ {
        logger.Info("性能测试",
            logger.Int("iteration", i),
            logger.String("data", strings.Repeat("x", 100)))
    }
    duration := time.Since(start)

    avgPerLog := duration / 10000
    logsPerSecond := float64(time.Second) / float64(avgPerLog)

    fmt.Printf("性能测试结果:\n")
    fmt.Printf("  总耗时: %v\n", duration)
    fmt.Printf("  平均每条日志: %v\n", avgPerLog)
    fmt.Printf("  每秒日志数: %.0f\n", logsPerSecond)

    // 提供性能优化建议
    if avgPerLog > 1*time.Millisecond {
        fmt.Println("性能建议:")
        fmt.Println("  - 考虑启用日志采样")
        fmt.Println("  - 增加批量处理大小")
        fmt.Println("  - 使用异步日志记录")
        fmt.Println("  - 减少日志字段数量")
    }

    // 检查内存使用
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    fmt.Printf("\n内存使用情况:\n")
    fmt.Printf("  分配内存: %d KB\n", m.Alloc/1024)
    fmt.Printf("  系统内存: %d KB\n", m.Sys/1024)
    fmt.Printf("  GC 次数: %d\n", m.NumGC)
}

// 缓冲区诊断
func diagnoseBuffer() {
    stats := logger.GetStats().(map[string]interface{})

    if bufferUsage, exists := stats["buffer_usage"]; exists {
        usage := bufferUsage.(float64)
        fmt.Printf("缓冲区使用率: %.1f%%\n", usage*100)

        if usage > 0.8 {
            fmt.Println("缓冲区使用率过高，建议:")
            fmt.Println("  - 增加批量刷新频率")
            fmt.Println("  - 启用日志采样")
            fmt.Println("  - 检查日志提供者性能")
        }
    }
}
```

### 日志分析和查询

```go
// 生成日志分析报告
func generateLogAnalysisReport() {
    report := map[string]interface{}{
        "timestamp": time.Now(),
        "period": "last_24_hours",
        "statistics": collectLogStatistics(),
        "trends": analyzeLogTrends(),
        "alerts": identifyLogAnomalies(),
        "recommendations": generateRecommendations(),
    }

    // 将报告输出到日志
    logger.Info("日志分析报告", logger.Any("report", report))
}

// 收集日志统计
func collectLogStatistics() map[string]interface{} {
    return map[string]interface{}{
        "total_logs":     getTotalLogCount(),
        "error_logs":     getErrorLogCount(),
        "warning_logs":   getWarningLogCount(),
        "info_logs":      getInfoLogCount(),
        "debug_logs":     getDebugLogCount(),
        "unique_users":   getUniqueUserCount(),
        "unique_traces":  getUniqueTraceCount(),
        "error_rate":     calculateErrorRate(),
    }
}

// 分析日志趋势
func analyzeLogTrends() map[string]interface{} {
    hourlyStats := getHourlyLogStats()

    return map[string]interface{}{
        "peak_hour":     findPeakHour(hourlyStats),
        "growth_rate":   calculateGrowthRate(hourlyStats),
        "patterns":      identifyPatterns(hourlyStats),
    }
}
```

## 📚 参考资料

- [EchoMind 日志框架 API 文档](./README.md)
- [EchoMind 日志框架使用示例](./EXAMPLES.md)
- [Elasticsearch 日志最佳实践](https://www.elastic.co/guide/en/elasticsearch/guide/current/logging.html)
- [Grafana Loki 日志聚合](https://grafana.com/docs/loki/)
- [Splunk 日志管理](https://docs.splunk.com/Documentation/Splunk/8.0.0/Data/GetStarted/DataLog/)
- [Go 日志最佳实践](https://github.com/golang/go/wiki/CodeReviewComments#logging)