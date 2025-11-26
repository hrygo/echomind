# OpenTelemetry 集成指南

**版本**: v1.2.0  
**更新日期**: 2025-11-26

---

## 📖 目录

1. [概述](#概述)
2. [快速开始](#快速开始)
3. [配置说明](#配置说明)
4. [使用指南](#使用指南)
5. [指标说明](#指标说明)
6. [最佳实践](#最佳实践)
7. [故障排查](#故障排查)

---

## 概述

EchoMind v1.2.0 集成了 OpenTelemetry 可观测性框架,提供:

- **Distributed Tracing** (分布式追踪) - 请求链路追踪
- **Metrics** (指标) - 性能和业务指标收集
- **Logs Correlation** (日志关联) - 日志与追踪关联

### 架构

```
┌─────────────┐
│   Backend   │
│  Services   │
└──────┬──────┘
       │
       ├─→ Traces (OpenTelemetry SDK)
       ├─→ Metrics (OpenTelemetry SDK)
       └─→ Logs (Zap + TraceID)
              │
              ↓
       ┌─────────────┐
       │   Exporter  │
       │  (Console/  │
       │   File/     │
       │   OTLP)     │
       └─────────────┘
```

---

## 快速开始

### 1. 启用 OpenTelemetry

编辑 `backend/configs/config.yaml`:

```yaml
telemetry:
  enabled: true  # 开启可观测性
  service_name: "echomind-backend"
  service_version: "v1.2.0"
  environment: "development"
  
  exporter:
    type: "console"  # 开发环境使用控制台输出
```

### 2. 启动应用

```bash
cd backend
go run cmd/main.go
```

### 3. 验证

执行搜索请求后,控制台会输出 trace 和 metrics:

```json
{
  "Name": "SearchService.Search",
  "SpanContext": {
    "TraceID": "...",
    "SpanID": "..."
  },
  "Attributes": {
    "user.id": "...",
    "search.query": "project",
    "results.total": 5
  }
}
```

---

## 配置说明

### 完整配置示例

```yaml
telemetry:
  enabled: true
  service_name: "echomind-backend"
  service_version: "v1.2.0"
  environment: "development"  # development, staging, production
  
  # Exporter 配置
  exporter:
    type: "console"  # console, file, otlp
    
    # 控制台输出 (开发环境)
    console:
      enable_color: true
      pretty_print: true
    
    # 文件输出
    file:
      traces_path: "./logs/traces.jsonl"
      metrics_path: "./logs/metrics.jsonl"
    
    # OTLP Collector (生产环境)
    otlp:
      endpoint: "localhost:4318"
      insecure: true
      timeout: "10s"
  
  # 采样策略
  sampling:
    type: "always_on"  # always_on, always_off, traceidratio
    ratio: 1.0         # 采样率 (0.0-1.0)
  
  # Metrics 配置
  metrics:
    export_interval: "60s"  # 指标导出间隔
    export_timeout: "30s"
```

### 配置项说明

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `enabled` | bool | false | 是否启用 OTel |
| `service_name` | string | echomind-backend | 服务名称 |
| `environment` | string | development | 环境标识 |
| `exporter.type` | string | console | 导出器类型 |
| `sampling.type` | string | always_on | 采样策略 |
| `sampling.ratio` | float | 1.0 | 采样比例 |

---

## 使用指南

### 1. Tracing (分布式追踪)

#### 在服务中使用

```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/codes"
)

var tracer = otel.Tracer("echomind.myservice")

func (s *MyService) DoSomething(ctx context.Context) error {
    // 创建 span
    ctx, span := tracer.Start(ctx, "MyService.DoSomething",
        trace.WithAttributes(
            attribute.String("user.id", userID),
        ),
    )
    defer span.End()

    // 业务逻辑
    result, err := s.process(ctx)
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, "processing failed")
        return err
    }

    // 记录成功
    span.SetStatus(codes.Ok, "success")
    span.SetAttributes(
        attribute.Int("result.count", len(result)),
    )

    return nil
}
```

#### 嵌套 Span

```go
func (s *Service) ComplexOperation(ctx context.Context) error {
    ctx, span := tracer.Start(ctx, "ComplexOperation")
    defer span.End()

    // 子操作 1
    ctx, span1 := tracer.Start(ctx, "step1")
    step1Result := s.step1(ctx)
    span1.End()

    // 子操作 2
    ctx, span2 := tracer.Start(ctx, "step2")
    step2Result := s.step2(ctx)
    span2.End()

    return nil
}
```

### 2. Metrics (指标收集)

#### 使用预定义 Metrics

```go
import "github.com/hrygo/echomind/pkg/telemetry"

// 在服务初始化时创建 metrics
metrics, err := telemetry.NewSearchMetrics(ctx)

// 记录延迟
start := time.Now()
// ... 业务逻辑 ...
metrics.RecordSearchLatency(ctx, time.Since(start))

// 增加计数器
metrics.IncrementSearchRequests(ctx)

// 记录活跃数
metrics.IncrementActiveSearches(ctx)
defer metrics.DecrementActiveSearches(ctx)
```

#### 自定义 Metrics

```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/metric"
)

meter := otel.Meter("echomind.custom")

counter, _ := meter.Int64Counter(
    "custom.requests.total",
    metric.WithDescription("Total custom requests"),
)

counter.Add(ctx, 1)
```

### 3. Logs Correlation (日志关联)

#### 自动注入 TraceID

```go
import "github.com/hrygo/echomind/pkg/logger"

func Handler(ctx context.Context) {
    // 自动包含 trace_id 和 span_id
    logger.InfoCtx(ctx, "Processing request",
        logger.String("user_id", userID),
    )
}
```

输出示例:

```json
{
  "level": "info",
  "msg": "Processing request",
  "trace_id": "4bf92f3577b34da6a3ce929d0e0e4736",
  "span_id": "00f067aa0ba902b7",
  "user_id": "123"
}
```

#### 手动创建带 Trace 的日志

```go
// 获取带 trace 字段的 logger
traceLogger := logger.WithTraceContext(ctx)
traceLogger.Info("Operation completed")
```

---

## 指标说明

### SearchService Metrics

| 指标名称 | 类型 | 单位 | 说明 |
|---------|------|------|------|
| `search.latency` | Histogram | ms | 搜索端到端延迟 |
| `embedding.latency` | Histogram | ms | 嵌入生成延迟 |
| `db.query.latency` | Histogram | ms | 数据库查询延迟 |
| `search.requests.total` | Counter | count | 搜索请求总数 |
| `search.errors.total` | Counter | count | 搜索错误总数 |
| `cache.hits.total` | Counter | count | 缓存命中次数 |
| `cache.misses.total` | Counter | count | 缓存未命中次数 |
| `search.active` | UpDownCounter | count | 当前活跃搜索数 |
| `search.results.total` | Counter | count | 返回结果总数 |

### SyncService Metrics

| 指标名称 | 类型 | 单位 | 说明 |
|---------|------|------|------|
| `sync.latency` | Histogram | ms | 同步操作延迟 |
| `sync.requests.total` | Counter | count | 同步请求总数 |
| `sync.emails.processed` | Counter | count | 处理邮件数 |

### AIService Metrics

| 指标名称 | 类型 | 单位 | 说明 |
|---------|------|------|------|
| `ai.request.latency` | Histogram | ms | AI API 延迟 |
| `ai.requests.total` | Counter | count | AI 请求总数 |
| `ai.errors.total` | Counter | count | AI 错误总数 |
| `ai.tokens.used` | Counter | count | Token 消耗量 |

---

## 最佳实践

### 1. Span 命名规范

- **格式**: `ServiceName.MethodName`
- **示例**: `SearchService.Search`, `AIProvider.Embed`

### 2. Attribute 命名规范

使用 OpenTelemetry 语义约定:

- `user.id` - 用户 ID
- `http.method` - HTTP 方法
- `db.statement` - 数据库查询
- `error.type` - 错误类型

### 3. 采样策略

| 环境 | 采样策略 | 采样率 | 原因 |
|------|----------|--------|------|
| **Development** | always_on | 1.0 | 完整追踪便于调试 |
| **Staging** | traceidratio | 0.5 | 平衡成本和可见性 |
| **Production** | traceidratio | 0.1 | 降低存储成本 |

### 4. 性能考虑

- Span 创建开销: ~1-2µs
- Metrics 记录开销: ~0.5µs
- 总体性能影响: < 5%

### 5. 错误处理

```go
if err != nil {
    span.RecordError(err)
    span.SetStatus(codes.Error, err.Error())
    logger.ErrorCtx(ctx, "Operation failed",
        logger.Error(err),
    )
    return err
}
```

---

## 故障排查

### 问题 1: Telemetry 初始化失败

**症状**: 日志显示 "Failed to initialize telemetry"

**解决**:
1. 检查配置文件路径
2. 验证 YAML 语法
3. 查看详细错误信息

```bash
# 验证配置
cat backend/configs/config.yaml | grep -A20 telemetry
```

### 问题 2: Traces 未输出

**症状**: 控制台没有 trace 输出

**解决**:
1. 确认 `telemetry.enabled: true`
2. 检查 `exporter.type` 配置
3. 验证 span 是否被创建

```go
// 添加调试日志
span := trace.SpanFromContext(ctx)
if !span.IsRecording() {
    log.Warn("Span is not recording")
}
```

### 问题 3: Metrics 数据异常

**症状**: Metrics 值不正确

**解决**:
1. 检查 Metrics 初始化
2. 验证 Context 传递
3. 确认导出间隔配置

```yaml
metrics:
  export_interval: "10s"  # 缩短间隔便于调试
```

### 问题 4: 性能下降

**症状**: 启用 OTel 后性能下降 > 5%

**解决**:
1. 降低采样率
2. 使用批量导出
3. 检查 span 数量

```yaml
sampling:
  type: "traceidratio"
  ratio: 0.1  # 仅采样 10%
```

---

## 进阶主题

### 1. 集成 Jaeger

```yaml
exporter:
  type: "otlp"
  otlp:
    endpoint: "jaeger:4318"
    insecure: true
```

部署 Jaeger:

```bash
docker run -d --name jaeger \
  -p 16686:16686 \
  -p 4318:4318 \
  jaegertracing/all-in-one:latest
```

访问: http://localhost:16686

### 2. 集成 Prometheus

```yaml
exporter:
  type: "otlp"
  otlp:
    endpoint: "otel-collector:4318"
```

配置 OpenTelemetry Collector:

```yaml
receivers:
  otlp:
    protocols:
      http:
        endpoint: 0.0.0.0:4318

exporters:
  prometheus:
    endpoint: "0.0.0.0:8889"

service:
  pipelines:
    metrics:
      receivers: [otlp]
      exporters: [prometheus]
```

### 3. 自定义 Exporter

```go
import "go.opentelemetry.io/otel/exporters/..."

// 实现自定义 exporter
type CustomExporter struct {
    // ...
}

func (e *CustomExporter) Export(ctx context.Context, spans []trace.ReadOnlySpan) error {
    // 自定义导出逻辑
    return nil
}
```

---

## 参考资源

- [OpenTelemetry 官方文档](https://opentelemetry.io/docs/)
- [Go SDK 文档](https://pkg.go.dev/go.opentelemetry.io/otel)
- [语义约定](https://github.com/open-telemetry/semantic-conventions)
- [EchoMind 时序图文档](./api_search_sequence_diagram.md)

---

**维护者**: EchoMind Team  
**更新频率**: 随版本更新
