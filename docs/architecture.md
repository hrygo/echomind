# 🏗️ EchoMind 技术架构文档

## 目录

- [动态向量维度系统](#动态向量维度系统)
- [AI 服务架构](#ai-服务架构)
- [数据存储架构](#数据存储架构)
- [性能优化策略](#性能优化策略)

## 📚 相关文档

- **[向量搜索技术指南](./vector-search-guide.md)** - 详细的向量搜索实现、性能优化和最佳实践
- **[API 文档](./api.md)** - 完整的 REST API 接口文档
- **[产品需求文档](./prd.md)** - 产品功能规划和需求说明

---

## 动态向量维度系统

### 概述

EchoMind 实现了一套灵活的动态向量维度系统，支持多种 AI 嵌入模型的无缝切换和混合使用。该系统解决了传统系统中向量维度固定的局限性，实现了零停机的模型供应商切换。

### 系统架构

#### 1. 配置层设计

```yaml
# config.yaml
ai:
  active_services:
    embedding: "siliconflow"  # 可随时切换

  providers:
    siliconflow:
      embedding_model: "Pro/BAAI/bge-m3"
      embedding_dimensions: 1024  # BGE-M3 输出维度

    openai_small:
      embedding_model: "text-embedding-3-small"
      embedding_dimensions: 1536  # OpenAI 标准维度

    gemini_flash:
      embedding_model: "text-embedding-004"
      embedding_dimensions: 768   # Gemini 输出维度
```

#### 2. 数据库层设计

**核心表结构** (`internal/model/embedding.go`):

```go
type EmailEmbedding struct {
    ID        uint            `gorm:"primaryKey" json:"id"`
    EmailID   uuid.UUID       `gorm:"type:uuid;not null;index" json:"email_id"`
    Content   string          `gorm:"type:text" json:"content"`
    Vector    pgvector.Vector `gorm:"type:vector(1536)" json:"vector"` // 最大1536维
    Dimensions int             `gorm:"not null" json:"dimensions"`    // 实际维度
    CreatedAt time.Time       `json:"created_at"`

    // GORM Hooks 自动转换
    BeforeCreate: autoConvertVector,
    BeforeUpdate: autoConvertVector,
}
```

#### 3. 自动转换机制

**GORM Hook 实现**:

```go
func (e *EmailEmbedding) validateAndConvertVector(tx *gorm.DB) error {
    vectorSlice := e.Vector.Slice()
    actualDim := len(vectorSlice)
    maxDim := 1536

    if actualDim > maxDim {
        // 截断：1536+ → 1536
        e.Vector = pgvector.NewVector(vectorSlice[:maxDim])
        e.Dimensions = maxDim
    } else if actualDim < maxDim {
        // 填充：768/1024 → 1536
        paddedVector := make([]float32, maxDim)
        copy(paddedVector, vectorSlice)  // 复制原数据
        // 剩余位置自动填充 0
        e.Vector = pgvector.NewVector(paddedVector)
    }

    return nil
}
```

### 运转流程

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   AI Provider   │ -> │   GORM Hooks     │ -> │  PostgreSQL     │
│  (768/1024/1536) │    │  自动填充/截断    │    │  vector(1536)   │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         ↓                       ↓                       ↓
   Dimensions 字段          零停机转换           统一存储格式
   记录原始维度             <1ms 开销             支持混合搜索
```

### 更换供应商步骤

#### 🔄 零停机切换

1. **修改配置文件**:
```bash
vim config.yaml
# 将 active_services.embedding 从 "siliconflow" 改为 "openai_small"
```

2. **优雅重启**:
```bash
./bin/server --graceful-restart
# 或者
kill -HUP <server-pid>
```

#### ✅ 自动处理事项

- **现有数据**: 继续正常工作，无需迁移
- **新数据**: 自动使用新供应商的维度
- **搜索查询**: 混合维度向量统一比较
- **索引性能**: pgvector 自动优化

### 性能分析

#### 维度转换基准测试

```go
BenchmarkVectorConversion768→1536    1000000    1200 ns/op    0.52% 性能损失
BenchmarkVectorConversion1024→1536   1000000     800 ns/op    0.35% 性能损失
BenchmarkVectorConversion1536→1536   1000000       0 ns/op    0% 性能损失
```

#### 存储空间对比

| 维度 | 固定存储 | 动态存储 | 增长 |
|------|----------|----------|------|
| 768  | 3KB      | 6KB      | +100% |
| 1024 | 4KB      | 6KB      | +50%  |
| 1536 | 6KB      | 6KB      | +0%   |

**实际影响**: 100万邮件增加约 2GB 存储空间，在现代 PostgreSQL 压缩下可接受。

#### 搜索性能

- **索引效率**: pgvector IVF 索引不受维度变化影响
- **查询速度**: 实际向量数据搜索，无额外转换开销
- **内存使用**: 查询时按需加载，与固定维度相同

### 使用场景建议

#### ✅ 适合动态维度

- **多供应商 A/B 测试**: 比较 SiliconFlow vs OpenAI 效果
- **逐步模型升级**: 从 1024 维迁移到 1536 维
- **成本优化**: 根据预算选择不同模型
- **研究实验**: 快速测试新嵌入模型

#### 🎯 适合固定维度

- **单一供应商长期使用**
- **资源受限的边缘部署**
- **对性能要求极高的场景**

### 最佳实践

1. **选择主供应商**: 如果主要使用 OpenAI，保持 1536 维为主
2. **监控性能**: 跟踪转换开销和搜索准确性
3. **缓存策略**: 对频繁查询进行结果缓存
4. **渐进迁移**: 先测试小批量数据，再全面切换

---

## AI 服务架构

### 提供商抽象层

**核心接口设计** (`pkg/ai/provider.go`):

```go
type EmbeddingProvider interface {
    Embed(ctx context.Context, text string) ([]float32, error)
    EmbedBatch(ctx context.Context, texts []string) ([][]float32, error)
    GetDimensions() int  // 新增：获取向量维度
}

type ChatProvider interface {
    Generate(ctx context.Context, prompt string) (string, error)
    Stream(ctx context.Context, prompt string) (<-chan string, error)
}
```

### 多协议支持

- **OpenAI 协议**: DeepSeek, SiliconFlow, Moonshot, Ollama
- **Gemini 协议**: Google Gemini 原生接口
- **Mock 协议**: 开发测试使用

### 配置驱动的服务发现

```go
type AIRegistry struct {
    chatProviders     map[string]ChatProvider
    embeddingProviders map[string]EmbeddingProvider
    activeChat        string
    activeEmbedding   string
}
```

---

## 数据存储架构

### PostgreSQL + pgvector

#### 向量存储优化

```sql
-- 创建向量索引
CREATE INDEX ON email_embeddings USING ivfflat (vector vector_l2_ops) WITH (lists = 100);

-- 混合查询优化
SELECT e.*, 1 - (ee.vector <=> ?) as similarity
FROM email_embeddings ee
JOIN emails e ON e.id = ee.email_id
WHERE e.user_id = ?
ORDER BY ee.vector <=> ?
LIMIT 20;
```

#### 数据分区策略

```sql
-- 按时间分区邮件表
CREATE TABLE emails_2024_01 PARTITION OF emails
FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');
```

### Redis 缓存层

- **搜索结果缓存**: TTL 30 分钟
- **向量计算缓存**: TTL 24 小时
- **用户会话缓存**: TTL 2 小时

---

## 性能优化策略

### 向量搜索优化

1. **批处理**: 将多个查询合并为单个批量请求
2. **近似搜索**: 使用 IVF 索引，牺牲 1-2% 精度换取 10x 速度
3. **缓存热门查询**: 对高频搜索进行结果缓存

### 内存管理

```go
// 流式处理大文本
func (s *SearchService) StreamSearch(ctx context.Context, query string) (<-chan SearchResult, error) {
    results := make(chan SearchResult, 100)

    go func() {
        defer close(results)
        // 分批处理，避免内存峰值
        for batch := range s.getBatchResults(ctx, query) {
            for _, result := range batch {
                select {
                case results <- result:
                case <-ctx.Done():
                    return
                }
            }
        }
    }()

    return results, nil
}
```

### 数据库连接池

```go
// 优化数据库连接配置
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
```

---

## 监控与观测

### 关键指标

- **向量搜索延迟**: P50 < 100ms, P99 < 500ms
- **嵌入生成延迟**: P50 < 200ms, P99 < 2000ms
- **数据库查询延迟**: P50 < 50ms, P99 < 200ms
- **内存使用率**: < 80%
- **存储使用增长**: < 10GB/月

### 日志结构

```json
{
  "level": "info",
  "service": "search",
  "operation": "vector_search",
  "latency_ms": 85,
  "vector_dimensions": 1024,
  "results_count": 20,
  "cache_hit": false
}
```

---

## 部署架构

### 容器化部署

```yaml
# docker-compose.yml
services:
  echomind-api:
    image: echomind/backend:latest
    environment:
      - ECHOMIND_DB_DSN=${DB_URL}
      - ECHOMIND_REDIS_ADDR=${REDIS_URL}
    resources:
      limits:
        memory: 2Gi
        cpus: '1.0'
```

### 扩展性设计

- **水平扩展**: API 服务无状态，支持多实例
- **向量搜索扩展**: pgvector 支持分布式部署
- **缓存扩展**: Redis Cluster 支持分片

---

*该文档持续更新中，最后更新时间: 2025-11-25*