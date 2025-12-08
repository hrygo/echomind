# 🏗️ EchoMind 技术架构文档

## 目录


- [AI 服务架构](#ai-服务架构)
- [数据存储架构](#数据存储架构)
- [性能优化策略](#性能优化策略)

## 📚 相关文档

- **[EchoMind 邮件处理系统时序图](./api_search_sequence_diagram.md)** - 完整的邮件搜索、同步和Reindex流程时序图
- **[向量搜索技术指南](./vector-search-guide.md)** - 详细的向量搜索实现、性能优化和最佳实践
- **[API 文档](./api.md)** - 完整的 REST API 接口文档
- **[产品需求文档](./prd.md)** - 产品功能规划和需求说明

---

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