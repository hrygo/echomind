# 数据库迁移指南

## 问题描述

当执行 `make reindex` 时出现错误：
```
ERROR: expected 768 dimensions, not 1536 (SQLSTATE 22000)
```

这是因为数据库表 `email_embeddings` 仍然使用旧的向量维度（768/1024），而新的代码期望 1536 维度。

## 解决方案

### 方法一：使用 Makefile 自动迁移（推荐）

```bash
# 执行数据库迁移
make migrate-db

# 重启后端服务
make stop-apps
make run-backend

# 重新生成所有嵌入向量
make reindex
```

### 方法二：手动执行 SQL

1. **连接数据库**:
```bash
make db-shell
```

2. **执行迁移脚本**:
```sql
-- 备份现有数据
CREATE TABLE email_embeddings_backup_20251125 AS SELECT * FROM email_embeddings;

-- 删除旧表
DROP TABLE IF EXISTS email_embeddings CASCADE;

-- 退出数据库
\q
```

3. **重启应用**: 应用会自动创建正确的新表结构

4. **重新生成嵌入**:
```bash
make reindex
```

## 迁移后的表结构

迁移后的 `email_embeddings` 表结构：

```sql
CREATE TABLE email_embeddings (
    id SERIAL PRIMARY KEY,
    email_id UUID NOT NULL,
    content TEXT,
    vector vector(1536),      -- 支持 1536 维度
    dimensions INTEGER NOT NULL DEFAULT 1024,  -- 跟踪实际维度
    created_at TIMESTAMP
);
```

## 验证迁移

执行以下 SQL 验证迁移是否成功：

```sql
-- 检查表结构
SELECT column_name, data_type
FROM information_schema.columns
WHERE table_name = 'email_embeddings';

-- 检查索引
SELECT indexname, indexdef
FROM pg_indexes
WHERE tablename = 'email_embeddings';
```

## 注意事项

⚠️ **警告**: 迁移过程会删除现有的嵌入向量，需要重新生成。

- **数据备份**: 迁移前会自动创建备份表 `email_embeddings_backup_YYYYMMDD`
- **重新索引**: 迁移后必须执行 `make reindex` 重新生成所有嵌入向量
- **服务重启**: 迁移后需要重启后端服务以使用新的表结构

## 恢复方案

如果迁移失败，可以使用备份表恢复：

```sql
-- 删除新表
DROP TABLE IF EXISTS email_embeddings;

-- 恢复备份表
ALTER TABLE email_embeddings_backup_20251125 RENAME TO email_embeddings;
```

## 技术细节

- **向量类型**: `vector(1536)` 支持 OpenAI 的最大维度
- **自动转换**: 应用层会自动处理 768/1024 维度到 1536 的填充
- **向后兼容**: 新结构支持所有嵌入供应商（Gemini 768、SiliconFlow 1024、OpenAI 1536）

## 相关文件

- `backend/internal/model/embedding.go` - 数据模型定义
- `backend/internal/service/search.go` - 搜索服务实现
- `docs/architecture.md` - 技术架构文档
- `docs/vector-search-guide.md` - 向量搜索指南