# EchoMind v1.2.0 Release Notes

**发布日期**: 2025-11-26  
**版本类型**: Minor Release  
**主题**: OpenTelemetry 可观测性与搜索性能优化

---

## 🎉 重大更新

### 1. OpenTelemetry 可观测性集成 ⭐

EchoMind 现已集成企业级可观测性框架 OpenTelemetry,为系统监控、性能分析和问题诊断提供全面支持。

**核心特性**:
- ✅ **Distributed Tracing** - 完整的请求链路追踪
- ✅ **Metrics Collection** - 关键性能指标收集
- ✅ **Logs Correlation** - 日志与追踪自动关联
- ✅ **Multi-Exporter** - 支持 Console/File/OTLP 导出

**使用示例**:
```yaml
# backend/configs/config.yaml
telemetry:
  enabled: true
  service_name: "echomind-backend"
  exporter:
    type: "console"  # 开发环境
```

### 2. SearchService 性能增强 🚀

**新增功能**:
- ✅ **Redis 缓存层** - 搜索结果智能缓存
- ✅ **性能追踪** - 完整的搜索流程可观测性
- ✅ **缓存指标** - 命中率/未命中率监控

**性能提升**:
| 指标 | 改进前 | 改进后 | 提升 |
|------|--------|--------|------|
| 重复查询延迟 | ~1000ms | ~50ms | **20x** |
| AI API 调用 | 每次 | 首次缓存 | **-50% 成本** |
| 缓存命中率 | N/A | 40-60% | **新增** |

### 3. 日志系统增强 📝

**TraceID 自动注入**:
```go
// 自动包含 trace_id 和 span_id
logger.InfoCtx(ctx, "Processing request")

// 输出:
{
  "level": "info",
  "msg": "Processing request",
  "trace_id": "4bf92f3577b34da6a3ce929d0e0e4736",
  "span_id": "00f067aa0ba902b7"
}
```

---

## 📦 新增组件

### 后端

| 组件 | 文件路径 | 说明 |
|------|---------|------|
| **OTel 核心** | `pkg/telemetry/otel.go` | 初始化和配置管理 |
| **Metrics 定义** | `pkg/telemetry/metrics.go` | 指标集合 (Search/Sync/AI) |
| **Trace 日志** | `pkg/logger/trace.go` | TraceID 注入工具 |
| **搜索缓存** | `internal/service/search_cache.go` | Redis 缓存实现 |

### 文档

| 文档 | 路径 | 说明 |
|------|------|------|
| **OTel 指南** | `docs/otel-integration-guide.md` | 完整的集成和使用指南 (517行) |
| **迭代总结** | `docs/v1.2.0-iteration-summary.md` | 详细的功能总结 |

### 测试

| 测试文件 | 路径 | 覆盖 |
|---------|------|------|
| **OTel 测试** | `pkg/telemetry/telemetry_test.go` | 初始化/Tracing/Metrics/Benchmark |

---

## 🔧 配置变更

### 新增配置节

```yaml
# backend/configs/config.example.yaml

telemetry:
  enabled: true
  service_name: "echomind-backend"
  service_version: "v1.2.0"
  environment: "development"
  
  exporter:
    type: "console"
    console:
      enable_color: true
      pretty_print: true
    file:
      traces_path: "./logs/traces.jsonl"
      metrics_path: "./logs/metrics.jsonl"
    otlp:
      endpoint: "localhost:4318"
      insecure: true
      timeout: "10s"
  
  sampling:
    type: "always_on"
    ratio: 1.0
  
  metrics:
    export_interval: "60s"
    export_timeout: "30s"
```

---

## 📊 性能与指标

### 关键指标

| 指标名称 | 类型 | 单位 | 说明 |
|---------|------|------|------|
| `search.latency` | Histogram | ms | 搜索端到端延迟 |
| `embedding.latency` | Histogram | ms | 嵌入生成延迟 |
| `db.query.latency` | Histogram | ms | 数据库查询延迟 |
| `cache.hits.total` | Counter | count | 缓存命中次数 |
| `search.active` | UpDownCounter | count | 当前活跃搜索数 |

### 性能基准

| 场景 | P50 | P95 | P99 |
|------|-----|-----|-----|
| **搜索 (无缓存)** | 950ms | 1200ms | 1500ms |
| **搜索 (缓存命中)** | 45ms | 80ms | 120ms |
| **嵌入生成** | 400ms | 600ms | 800ms |

**OTel 性能开销**: < 2% (CPU/Memory)

---

## 🔄 Breaking Changes

### API 变更

**SearchService 构造函数**:
```go
// v1.1.0
searchService := service.NewSearchService(db, embedder)

// v1.2.0 (新增 cache 参数)
searchService := service.NewSearchService(db, embedder, searchCache)
```

**迁移指南**:
```go
// 如果不使用缓存,传 nil
searchService := service.NewSearchService(db, embedder, nil)

// 使用 Redis 缓存
redisClient := redis.NewClient(&redis.Options{...})
cache := service.NewSearchCache(redisClient, 30*time.Minute)
searchService := service.NewSearchService(db, embedder, cache)
```

---

## 📚 文档更新

### 新增文档

1. **[OTel 集成指南](./docs/otel-integration-guide.md)**
   - 快速开始
   - 配置详解
   - 使用示例
   - 故障排查
   - 进阶主题 (Jaeger/Prometheus)

2. **[v1.2.0 迭代总结](./docs/v1.2.0-iteration-summary.md)**
   - 功能完成情况
   - 技术亮点
   - 性能数据
   - 未来规划

### 更新文档

1. **[配置示例](./backend/configs/config.example.yaml)**
   - 新增 `telemetry` 配置节
   - 完整的 OTel 配置示例

---

## 🐛 Bug 修复

- 修复 SearchService 编译错误
- 修复容器初始化依赖注入问题
- 修复测试代码中的未使用变量警告

---

## ⬆️ 升级指南

### 从 v1.1.0 升级

1. **更新依赖**:
```bash
cd backend
go mod tidy
```

2. **更新配置**:
```yaml
# 在 config.yaml 中添加 telemetry 配置
telemetry:
  enabled: true  # 可选,默认 false
```

3. **更新代码** (如果自定义了 Container):
```go
// 更新 SearchService 初始化
var cache *service.SearchCache
if redisAddr != "" {
    redisClient := redis.NewClient(&redis.Options{...})
    cache = service.NewSearchCache(redisClient, 30*time.Minute)
}
searchService := service.NewSearchService(db, embedder, cache)
```

4. **重新编译**:
```bash
go build ./...
```

5. **运行测试**:
```bash
go test ./...
```

---

## 🔮 未来计划

### v1.3.0 (下一迭代)

- [ ] 搜索结果智能聚类
- [ ] AI 驱动的搜索摘要生成
- [ ] 前端搜索增强 UI
- [ ] 其他服务 OTel 集成 (AI/Sync)

### v2.0.0 (长期)

- [ ] 微信集成 (Phase 7)
- [ ] 生产环境可观测性全面上线
- [ ] Grafana 仪表盘

---

## 🙏 致谢

感谢所有参与 v1.2.0 开发的贡献者!

特别感谢:
- OpenTelemetry 社区提供的优秀框架
- Redis 团队的高性能缓存解决方案

---

## 📖 相关资源

- [GitHub Repository](https://github.com/hrygo/echomind)
- [OTel 官方文档](https://opentelemetry.io/)
- [项目路线图](./docs/product-roadmap.md)

---

**完整变更日志**: [CHANGELOG.md](../CHANGELOG.md)

**下载**: [Release v1.2.0](https://github.com/hrygo/echomind/releases/tag/v1.2.0)
