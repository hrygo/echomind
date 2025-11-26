# 搜索增强 API 使用示例

**版本**: v1.2.0  
**更新日期**: 2025-11-26

## 概述

本文档提供搜索增强功能的 API 使用示例和最佳实践。

---

## 基础搜索 (向后兼容)

### 请求

```bash
curl -X GET "http://localhost:8080/api/v1/search?q=项目进展&limit=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 响应

```json
{
  "query": "项目进展",
  "results": [
    {
      "email_id": "123e4567-e89b-12d3-a456-426614174000",
      "subject": "Q4项目进展汇报",
      "sender": "张三 <zhangsan@example.com>",
      "snippet": "本季度项目进展顺利，已完成核心功能开发...",
      "date": "2025-11-20T10:30:00Z",
      "score": 0.92
    }
  ],
  "count": 10
}
```

---

## 启用聚类功能

### 按发件人聚类

```bash
curl -X GET "http://localhost:8080/api/v1/search?q=项目&enable_clustering=true&cluster_type=sender" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**响应**:

```json
{
  "query": "项目",
  "results": [...],
  "count": 15,
  "cluster_type": "sender",
  "clusters": [
    {
      "id": "cluster_sender_1",
      "type": "sender",
      "label": "张三 (5封)",
      "count": 5,
      "results": [...]
    },
    {
      "id": "cluster_sender_2",
      "type": "sender",
      "label": "李四 (3封)",
      "count": 3,
      "results": [...]
    }
  ]
}
```

### 按时间聚类

```bash
curl -X GET "http://localhost:8080/api/v1/search?q=会议&enable_clustering=true&cluster_type=time" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**响应**:

```json
{
  "query": "会议",
  "cluster_type": "time",
  "clusters": [
    {
      "id": "cluster_time_today",
      "type": "time",
      "label": "今天",
      "count": 3,
      "results": [...]
    },
    {
      "id": "cluster_time_this_week",
      "type": "time",
      "label": "本周",
      "count": 5,
      "results": [...]
    },
    {
      "id": "cluster_time_this_month",
      "type": "time",
      "label": "本月",
      "count": 8,
      "results": [...]
    },
    {
      "id": "cluster_time_older",
      "type": "time",
      "label": "更早",
      "count": 4,
      "results": [...]
    }
  ]
}
```

### 按主题聚类

```bash
curl -X GET "http://localhost:8080/api/v1/search?q=技术&enable_clustering=true&cluster_type=topic" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**响应**:

```json
{
  "query": "技术",
  "cluster_type": "topic",
  "clusters": [
    {
      "id": "cluster_topic_1",
      "type": "topic",
      "label": "架构设计 (6封)",
      "count": 6,
      "results": [...]
    },
    {
      "id": "cluster_topic_2",
      "type": "topic",
      "label": "性能优化 (4封)",
      "count": 4,
      "results": [...]
    }
  ]
}
```

---

## 启用 AI 智能摘要

### 基础摘要

```bash
curl -X GET "http://localhost:8080/api/v1/search?q=项目&enable_summary=true" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**响应**:

```json
{
  "query": "项目",
  "results": [...],
  "count": 15,
  "summary": {
    "natural_summary": "找到15封相关邮件，主要来自3位联系人，涵盖项目进展、技术评审和需求变更等主题。",
    "key_topics": [
      "项目进展",
      "技术评审",
      "需求变更",
      "资源分配",
      "风险管理"
    ],
    "important_people": [
      "张三",
      "李四",
      "王五"
    ]
  }
}
```

---

## 组合使用 (聚类 + 摘要)

### 完整增强搜索

```bash
curl -X GET "http://localhost:8080/api/v1/search?q=项目&enable_clustering=true&cluster_type=sender&enable_summary=true&limit=20" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**响应**:

```json
{
  "query": "项目",
  "results": [...],
  "count": 20,
  "cluster_type": "sender",
  "clusters": [
    {
      "id": "cluster_sender_1",
      "type": "sender",
      "label": "张三 (8封)",
      "count": 8,
      "results": [...]
    }
  ],
  "summary": {
    "natural_summary": "找到20封项目相关邮件，主要讨论Q4项目进展、技术架构评审和需求变更等内容。",
    "key_topics": [
      "项目进展",
      "技术架构",
      "需求变更"
    ],
    "important_people": [
      "张三",
      "李四"
    ]
  }
}
```

---

## 高级过滤

### 组合过滤器

```bash
curl -X GET "http://localhost:8080/api/v1/search?\
q=项目&\
sender=zhangsan@example.com&\
start_date=2025-11-01&\
end_date=2025-11-30&\
enable_clustering=true&\
enable_summary=true&\
limit=50" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**支持的过滤参数**:

| 参数 | 类型 | 说明 | 示例 |
|------|------|------|------|
| `q` | string | 搜索查询 (必填) | `项目进展` |
| `sender` | string | 发件人过滤 | `zhangsan@example.com` |
| `context_id` | uuid | 上下文 ID | `123e4567-...` |
| `start_date` | date | 开始日期 | `2025-11-01` |
| `end_date` | date | 结束日期 | `2025-11-30` |
| `limit` | int | 结果数量 (1-100) | `20` |
| `enable_clustering` | boolean | 启用聚类 | `true` |
| `cluster_type` | string | 聚类类型 | `sender/time/topic` |
| `enable_summary` | boolean | 启用摘要 | `true` |

---

## 性能优化建议

### 1. 缓存利用

重复查询会自动使用 Redis 缓存：

```bash
# 第一次查询 - 慢 (~1000ms)
curl -X GET "http://localhost:8080/api/v1/search?q=项目" -H "Authorization: Bearer TOKEN"

# 第二次相同查询 - 快 (~50ms)
curl -X GET "http://localhost:8080/api/v1/search?q=项目" -H "Authorization: Bearer TOKEN"
```

**缓存特性**:
- TTL: 30分钟
- 用户隔离: 每个用户独立缓存
- 自动失效: 邮件更新时清理

### 2. 按需启用增强功能

```bash
# 快速查看 - 不启用增强功能
curl -X GET "http://localhost:8080/api/v1/search?q=项目"

# 深度分析 - 启用完整增强
curl -X GET "http://localhost:8080/api/v1/search?q=项目&enable_clustering=true&enable_summary=true"
```

**性能对比**:
| 场景 | 延迟 | 说明 |
|------|------|------|
| 基础搜索 | ~50ms (缓存) / ~1000ms (首次) | 仅向量搜索 |
| + 聚类 | +20ms | 本地计算 |
| + AI 摘要 | +500ms | AI API 调用 |

### 3. 限制结果数量

```bash
# 快速预览 - 限制10条
curl -X GET "http://localhost:8080/api/v1/search?q=项目&limit=10"

# 完整结果 - 限制100条
curl -X GET "http://localhost:8080/api/v1/search?q=项目&limit=100"
```

---

## 错误处理

### AI 摘要失败自动降级

如果 AI Provider 不可用，系统会自动使用快速摘要：

```json
{
  "query": "项目",
  "results": [...],
  "summary": {
    "natural_summary": "找到15封相关邮件，来自3位不同的发件人。",
    "key_topics": [],
    "important_people": ["张三", "李四", "王五"]
  }
}
```

### 错误响应示例

```json
{
  "error": "search failed",
  "details": "database connection timeout"
}
```

---

## 集成示例

### JavaScript/TypeScript

```typescript
interface SearchOptions {
  query: string;
  enableClustering?: boolean;
  clusterType?: 'sender' | 'time' | 'topic';
  enableSummary?: boolean;
  limit?: number;
}

async function enhancedSearch(options: SearchOptions) {
  const params = new URLSearchParams({
    q: options.query,
    limit: String(options.limit || 10),
  });

  if (options.enableClustering) {
    params.append('enable_clustering', 'true');
    params.append('cluster_type', options.clusterType || 'sender');
  }

  if (options.enableSummary) {
    params.append('enable_summary', 'true');
  }

  const response = await fetch(`/api/v1/search?${params}`, {
    headers: {
      'Authorization': `Bearer ${getToken()}`,
    },
  });

  return response.json();
}

// 使用示例
const results = await enhancedSearch({
  query: '项目',
  enableClustering: true,
  clusterType: 'sender',
  enableSummary: true,
  limit: 20,
});
```

### Python

```python
import requests

def enhanced_search(
    query: str,
    enable_clustering: bool = False,
    cluster_type: str = 'sender',
    enable_summary: bool = False,
    limit: int = 10,
    token: str = None
):
    params = {
        'q': query,
        'limit': limit,
    }
    
    if enable_clustering:
        params['enable_clustering'] = 'true'
        params['cluster_type'] = cluster_type
    
    if enable_summary:
        params['enable_summary'] = 'true'
    
    response = requests.get(
        'http://localhost:8080/api/v1/search',
        params=params,
        headers={'Authorization': f'Bearer {token}'}
    )
    
    return response.json()

# 使用示例
results = enhanced_search(
    query='项目',
    enable_clustering=True,
    cluster_type='sender',
    enable_summary=True,
    limit=20,
    token='YOUR_TOKEN'
)
```

---

## 最佳实践

### 1. 渐进式增强

```bash
# Step 1: 基础搜索
/api/v1/search?q=项目

# Step 2: 添加聚类 (用户点击"按发件人分组")
/api/v1/search?q=项目&enable_clustering=true&cluster_type=sender

# Step 3: 添加摘要 (用户点击"生成摘要")
/api/v1/search?q=项目&enable_clustering=true&enable_summary=true
```

### 2. 用户体验优化

- **快速响应**: 先显示基础搜索结果
- **按需加载**: 用户需要时才调用增强功能
- **缓存友好**: 相同查询自动复用缓存

### 3. 监控和日志

查看 OTel Traces 了解性能：

```bash
# 查看搜索请求的完整追踪
tail -f backend/logs/backend.log | grep trace_id
```

---

## 未来扩展

计划中的功能（v1.3.0）：

- [ ] 前端 UI 组件 (SearchSummaryCard)
- [ ] 聚类视图切换
- [ ] 搜索历史和建议
- [ ] 导出搜索结果

---

**最后更新**: 2025-11-26  
**维护者**: EchoMind Backend Team
