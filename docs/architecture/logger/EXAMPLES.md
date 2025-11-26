# EchoMind 日志框架使用示例

本文档提供了 EchoMind 企业级日志框架的详细使用示例和最佳实践。

## 🚀 快速开始

### 基础初始化

```go
package main

import (
    "github.com/hrygo/echomind/pkg/logger"
)

func main() {
    // 1. 使用默认配置
    if err := logger.Init(logger.DefaultConfig()); err != nil {
        panic(err)
    }

    // 2. 使用生产环境配置
    if err := logger.Init(logger.ProductionConfig()); err != nil {
        panic(err)
    }

    // 3. 使用开发环境配置
    if err := logger.Init(logger.DevelopmentConfig()); err != nil {
        panic(err)
    }

    // 4. 从环境变量加载配置
    config := logger.LoadConfigFromEnv()
    if err := logger.Init(config); err != nil {
        panic(err)
    }

    // 记录日志
    logger.Info("应用启动成功")

    // 清理资源
    defer logger.Close()
}
```

## 📝 基础日志记录

### 简单日志

```go
// 基础日志级别
logger.Debug("调试信息")
logger.Info("一般信息")
logger.Warn("警告信息")
logger.LogError("错误信息") // 使用 LogError 避免与 logger.Error 冲突
logger.Fatal("致命错误") // 这会调用 os.Exit(1)
```

### 结构化日志

```go
logger.Info("用户登录成功",
    logger.String("user_id", "12345"),
    logger.String("email", "user@example.com"),
    logger.String("ip", "192.168.1.100"),
    logger.Duration("response_time", 150*time.Millisecond),
    logger.Bool("new_user", false),
)

logger.Error("数据库连接失败",
    logger.Error(err),
    logger.String("database", "postgres"),
    logger.Int("retry_count", 3),
    logger.Time("last_attempt", time.Now()),
)
```

### 使用字段创建函数

```go
// 字段类型
logger.String("key", "value")
logger.Int("count", 42)
logger.Float64("price", 19.99)
logger.Bool("enabled", true)
logger.Duration("latency", 200*time.Millisecond)
logger.Time("created_at", time.Now())
logger.Any("data", complexObject)
logger.Error(err)

// 从 map 创建字段
fields := logger.FieldsFromMap(map[string]interface{}{
    "user_id":  "12345",
    "action":   "login",
    "success":  true,
})

logger.Info("操作完成", fields...)
```

## 🔄 上下文感知日志

### 上下文传递

```go
import "github.com/hrygo/echomind/pkg/logger"

func handleRequest(r *http.Request) {
    // 创建上下文
    ctx := r.Context()

    // 添加追踪信息
    traceID := generateTraceID()
    ctx = logger.WithTraceID(ctx, traceID)

    // 添加用户信息
    if userID := getUserID(r); userID != "" {
        ctx = logger.WithUserID(ctx, userID)
    }

    // 添加组织信息
    if orgID := getOrgID(r); orgID != "" {
        ctx = logger.WithOrgID(ctx, orgID)
    }

    // 添加会话信息
    if sessionID := getSessionID(r); sessionID != "" {
        ctx = logger.WithSessionID(ctx, sessionID)
    }

    // 使用上下文记录日志
    processRequest(ctx, r)
}

func processRequest(ctx context.Context, r *http.Request) {
    logger.InfoContext(ctx, "开始处理请求",
        logger.String("path", r.URL.Path),
        logger.String("method", r.Method))

    // 业务逻辑处理...

    logger.InfoContext(ctx, "请求处理完成",
        logger.Duration("processing_time", time.Since(startTime)))
}
```

### 一次设置多个上下文值

```go
// 一次性设置多个上下文值
ctx := logger.WithContext(
    context.Background(),
    "trace-123",
    "user-456",
    "org-789",
    "session-012",
    "req-345",
)

logger.InfoContext(ctx, "所有上下文信息已设置")
```

### With 方法创建专用日志器

```go
// 为特定组件创建带固定字段的日志器
emailLogger := logger.With(
    logger.String("component", "email_processor"),
    logger.String("version", "1.0.0"),
)

// 使用专用日志器
emailLogger.Info("开始处理邮件",
    logger.String("email_id", "abc-123"))

emailLogger.Error("邮件处理失败",
    logger.Error(err),
    logger.String("email_id", "abc-123"))
```

## 🔧 高级配置

### 自定义配置

```go
config := &logger.Config{
    Level:      logger.InfoLevel,
    Production: true,
    Output: logger.OutputConfig{
        File: logger.FileOutputConfig{
            Enabled:   true,
            Path:      "/var/log/echomind/app.log",
            MaxSize:   500, // MB
            MaxAge:    30,  // days
            MaxBackups: 10,
            Compress:  true,
        },
        Console: logger.ConsoleOutputConfig{
            Enabled: false,
            Format:  "json",
            Color:   false,
        },
    },
    Context: logger.ContextConfig{
        AutoFields: []string{
            "trace_id", "user_id", "org_id", "session_id", "request_id",
        },
        GlobalFields: map[string]interface{}{
            "service":     "echomind-backend",
            "environment": "production",
            "version":     "1.0.0",
        },
    },
    Sampling: logger.SamplingConfig{
        Enabled: true,
        Rate:    1000, // 每秒最多 1000 条日志
        Burst:   100,  // 突发最多 100 条
        Levels:  []logger.Level{logger.DebugLevel, logger.InfoLevel},
    },
    Providers: []logger.ProviderConfig{
        {
            Name:    "elasticsearch",
            Type:    "elasticsearch",
            Enabled: true,
            Settings: map[string]interface{}{
                "url":        "http://elasticsearch:9200",
                "index":      "echomind-logs",
                "batch_size": 500,
                "username":   "elastic",
                "password":   "changeme",
            },
        },
        {
            Name:    "loki",
            Type:    "loki",
            Enabled: false, // 暂时禁用
            Settings: map[string]interface{}{
                "url": "http://loki:3100/loki/api/v1/push",
                "labels": map[string]string{
                    "service": "echomind",
                },
            },
        },
    },
}

if err := logger.Init(config); err != nil {
    panic(err)
}
```

### 企业日志平台集成

#### Elasticsearch 集成

```go
config := &logger.Config{
    Providers: []logger.ProviderConfig{
        {
            Name:    "elasticsearch",
            Type:    "elasticsearch",
            Enabled: true,
            Settings: map[string]interface{}{
                "url":        "http://elasticsearch:9200",
                "index":      "echomind-logs",
                "batch_size": 200,
                "timeout":    "10s",
                "username":   "elastic",
                "password":   "your_password",
            },
        },
    },
}

// 或使用优化版提供者
esConfig := map[string]interface{}{
    "url":        "http://elasticsearch:9200",
    "index":      "echomind-logs",
    "batch_size": 500,
    "username":   "elastic",
    "password":   "your_password",
}

provider, err := logger.NewOptimizedElasticsearchProvider(esConfig)
if err != nil {
    panic(err)
}

// 手动写入到提供者
entry := &logger.LogEntry{
    Timestamp: time.Now(),
    Level:     logger.InfoLevel,
    Message:   "直接写入 Elasticsearch",
    Fields:    map[string]interface{}{"source": "manual"},
}

if err := provider.Write(context.Background(), entry); err != nil {
    // 处理错误
}
```

#### Grafana Loki 集成

```go
config := &logger.Config{
    Providers: []logger.ProviderConfig{
        {
            Name:    "loki",
            Type:    "loki",
            Enabled: true,
            Settings: map[string]interface{}{
                "url": "http://loki:3100/loki/api/v1/push",
                "labels": map[string]string{
                    "service":     "echomind",
                    "environment": "production",
                    "version":     "1.0.0",
                },
                "batch_size": 100,
            },
        },
    },
}
```

#### Splunk 集成

```go
config := &logger.Config{
    Providers: []logger.ProviderConfig{
        {
            Name:    "splunk",
            Type:    "splunk",
            Enabled: true,
            Settings: map[string]interface{}{
                "url":     "https://splunk.example.com:8088/services/collector/event",
                "token":   "your_hec_token",
                "index":   "echomind",
                "source":  "backend",
                "sourcetype": "json",
            },
        },
    },
}
```

## 🎯 最佳实践

### 1. 日志级别使用指南

```go
// ✅ 推荐：合适的日志级别
logger.Debug("开始执行 SQL 查询", logger.String("sql", sql))        // 调试信息
logger.Info("用户登录成功", logger.String("user_id", userID))         // 一般信息
logger.Warn("API 调用接近速率限制", logger.Int("current_rate", 95))    // 警告信息
logger.Error("数据库连接失败", logger.Error(err))                    // 错误信息

// ❌ 不推荐：级别滥用
logger.Info("调试信息")                                           // 应该使用 Debug
logger.Error("一般信息")                                           // 应该使用 Info
```

### 2. 字段命名规范

```go
// ✅ 推荐：清晰的字段命名
logger.Info("处理订单",
    logger.String("order_id", "ORD-2024-001"),
    logger.String("customer_id", "CUST-12345"),
    logger.Float64("amount", 99.99),
    logger.String("currency", "USD"),
    logger.Time("order_time", time.Now()),
    logger.String("status", "pending"),
)

// ❌ 不推荐：模糊的字段命名
logger.Info("处理订单",
    logger.String("id", "ORD-2024-001"),     // 不清楚是什么 ID
    logger.String("value", "99.99"),        // 不清楚是什么值
    logger.String("time", "2024-01-01"),     // 不清楚是什么时间
)
```

### 3. 上下文传播

```go
// ✅ 推荐：在请求开始时设置上下文
func middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        traceID := generateTraceID()
        ctx := logger.WithTraceID(r.Context(), traceID)

        // 如果有用户信息
        if userID := getUserIDFromToken(r.Header.Get("Authorization")); userID != "" {
            ctx = logger.WithUserID(ctx, userID)
        }

        // 传递到下一个处理器
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// 在业务逻辑中使用上下文
func processOrder(ctx context.Context, orderID string) error {
    logger.InfoContext(ctx, "开始处理订单",
        logger.String("order_id", orderID))

    // 业务逻辑...

    logger.InfoContext(ctx, "订单处理完成")
    return nil
}
```

### 4. 错误处理

```go
// ✅ 推荐：结构化错误信息
func (s *Service) ProcessData(ctx context.Context, data []byte) error {
    if err := validateData(data); err != nil {
        logger.ErrorContext(ctx, "数据验证失败",
            logger.Error(err),
            logger.String("data_hash", hashData(data)),
            logger.Int("data_size", len(data)),
            logger.String("component", "data_processor"),
        )
        return fmt.Errorf("数据验证失败: %w", err)
    }

    // 处理逻辑...
    return nil
}

// ❌ 不推荐：简单错误日志
func (s *Service) ProcessData(ctx context.Context, data []byte) error {
    if err := validateData(data); err != nil {
        logger.Error("验证失败") // 缺少上下文和详细信息
        return err
    }

    // 处理逻辑...
    return nil
}
```

### 5. 性能优化

```go
// ✅ 推荐：批量处理和延迟序列化
func processBatch(ctx context.Context, items []Item) {
    // 创建带组件字段的日志器
    itemLogger := logger.With(
        logger.String("component", "batch_processor"),
        logger.String("batch_id", generateBatchID()),
    )

    for i, item := range items {
        itemLogger.InfoContext(ctx, "处理项目",
            logger.Int("index", i),
            logger.String("item_id", item.ID),
            logger.String("item_type", item.Type))

        // 处理项目...
    }

    itemLogger.InfoContext(ctx, "批量处理完成",
        logger.Int("total_items", len(items)),
        logger.Duration("processing_time", time.Since(startTime)))
}

// ✅ 推荐：避免在热路径中创建大对象
func handleRequest(ctx context.Context, req *Request) {
    // 避免创建大的结构体用于日志
    logger.InfoContext(ctx, "处理请求",
        logger.String("method", req.Method),
        logger.String("path", req.Path),
        logger.Int("content_length", len(req.Body)),
        // 不要记录整个请求体
    )
}

// ❌ 不推荐：在热路径中创建大对象
func handleRequest(ctx context.Context, req *Request) {
    // 不要这样做：创建大对象用于日志记录
    requestDetails := map[string]interface{}{
        "method":      req.Method,
        "path":        req.Path,
        "headers":     req.Headers, // 可能很大
        "body":        string(req.Body), // 可能很大
        "query_params": req.QueryParams,
    }

    logger.InfoContext(ctx, "处理请求", logger.Any("request", requestDetails))
}
```

### 6. 敏感信息处理

```go
// ✅ 推荐：避免记录敏感信息
func processPayment(ctx context.Context, payment *Payment) error {
    logger.InfoContext(ctx, "处理支付",
        logger.String("payment_id", payment.ID),
        logger.String("currency", payment.Currency),
        logger.Float64("amount", payment.Amount),
        logger.String("merchant_id", payment.MerchantID),
        // 不要记录卡号、CVV 等敏感信息
    )

    if err := validatePayment(payment); err != nil {
        logger.ErrorContext(ctx, "支付验证失败",
            logger.Error(err),
            logger.String("payment_id", payment.ID),
            // 不要暴露具体的验证失败原因，可能包含敏感信息
        )
        return err
    }

    return nil
}

// ❌ 不推荐：记录敏感信息
func processPayment(ctx context.Context, payment *Payment) error {
    logger.InfoContext(ctx, "处理支付",
        logger.String("card_number", payment.CardNumber),  // 敏感信息！
        logger.String("cvv", payment.CVV),               // 敏感信息！
        logger.String("card_holder", payment.CardHolder), // 可能的敏感信息
    )

    // 处理逻辑...
    return nil
}
```

## 🔍 监控和告警

### 关键指标监控

```go
// 在应用启动时设置监控指标
func setupLoggingMetrics() {
    go func() {
        ticker := time.NewTicker(1 * time.Minute)
        defer ticker.Stop()

        for range ticker.C {
            // 获取日志统计信息
            if stats, ok := logger.GetStats().(map[string]interface{}); ok {
                // 记录到监控系统
                recordMetrics("logging", stats)

                // 检查告警条件
                if bufferUsage, exists := stats["buffer_usage"]; exists {
                    if usage, ok := bufferUsage.(float64); ok && usage > 0.8 {
                        // 日志缓冲区使用率过高，发送告警
                        sendAlert("日志缓冲区使用率过高", map[string]interface{}{
                            "usage": usage,
                            "threshold": 0.8,
                        })
                    }
                }
            }
        }
    }()
}
```

### 错误监控

```go
// 创建错误监控日志器
errorLogger := logger.With(
    logger.String("component", "error_monitor"),
    logger.String("alert_level", "high"),
)

func monitorErrors() {
    // 监控错误率
    errorCount := 0
    totalRequests := 0

    go func() {
        ticker := time.NewTicker(5 * time.Minute)
        defer ticker.Stop()

        for range ticker.C {
            if totalRequests > 0 {
                errorRate := float64(errorCount) / float64(totalRequests) * 100

                if errorRate > 5.0 { // 错误率超过 5%
                    errorLogger.Error("系统错误率过高",
                        logger.Float64("error_rate", errorRate),
                        logger.Int("error_count", errorCount),
                        logger.Int("total_requests", totalRequests),
                        logger.String("threshold", "5%"))
                }
            }

            // 重置计数器
            errorCount = 0
            totalRequests = 0
        }
    }()
}
```

## 🛠️ 故障排除

### 常见问题和解决方案

#### 1. 日志没有输出

```go
// 检查日志级别设置
currentLevel := logger.GetLevel()
fmt.Printf("当前日志级别: %s\n", currentLevel.String())

// 测试不同级别的日志
logger.Debug("这是调试信息 - 只有在 Debug 级别才会显示")
logger.Info("这是信息 - 应该始终显示")
logger.Error("这是错误 - 应该始终显示")
```

#### 2. 企业日志平台连接失败

```go
// 测试提供者连接
config := &logger.Config{
    Providers: []logger.ProviderConfig{
        {
            Name:    "test-elasticsearch",
            Type:    "elasticsearch",
            Enabled: true,
            Settings: map[string]interface{}{
                "url": "http://localhost:9200",
            },
        },
    },
}

log, err := logger.NewLogger(config)
if err != nil {
    fmt.Printf("创建日志器失败: %v\n", err)
}

// 检查提供者健康状态
for _, provider := range log.GetConfig().Providers {
    if provider.Enabled {
        if p, err := logger.createProvider(provider); err == nil {
            if err := p.Ping(); err != nil {
                fmt.Printf("提供者 %s 健康检查失败: %v\n", provider.Name, err)
            } else {
                fmt.Printf("提供者 %s 连接正常\n", provider.Name)
            }
        }
    }
}
```

#### 3. 性能问题

```go
// 检查性能指标
func checkPerformance() {
    start := time.Now()

    // 记录大量日志进行性能测试
    for i := 0; i < 1000; i++ {
        logger.Info("性能测试",
            logger.Int("iteration", i),
            logger.String("data", strings.Repeat("x", 100)))
    }

    duration := time.Since(start)
    fmt.Printf("记录 1000 条日志耗时: %v\n", duration)
    fmt.Printf("平均每条日志耗时: %v\n", duration/1000)
}
```

## 📚 参考资源

- [EchoMind 日志框架 API 文档](./README.md)
- [企业日志最佳实践指南](./BEST_PRACTICES.md)
- [性能优化建议](./PERFORMANCE_OPTIMIZATION.md)
- [故障排除指南](./TROUBLESHOOTING.md)