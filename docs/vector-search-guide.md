# 🔍 EchoMind 向量搜索技术指南

## 概述

EchoMind 的向量搜索系统基于 PostgreSQL + pgvector 构建，实现了高性能的语义搜索和 RAG (Retrieval-Augmented Generation) 功能。本指南详细介绍了向量搜索的技术实现、性能优化和最佳实践。

## 🏗️ 系统架构

### 向量搜索流程图

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   用户查询      │ -> │  文本预处理      │ -> │  嵌入生成       │ -> │  向量搜索       │
│   "项目进展"     │    │  清洗、分块      │    │  768/1024/1536  │    │  相似度计算     │
└─────────────────┘    └─────────────────┘    └─────────────────┘    └─────────────────┘
                                                                               ↓
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   结果排序      │ <- │  分数过滤       │ <- │  后处理         │ <- │  数据库查询     │
│   相关性排序     │    │  阈值过滤       │    │  去重、聚合      │    │  近似最近邻     │
└─────────────────┘    └─────────────────┘    └─────────────────┘    └─────────────────┘
```

### 核心组件

#### 1. 嵌入生成服务 (`pkg/ai/embedding.go`)

```go
type EmbeddingService struct {
    provider ai.EmbeddingProvider
    cache    *redis.Client
    metrics  *prometheus.Registry
}

// 批量生成嵌入，提高效率
func (s *EmbeddingService) GenerateEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
    // 1. 检查缓存
    cached := s.getCachedEmbeddings(texts)

    // 2. 批量生成未缓存的
    uncached := s.getUncachedTexts(texts, cached)
    vectors, err := s.provider.EmbedBatch(ctx, uncached)

    // 3. 缓存结果
    s.cacheEmbeddings(uncached, vectors)

    // 4. 合并结果
    return s.mergeResults(cached, vectors), nil
}
```

#### 2. 向量搜索引擎 (`internal/service/search.go`)

```go
type SearchEngine struct {
    db          *gorm.DB
    embedder    ai.EmbeddingProvider
    cache       *redis.Client
    indexConfig *IndexConfig
}

// 高性能向量搜索
func (s *SearchEngine) SemanticSearch(ctx context.Context, req SearchRequest) (*SearchResponse, error) {
    // 1. 查询向量化
    queryVector, err := s.embedder.Embed(ctx, req.Query)

    // 2. 构建优化查询
    sql := s.buildOptimizedQuery(req)

    // 3. 执行向量相似度搜索
    var results []EmailEmbedding
    err = s.db.WithContext(ctx).Raw(sql,
        pgvector.NewVector(queryVector),
        req.UserID,
        req.Limit,
    ).Scan(&results).Error

    // 4. 后处理和排序
    return s.postProcessResults(results, req), nil
}
```

## 🎯 向量数据库优化

### 索引策略

#### 1. IVFFlat 索引 (适合中等规模数据)

```sql
-- 创建 IVFFlat 索引
CREATE INDEX CONCURRENTLY email_embeddings_vector_idx
ON email_embeddings
USING ivfflat (vector vector_l2_ops)
WITH (lists = 100);

-- 监控索引效果
EXPLAIN (ANALYZE, BUFFERS)
SELECT id, 1 - (vector <=> '[0.1,0.2,...]') as similarity
FROM email_embeddings
ORDER BY vector <=> '[0.1,0.2,...]'
LIMIT 20;
```

#### 2. HNSW 索引 (适合大规模数据)

```sql
-- 创建 HNSW 索引 (更高的召回率)
CREATE INDEX CONCURRENTLY email_embeddings_hnsw_idx
ON email_embeddings
USING hnsw (vector vector_cosine_ops)
WITH (m = 16, ef_construction = 64);

-- 调整查询时的 ef 参数以提高精度
SET hnsw.ef_search = 100;
```

### 数据分区优化

```sql
-- 按时间分区向量表
CREATE TABLE email_embeddings_2024_q1 PARTITION OF email_embeddings
FOR VALUES FROM ('2024-01-01') TO ('2024-04-01');

-- 自动创建新分区
CREATE OR REPLACE FUNCTION create_monthly_partition()
RETURNS void AS $$
DECLARE
    start_date date;
    end_date date;
BEGIN
    start_date := date_trunc('month', CURRENT_DATE);
    end_date := start_date + interval '1 month';

    EXECUTE format('CREATE TABLE IF NOT EXISTS email_embeddings_%s PARTITION OF email_embeddings
                    FOR VALUES FROM (%L) TO (%L)',
                   to_char(start_date, 'YYYY_MM'),
                   start_date, end_date);
END;
$$ LANGUAGE plpgsql;
```

## ⚡ 性能优化策略

### 1. 查询优化

#### 缓存策略

```go
type SearchCache struct {
    redis  *redis.Client
    local  *sync.Map
    ttl    time.Duration
}

// 多层缓存策略
func (c *SearchCache) Get(key string) (*SearchResult, bool) {
    // L1: 内存缓存 (最快)
    if result, ok := c.local.Load(key); ok {
        return result.(*SearchResult), true
    }

    // L2: Redis 缓存 (中等)
    if cached, err := c.redis.Get(ctx, key).Result(); err == nil {
        var result SearchResult
        json.Unmarshal([]byte(cached), &result)
        c.local.Store(key, &result) // 回填 L1
        return &result, true
    }

    return nil, false
}
```

#### 批量搜索

```go
// 批量处理多个搜索查询
func (s *SearchEngine) BatchSearch(ctx context.Context, queries []string) ([]*SearchResponse, error) {
    // 1. 批量生成查询向量
    queryVectors, err := s.embedder.EmbedBatch(ctx, queries)

    // 2. 并行执行搜索
    var wg sync.WaitGroup
    results := make([]*SearchResponse, len(queries))

    for i, query := range queries {
        wg.Add(1)
        go func(idx int, q string) {
            defer wg.Done()
            results[idx], _ = s.SingleSearch(ctx, q, queryVectors[idx])
        }(i, query)
    }

    wg.Wait()
    return results, nil
}
```

### 2. 内存优化

#### 流式处理

```go
// 流式处理大量搜索结果
func (s *SearchEngine) StreamSearch(ctx context.Context, req SearchRequest) (<-chan SearchResult, error) {
    resultChan := make(chan SearchResult, 100)

    go func() {
        defer close(resultChan)

        // 分页处理避免内存峰值
        pageSize := 1000
        for offset := 0; ; offset += pageSize {
            batch, err := s.searchBatch(ctx, req, offset, pageSize)
            if err != nil || len(batch) == 0 {
                break
            }

            for _, result := range batch {
                select {
                case resultChan <- result:
                case <-ctx.Done():
                    return
                }
            }
        }
    }()

    return resultChan, nil
}
```

### 3. 向量化优化

#### 文本预处理

```go
type TextPreprocessor struct {
    chunker    *TextChunker
    cleaner    *TextCleaner
    tokenizer  *Tokenizer
}

// 优化文本分块策略
func (p *TextPreprocessor) ProcessText(text string, maxTokens int) []string {
    // 1. 文本清洗
    cleaned := p.cleaner.Clean(text)

    // 2. 智能分块 (保持语义完整性)
    chunks := p.chunker.ChunkWithOverlap(cleaned, maxTokens, 0.1)

    // 3. 质量过滤
    filtered := p.filterLowQualityChunks(chunks)

    return filtered
}

// 质量过滤标准
func (p *TextPreprocessor) filterLowQualityChunks(chunks []string) []string {
    var result []string
    for _, chunk := range chunks {
        if len(chunk) < 50 ||
           strings.Count(chunk, " ") < 5 ||
           p.hasRepetitiveContent(chunk) {
            continue
        }
        result = append(result, chunk)
    }
    return result
}
```

## 📊 性能监控

### 关键指标

#### 1. 搜索性能指标

```go
var (
    searchDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "vector_search_duration_seconds",
            Help:    "Duration of vector search operations",
            Buckets: prometheus.DefBuckets,
        },
        []string{"user_id", "query_length"},
    )

    cacheHitRate = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "search_cache_hits_total",
            Help: "Total number of search cache hits",
        },
        []string{"cache_level"},
    )

    indexUsage = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "vector_index_usage_ratio",
            Help: "Ratio of index usage in vector queries",
        },
        []string{"index_type"},
    )
)
```

#### 2. 数据库监控

```sql
-- 监控向量索引使用情况
SELECT
    schemaname,
    tablename,
    indexname,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch
FROM pg_stat_user_indexes
WHERE tablename = 'email_embeddings';

-- 监控查询性能
SELECT
    query,
    calls,
    total_exec_time,
    mean_exec_time,
    rows
FROM pg_stat_statements
WHERE query LIKE '%vector%'
ORDER BY total_exec_time DESC;
```

### 性能基准

#### 搜索延迟基准

| 数据量 | 维度 | 索引类型 | P50 延迟 | P99 延迟 | QPS |
|--------|------|----------|----------|----------|-----|
| 10万   | 768  | IVFFlat  | 15ms     | 45ms     | 200 |
| 100万  | 768  | IVFFlat  | 25ms     | 80ms     | 150 |
| 1000万 | 768  | HNSW     | 35ms     | 120ms    | 100 |
| 10万   | 1536 | IVFFlat  | 20ms     | 60ms     | 180 |
| 100万  | 1536 | HNSW     | 30ms     | 100ms    | 120 |

#### 存储性能基准

| 维度 | 每向量大小 | 100万向量存储 | 索引大小 | 搜索内存 |
|------|------------|----------------|----------|----------|
| 768  | 3KB        | 3GB            | 500MB    | 200MB    |
| 1024 | 4KB        | 4GB            | 700MB    | 300MB    |
| 1536 | 6KB        | 6GB            | 1.2GB    | 500MB    |

## 🔧 配置最佳实践

### 1. 环境配置

```yaml
# config.yaml
database:
  # 连接池优化
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: 300s

vector:
  # 索引配置
  index_type: "hnsw"  # ivfflat 或 hnsw
  hnsw:
    m: 16
    ef_construction: 64
  ivfflat:
    lists: 100

search:
  # 搜索配置
  default_limit: 20
  max_limit: 100
  cache_ttl: 1800s  # 30分钟

embedding:
  # 嵌入配置
  batch_size: 32
  chunk_size: 1000
  overlap_ratio: 0.1
```

### 2. 生产环境调优

```sql
-- PostgreSQL 配置优化
-- postgresql.conf
shared_buffers = 2GB                    -- 25% of RAM
effective_cache_size = 6GB              -- 75% of RAM
work_mem = 64MB                         -- Per query memory
maintenance_work_mem = 256MB
random_page_cost = 1.1                  -- SSD optimization
seq_page_cost = 1.0

-- pgvector 特定优化
hnsw.ef_search = 100                    -- Higher ef for better recall
ivfflat.probes = 10                     -- Number of probes to examine
```

## 🧪 测试和调试

### 1. 单元测试

```go
func TestVectorSearch(t *testing.T) {
    // 准备测试数据
    testVectors := generateTestVectors(1000, 768)

    // 插入测试数据
    for _, vec := range testVectors {
        db.Create(&EmailEmbedding{
            Vector: pgvector.NewVector(vec),
        })
    }

    // 执行搜索测试
    query := generateTestVector(768)
    results, err := searchEngine.Search(query, 10)

    assert.NoError(t, err)
    assert.Len(t, results, 10)

    // 验证相似度排序
    for i := 1; i < len(results); i++ {
        assert.GreaterOrEqual(t, results[i-1].Score, results[i].Score)
    }
}
```

### 2. 集成测试

```go
func TestSearchPerformance(t *testing.T) {
    // 性能基准测试
    b := testing.B{}

    for i := 0; i < b.N; i++ {
        start := time.Now()
        _, err := searchEngine.SemanticSearch(ctx, SearchRequest{
            Query: "test query",
            Limit: 20,
        })
        duration := time.Since(start)

        // 记录性能指标
        prometheus.RecordDuration(duration)

        assert.NoError(t, err)
        assert.Less(t, duration, 100*time.Millisecond)
    }
}
```

### 3. 调试工具

```go
// 搜索调试工具
type SearchDebugger struct {
    logger *log.Logger
}

func (d *SearchDebugger) DebugSearch(ctx context.Context, query string) {
    d.logger.Printf("=== Search Debug ===")
    d.logger.Printf("Query: %s", query)

    // 1. 记录查询向量
    vector, _ := d.embedder.Embed(ctx, query)
    d.logger.Printf("Query vector dimensions: %d", len(vector))

    // 2. 记录查询计划
    var explain []string
    db.Raw("EXPLAIN (ANALYZE, BUFFERS) " + searchSQL, vector).Scan(&explain)
    d.logger.Printf("Query plan: %v", explain)

    // 3. 记录搜索结果
    results, _ := d.searchEngine.Search(ctx, query, 5)
    d.logger.Printf("Results count: %d", len(results))
    for i, result := range results {
        d.logger.Printf("Result %d: score=%.4f, id=%s", i, result.Score, result.ID)
    }
}
```

## 🚀 未来优化方向

### 1. 高级索引策略
- **DiskANN**: 更大规模的向量索引
- **Quantization**: 向量量化减少存储空间
- **Hybrid Search**: 结合关键词搜索和向量搜索

### 2. 智能缓存
- **LRU Cache**: 更智能的缓存策略
- **Prefetching**: 基于用户行为的预取
- **Distributed Cache**: Redis Cluster 支持

### 3. 实时优化
- **Online Learning**: 基于用户反馈的实时优化
- **A/B Testing**: 不同搜索算法的对比测试
- **Auto-tuning**: 自动调整搜索参数

---

*本文档持续更新中，最后更新时间: 2025-11-25*