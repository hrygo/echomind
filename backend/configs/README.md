# 配置文件说明

本文档说明如何配置 EchoMind 后端服务。

## 文件说明

- `config.yaml` - **实际配置文件** (包含真实API密钥，不提交到版本控制)
- `config.example.yaml` - **示例配置文件** (已脱敏，可安全分享)

## 首次配置步骤

### 1. 复制示例配置
```bash
cp config.example.yaml config.yaml
```

### 2. 替换占位符为真实值

#### 必需配置
- `server.jwt.secret`: JWT密钥 (推荐64字符随机字符串)
- `security.encryption_key`: 加密密钥 (64字符十六进制字符串)

#### AI服务配置 (至少配置一个)
根据你的需求选择AI提供商：

**选项A: SiliconFlow (推荐 - 性价比高)**
```yaml
active_services:
  chat: "deepseek"
  embedding: "siliconflow"

providers:
  deepseek:
    settings:
      api_key: "sk-your-actual-deepseek-key"

  siliconflow:
    settings:
      api_key: "sk-your-actual-siliconflow-key"
```

**选项B: OpenAI (标准选择)**
```yaml
active_services:
  chat: "openai_small"
  embedding: "openai_small"

providers:
  openai_small:
    settings:
      api_key: "sk-your-actual-openai-key"
```

**选项C: Google Gemini (多模态)**
```yaml
active_services:
  chat: "gemini_flash"
  embedding: "gemini_flash"

providers:
  gemini_flash:
    settings:
      api_key: "your-actual-gemini-key"
```

### 3. 生成密钥

#### JWT密钥生成
```bash
# 生成64字符JWT密钥
openssl rand -hex 32
```

#### 加密密钥生成
```bash
# 生成64字符十六进制加密密钥
openssl rand -hex 32
```

## 嵌入模型维度配置

每个AI提供商都有特定的嵌入维度，必须与数据库schema匹配：

| 提供商 | 模型 | 默认维度 | 配置值 |
|--------|------|----------|--------|
| SiliconFlow | BGE-M3 | 1024 | `embedding_dimensions: 1024` |
| OpenAI | text-embedding-3-small | 1536 | `embedding_dimensions: 1536` |
| Google Gemini | text-embedding-004 | 768 | `embedding_dimensions: 768` |
| Ollama | nomic-embed-text | 768 | `embedding_dimensions: 768` |

**重要**: 数据库schema现在支持最大1536维度，包含自动转换：
- 矢量存储为 `vector(1536)` 以支持OpenAI标准
- 数据库层自动处理填充/截断至正确维度
- Dimensions字段跟踪原始向量维度
- 无需手动迁移数据

## 数据库配置

### PostgreSQL + pgvector
```yaml
database:
  dsn: "host=localhost user=user password=password dbname=echomind_db port=5432 sslmode=disable"
```

确保安装了pgvector扩展：
```sql
CREATE EXTENSION IF NOT EXISTS vector;
```

### Redis配置
```yaml
redis:
  addr: "localhost:6380"
  password: ""
  db: 0
```

## 生产环境配置

### 安全配置
```yaml
server:
  environment: production

security:
  encryption_key: "your-production-64-char-hex-key"
```

### JWT配置
```yaml
jwt:
  secret: "your-production-jwt-secret"
  expiration_hours: 24  # 生产环境建议更短
```

### 数据库连接
```yaml
database:
  dsn: "host=your-db-host user=your-user password=your-password dbname=echomind_prod port=5432 sslmode=require"
```

## 环境变量支持

你也可以通过环境变量覆盖配置：

```bash
export ECHOMIND_DB_DSN="your-production-dsn"
export ECHOMIND_JWT_SECRET="your-jwt-secret"
export ECHOMIND_ENCRYPTION_KEY="your-encryption-key"
```

## 常见问题

### Q: 如何切换AI提供商？
A: 修改 `active_services` 部分即可：

```yaml
# 从SiliconFlow切换到OpenAI
active_services:
  chat: "openai_small"
  embedding: "openai_small"
```

### Q: 嵌入维度错误怎么办？
A: 确保配置的 `embedding_dimensions` 与提供商实际输出匹配，且与数据库schema一致。

### Q: 如何获取API密钥？
A:
- **DeepSeek**: https://platform.deepseek.com/api_keys
- **SiliconFlow**: https://cloud.siliconflow.cn/key
- **OpenAI**: https://platform.openai.com/api-keys
- **Google Gemini**: https://makersuite.google.com/app/apikey

## 配置验证

配置完成后，运行以下命令验证：
```bash
# 检查YAML语法
python3 -c "import yaml; yaml.safe_load(open('config.yaml'))"

# 启动服务
go run cmd/main.go
```

## 安全提醒

- ✅ `config.example.yaml` 已脱敏，可安全分享
- ❌ `config.yaml` 包含敏感信息，请勿提交到版本控制
- 🔒 确保API密钥安全存储
- 🔄 定期更换JWT密钥
- 🛡️ 生产环境使用HTTPS连接