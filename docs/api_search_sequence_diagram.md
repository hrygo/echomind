# EchoMind 邮件处理系统时序图

## 📋 目录

### 一、系统概览
- [1.1 核心流程说明](#11-核心流程说明)
- [1.2 共享技术栈](#12-共享技术栈)

### 二、核心业务流程
- [2.1 邮件搜索流程](#21-邮件搜索流程)
  - [2.1.1 搜索流程时序图](#211-搜索流程时序图)
  - [2.1.2 关键组件详细说明](#212-关键组件详细说明)
- [2.2 邮件同步流程](#22-邮件同步流程)
  - [2.2.1 同步流程时序图](#221-同步流程时序图)
  - [2.2.2 同步流程关键组件](#222-同步流程关键组件)
- [2.3 Reindex工具流程](#23-reindex工具流程)
  - [2.3.1 Reindex工具时序图](#231-reindex工具时序图)
  - [2.3.2 Reindex工具关键特性](#232-reindex工具关键特性)

### 三、技术架构组件
- [3.1 向量嵌入处理流程](#31-向量嵌入处理流程)
- [3.2 AI Provider架构](#32-ai-provider架构)
- [3.3 数据库模式设计](#33-数据库模式设计)
- [3.4 事件驱动架构](#34-事件驱动架构)

### 四、系统增强方案
- [4.1 当前搜索流程分析](#41-当前搜索流程分析)
- [4.2 增强搜索流程时序图](#42-增强搜索流程时序图)
- [4.3 建议实施的搜索增强功能](#43-建议实施的搜索增强功能)

### 五、监控与运维
- [5.1 关键性能指标 (KPIs)](#51-关键性能指标-kpis)
- [5.2 错误处理流程](#52-错误处理流程)
- [5.3 运维最佳实践](#53-运维最佳实践)

---

## 一、系统概览

EchoMind是一个基于AI的智能邮件处理系统，主要由三个核心流程组成：

```mermaid
flowchart TB
    subgraph "EchoMind 智能邮件处理系统"
        direction LR

        subgraph "数据流"
            A[外部邮件] --> B[邮件同步]
            B --> C[(向量数据库)]
            C --> D[智能搜索]
        end

        subgraph "AI处理"
            E[嵌入生成] --> F[邮件分析]
            F --> G[向量搜索]
        end

        subgraph "管理工具"
            H[Reindex工具] --> C
        end
    end

    style B fill:#e1f5fe
    style D fill:#f3e5f5
    style H fill:#fff3e0
    style C fill:#e8f5e8
```

### 1.1 核心流程说明

| 流程 | 功能描述 | 触发方式 | 主要输出 |
|------|----------|----------|----------|
| **邮件搜索** | 基于AI的智能邮件检索 | 用户搜索请求 | 相关邮件列表及AI分析 |
| **邮件同步** | 从邮箱服务器获取新邮件 | 定时/手动同步 | 邮件数据 + 向量嵌入 |
| **Reindex工具** | 重建现有邮件的向量索引 | 管理员执行 | 更新的向量数据库 |

### 1.2 共享技术栈

- **向量数据库**: PostgreSQL + pgvector
- **AI嵌入**: 多Provider支持 (OpenAI, Gemini, SiliconFlow等)
- **异步处理**: Asynq任务队列
- **事件驱动**: 事件总线架构
- **缓存**: Redis (建议实施)

---

## 二、核心业务流程

### 2.1 邮件搜索流程

邮件搜索是系统的核心用户接口，提供基于AI的语义搜索功能，支持自然语言查询和智能结果排序。

#### 2.1.1 搜索流程时序图

```mermaid
sequenceDiagram
    participant Client as 客户端
    participant Router as 路由层
    participant AuthMW as 认证中间件
    participant SearchHandler as 搜索处理器
    participant SearchService as 搜索服务
    participant AIProvider as AI嵌入提供者
    participant Config as 配置系统
    participant DB as 数据库
    participant OpenAI as OpenAI API

    %% 1. 请求入口与认证
    Client->>Router: GET /api/v1/search?q=project&limit=5
    activate Router
    Router->>AuthMW: JWT认证中间件
    activate AuthMW
    AuthMW->>AuthMW: 验证Token并提取用户ID
    AuthMW->>Router: 设置用户ID到Context
    deactivate AuthMW

    %% 2. 搜索处理
    Router->>SearchHandler: Search(ctx, userID, query, filters, limit)
    activate SearchHandler
    SearchHandler->>SearchHandler: 解析查询参数
    Note right of SearchHandler: - q: "project"<br/>- sender: 可选<br/>- context_id: 可选<br/>- start_date/end_date: 可选<br/>- limit: 5

    %% 3. 调用搜索服务
    SearchHandler->>SearchService: Search(ctx, userID, query, filters, limit)
    activate SearchService

    %% 4. 嵌入查询生成
    SearchService->>AIProvider: Embed(ctx, query="project")
    activate AIProvider
    AIProvider->>Config: 获取嵌入维度配置
    activate Config
    Config->>Config: 读取active_services.embedding
    Note right of Config: 当前配置: "siliconflow"
    Config->>Config: 查找providers.siliconflow配置
    Config->>Config: 读取embedding_dimensions: 1024
    Config->>AIProvider: 返回维度配置
    deactivate Config

    %% 5. AI向量生成
    AIProvider->>AIProvider: 验证嵌入模型配置
    Note right of AIProvider: embedding_model: "Pro/BAAI/bge-m3"<br/>dimensions: 1024
    AIProvider->>OpenAI: CreateEmbeddings(ctx, req)
    activate OpenAI
    Note right of OpenAI: Request:<br/>- Model: "Pro/BAAI/bge-m3"<br/>- Input: "project"<br/>- Dimensions: 1024
    OpenAI->>AIProvider: 返回嵌入向量 [1024维]
    deactivate OpenAI

    %% 6. 向量验证
    AIProvider->>AIProvider: 验证向量维度
    Note right of AIProvider: 实际维度: 1024<br/>配置维度: 1024 ✓
    AIProvider->>SearchService: 返回查询向量 [1024维]
    deactivate AIProvider

    %% 7. 数据库向量搜索
    SearchService->>DB: 执行向量相似度搜索
    activate DB
    Note right of DB: SQL:<br/>SELECT e.id, e.subject, ee.content,<br/>1 - (ee.vector <=> ?) as score<br/>FROM email_embeddings ee<br/>JOIN emails e ON e.id = ee.email_id<br/>WHERE e.user_id = ?<br/>ORDER BY ee.vector <=> ?<br/>LIMIT ?

    DB->>DB: pgvector向量比较和相似度计算
    DB->>SearchService: 返回搜索结果列表
    deactivate DB

    %% 8. 结果格式化与响应
    SearchService->>SearchService: 格式化搜索结果
    Note right of SearchService: SearchResult[]:<br/>- EmailID, Subject, Snippet<br/>- Sender, Date, Score (0-1)

    SearchService->>SearchHandler: 返回格式化结果
    deactivate SearchService
    SearchHandler->>Router: JSON响应
    deactivate SearchHandler
    Router->>Client: HTTP 200 OK
    deactivate Router

    Note over Client,OpenAI: 搜索流程完成
```

#### 2.1.2 关键组件详细说明

##### 2.1.2.1 AI嵌入维度配置系统

```mermaid
flowchart TD
    A[搜索请求] --> B[AIProvider.Embed]
    B --> C[读取配置系统]
    C --> D[获取active_services.embedding]
    D --> E[查找对应Provider配置]
    E --> F[读取embedding_dimensions]
    F --> G{维度配置存在?}
    G -->|是| H[使用配置维度]
    G -->|否| I[使用默认维度1024]
    H --> J[验证嵌入模型兼容性]
    I --> J
    J --> K[设置API请求维度参数]
    K --> L[调用外部AI服务]
```

##### 2.1.2.2 AI Provider配置映射表

| Provider | 嵌入模型 | 配置维度 | 模型原生维度 | 处理方式 |
|----------|----------|----------|--------------|----------|
| **siliconflow** | Pro/BAAI/bge-m3 | 1024 | 1024 | 直接使用 |
| **openai_small** | text-embedding-3-small | 1536 | 1536 | 直接使用 |
| **gemini_flash** | text-embedding-004 | 768 | 768 | 直接使用 |
| **local_ollama** | nomic-embed-text | 768 | 768 | 直接使用 |
| **mock** | - | 1024 | 1024 | 模拟生成 |

##### 2.1.2.3 向量维度验证机制

```go
// backend/internal/model/embedding.go:41-69
func (e *EmailEmbedding) validateAndConvertVector(tx *gorm.DB) error {
    vectorSlice := e.Vector.Slice()
    actualDimensions := len(vectorSlice)
    e.Dimensions = actualDimensions

    maxDimensions := 1536 // OpenAI最大标准维度

    // 超过最大维度则截断
    if actualDimensions > maxDimensions {
        truncatedSlice := vectorSlice[:maxDimensions]
        e.Vector = pgvector.NewVector(truncatedSlice)
        e.Dimensions = maxDimensions
    }

    // 小于最大维度则用零填充
    if actualDimensions < maxDimensions {
        paddedVector := make([]float32, maxDimensions)
        copy(paddedVector, vectorSlice)
        e.Vector = pgvector.NewVector(paddedVector)
    }

    return nil
}
```

##### 2.1.2.4 数据库搜索算法

**核心SQL查询**:
```sql
SELECT
    e.id as email_id,
    e.subject,
    ee.content as snippet,
    e.sender,
    e.date,
    1 - (ee.vector <=> ?) as score  -- 向量相似度计算
FROM email_embeddings ee
JOIN emails e ON e.id = ee.email_id
WHERE e.user_id = ?
ORDER BY ee.vector <=> ?  -- 按距离排序
LIMIT ?
```

**搜索步骤**:
1. **距离计算**: 使用 pgvector 的 `<=>` 操作符计算欧几里得距离
2. **相似度转换**: `1 - 距离` 得到相似度分数 (0-1之间，1为最相似)
3. **排序优化**: 使用 HNSW 索引加速近似最近邻搜索
4. **用户过滤**: 应用用户权限和结果数量限制

---

### 2.2 邮件同步流程

邮件同步流程负责从用户的邮箱服务器获取新邮件，进行AI分析处理，并生成向量嵌入以支持搜索功能。采用事件驱动的异步架构。

#### 2.2.1 同步流程时序图

```mermaid
sequenceDiagram
    participant Client as 客户端
    participant SyncHandler as 同步处理器
    participant SyncService as 同步服务
    participant EmailAccount as 邮箱账户
    participant IMAPConnector as IMAP连接器
    participant IMAPServer as IMAP服务器
    participant DB as 数据库
    participant EventBus as 事件总线
    participant TaskQueue as 任务队列
    participant AIService as AI服务
    participant EmbeddingService as 嵌入服务

    %% 1. 同步请求
    Client->>SyncHandler: POST /api/v1/sync/emails
    activate SyncHandler
    SyncHandler->>SyncService: SyncEmails(ctx, userID, teamID, orgID)
    activate SyncService

    %% 2. 获取邮箱配置
    SyncService->>EmailAccount: 获取用户邮箱配置
    activate EmailAccount
    EmailAccount->>EmailAccount: 读取账户信息
    Note right of EmailAccount: IMAP服务器、加密密码等
    EmailAccount->>SyncService: 返回账户配置
    deactivate EmailAccount

    %% 3. 建立IMAP连接
    SyncService->>IMAPConnector: NewConnector(account)
    activate IMAPConnector
    IMAPConnector->>IMAPConnector: 解密密码
    IMAPConnector->>IMAPServer: 建立TLS连接
    activate IMAPServer
    IMAPConnector->>IMAPServer: IMAP登录
    IMAPServer-->>IMAPConnector: 登录成功
    deactivate IMAPServer
    IMAPConnector->>SyncService: 连接就绪
    deactivate IMAPConnector

    %% 4. 获取邮件数据
    SyncService->>IMAPConnector: FetchEmails(lastSyncTime)
    activate IMAPConnector
    IMAPConnector->>IMAPServer: SELECT INBOX
    IMAPConnector->>IMAPServer: SEARCH SINCE lastSyncTime
    IMAPServer-->>IMAPConnector: 返回邮件UID列表
    IMAPConnector->>IMAPConnector: 获取最新10封邮件
    IMAPConnector->>IMAPConnector: 提取元数据和正文
    Note right of IMAPConnector: 提取主题、发件人、日期<br/>Message-ID、正文(TEXT/HTML)
    IMAPConnector->>SyncService: 返回邮件数据
    deactivate IMAPConnector

    %% 5. 邮件存储和事件发布
    loop 每封邮件处理
        SyncService->>DB: 检查Message-ID是否存在
        alt 邮件不存在
            SyncService->>DB: INSERT INTO emails
            activate DB
            DB->>DB: 保存邮件基础信息
            DB-->>SyncService: 保存成功
            deactivate DB

            %% 发布同步事件
            SyncService->>EventBus: Publish(EmailSyncedEvent)
            activate EventBus
            EventBus->>EventBus: 触发事件监听器

            %% 创建AI分析任务
            EventBus->>TaskQueue: Enqueue(EmailAnalyzeTask)
            activate TaskQueue
            TaskQueue->>TaskQueue: 添加到异步队列
            TaskQueue-->>EventBus: 任务已入队
            deactivate TaskQueue
            deactivate EventBus
        else 邮件已存在
            SyncService->>SyncService: 跳过重复邮件
        end
    end

    %% 6. 更新同步状态
    SyncService->>EmailAccount: 更新LastSyncAt
    activate EmailAccount
    EmailAccount->>EmailAccount: 记录最后同步时间
    EmailAccount-->>SyncService: 更新完成
    deactivate EmailAccount

    SyncService->>SyncHandler: 同步完成
    deactivate SyncService
    SyncHandler->>Client: HTTP 200 OK
    deactivate SyncHandler

    Note over Client,EmbeddingService: 同步请求完成，后台任务继续处理

    %% 7. 异步AI分析流程
    Note over TaskQueue,AIService: 后台异步处理流程
    TaskQueue->>AIService: HandleEmailAnalyzeTask
    activate AIService

    %% 垃圾邮件检测
    AIService->>AIService: SpamDetection(rules)
    Note right of AIService: 基于规则检测垃圾邮件<br/>- 发件人黑名单<br/>- 可疑关键词<br/>- 发送频率异常

    %% AI分析处理
    AIService->>AIService: 生成邮件摘要
    AIService->>AIService: 分类和情感分析
    AIService->>AIService: 紧急程度评估
    AIService->>AIService: 智能上下文匹配
    AIService->>AIService: 提取待办事项

    %% 更新分析结果
    AIService->>DB: UPDATE emails SET ai_analysis
    activate DB
    DB-->>AIService: 更新成功
    deactivate DB

    %% 生成向量嵌入
    AIService->>EmbeddingService: GenerateEmbedding(email_content)
    activate EmbeddingService
    EmbeddingService->>EmbeddingService: 文本分块处理
    EmbeddingService->>EmbeddingService: 调用AI嵌入API
    EmbeddingService->>EmbeddingService: 向量维度验证和转换

    %% 保存向量嵌入
    EmbeddingService->>DB: INSERT INTO email_embeddings
    activate DB
    DB->>DB: pgvector向量存储和索引更新
    DB-->>EmbeddingService: 嵌入保存成功
    deactivate DB
    EmbeddingService-->>AIService: 嵌入生成完成
    deactivate EmbeddingService

    AIService->>TaskQueue: 标记任务完成
    deactivate AIService
    deactivate TaskQueue

    Note over Client,EmbeddingService: 完整邮件同步和AI处理流程完成
```

#### 2.2.2 同步流程关键组件

##### 2.2.2.1 IMAP连接器配置

| 配置项 | 说明 | 示例 |
|--------|------|------|
| **Server** | IMAP服务器地址 | "imap.gmail.com:993" |
| **Username** | 邮箱地址 | "user@gmail.com" |
| **Password** | 加密存储的密码 | AES-256加密 |
| **LastSyncAt** | 最后同步时间 | 2025-01-15 10:30:00 |
| **SyncLimit** | 同步邮件数量限制 | 默认10封 |

##### 2.2.2.2 同步过滤规则

```go
// 同步过滤逻辑
func shouldSyncEmail(email *Email, lastSync time.Time) bool {
    return email.Date.After(lastSync) &&           // 时间过滤
           !isDuplicate(email.MessageID) &&         // 重复检测
           !isSpam(email.Sender, email.Subject) &&  // 垃圾邮件过滤
           email.BodyText != ""                     // 内容非空
}
```

##### 2.2.2.3 AI分析维度

| 分析类型 | 功能说明 | 输出格式 | 应用场景 |
|----------|----------|----------|----------|
| **摘要生成** | 提取邮件核心内容 | 50-100字摘要 | 快速浏览 |
| **分类** | 邮件类别识别 | work/personal/newsletter等 | 自动分类 |
| **情感分析** | 情感倾向判断 | positive/neutral/negative | 情绪追踪 |
| **紧急度** | 重要性评估 | high/medium/low | 优先级处理 |
| **智能上下文** | 关联项目/客户 | project_context/client_context | 业务关联 |

---

### 2.3 Reindex工具流程

Reindex工具用于重建现有邮件的向量嵌入，主要应用场景：
- 更新AI嵌入模型后重新生成向量
- 修复损坏或不完整的向量数据
- 调整向量维度配置
- 系统迁移后的数据重建

#### 2.3.1 Reindex工具时序图

```mermaid
sequenceDiagram
    participant Admin as 管理员
    participant ReindexCLI as Reindex CLI
    participant Container as 应用容器
    participant DB as 数据库
    participant SearchService as 搜索服务
    participant AIService as AI嵌入服务
    participant EmbeddingModel as 嵌入模型
    participant PGVector as pgvector扩展

    %% 1. 工具启动和初始化
    Admin->>ReindexCLI: ./reindex --config config.yaml
    activate ReindexCLI
    ReindexCLI->>ReindexCLI: 解析命令行参数
    ReindexCLI->>Container: app.NewContainer(configPath)
    activate Container

    Container->>Container: 加载配置文件
    Container->>Container: 初始化数据库连接
    Container->>Container: 初始化AI服务提供者
    Container->>Container: 初始化搜索服务
    Container-->>ReindexCLI: 容器初始化完成
    deactivate Container

    %% 2. 批量获取邮件数据
    ReindexCLI->>DB: SELECT id, subject, snippet, body_text FROM emails
    activate DB
    DB->>DB: 查询所有邮件基础字段
    DB->>DB: 按创建时间排序
    DB-->>ReindexCLI: 返回邮件列表 (N封邮件)
    deactivate DB

    ReindexCLI->>ReindexCLI: 记录开始日志
    Note right of ReindexCLI: 开始重建索引<br/>邮件总数: N

    %% 3. 逐个处理邮件
    loop 每封邮件处理
        ReindexCLI->>ReindexCLI: 获取下一封邮件
        ReindexCLI->>SearchService: GenerateAndSaveEmbedding(ctx, email)
        activate SearchService

        %% 删除旧嵌入数据
        SearchService->>DB: DELETE FROM email_embeddings WHERE email_id = ?
        activate DB
        DB->>DB: 删除现有向量数据
        DB-->>SearchService: 删除完成
        deactivate DB

        %% 文本预处理
        SearchService->>SearchService: 合并邮件文本内容
        Note right of SearchService: content = subject + "\\n" + snippet + "\\n" + body_text

        SearchService->>SearchService: 文本清理和标准化
        Note right of SearchService: - 去除HTML标签<br/>- 移除多余空白<br/>- 编码标准化

        %% 文本分块
        SearchService->>SearchService: 分块处理长文本
        Note right of SearchService: chunkSize := 1000字符<br/>maxChunks := 10块<br/>避免超出API限制

        %% 生成向量嵌入
        SearchService->>AIService: Embed(content_chunks)
        activate AIService

        AIService->>EmbeddingModel: 批量嵌入生成
        activate EmbeddingModel
        EmbeddingModel->>EmbeddingModel: 查询模型配置
        Note right of EmbeddingModel: 当前模型: "siliconflow/bge-m3"<br/>向量维度: 1024

        EmbeddingModel->>EmbeddingModel: API调用生成向量
        EmbeddingModel-->>AIService: 返回嵌入向量 [1024维]
        deactivate EmbeddingModel

        AIService->>AIService: 向量聚合和验证
        AIService->>AIService: 维度标准化
        AIService-->>SearchService: 返回最终嵌入向量
        deactivate AIService

        %% 创建嵌入记录
        SearchService->>SearchService: 创建EmailEmbedding对象
        Note right of SearchService: EmailEmbedding{<br/>  EmailID: email.ID,<br/>  Vector: vector,<br/>  Dimensions: 1024,<br/>  Content: content<br/>}

        %% 保存到向量数据库
        SearchService->>DB: INSERT INTO email_embeddings
        activate DB
        DB->>PGVector: 插入向量数据
        activate PGVector
        PGVector->>PGVector: pgvector向量存储
        PGVector->>PGVector: 更新HNSW索引
        PGVector-->>DB: 存储完成
        deactivate PGVector

        DB->>DB: 验证向量完整性
        DB->>DB: 记录元数据
        DB-->>SearchService: 保存成功
        deactivate DB

        SearchService-->>ReindexCLI: 邮件处理完成
        deactivate SearchService

        %% 进度统计
        ReindexCLI->>ReindexCLI: 更新处理统计
        Note right of ReindexCLI: success++ 或 failed++
    end

    %% 4. 生成处理报告和清理
    ReindexCLI->>ReindexCLI: 输出最终统计
    Note right of ReindexCLI: Reindex完成<br/>成功: X 封<br/>失败: Y 封<br/>总耗时: Z 分钟

    ReindexCLI->>Container: Close() 清理资源
    activate Container
    Container->>Container: 关闭数据库连接
    Container->>Container: 关闭AI服务连接
    Container-->>ReindexCLI: 资源清理完成
    deactivate Container

    ReindexCLI->>Admin: 退出程序
    deactivate ReindexCLI

    Note over Admin,PGVector: Reindex工具执行完成
```

#### 2.3.2 Reindex工具关键特性

##### 2.3.2.1 性能优化策略

| 优化项 | 说明 | 效果 |
|--------|------|------|
| **批量查询** | 一次查询所有邮件基础字段 | 减少数据库往返次数 |
| **文本分块** | 长邮件分块处理避免API限制 | 提高处理成功率 |
| **连接池** | 复用数据库和AI服务连接 | 降低连接开销 |
| **进度追踪** | 实时记录处理进度 | 便于监控和调试 |

##### 2.3.2.2 错误处理机制

```go
// 错误处理逻辑
func (cli *ReindexCLI) processEmail(email *model.Email) error {
    defer func() {
        if r := recover(); r != nil {
            cli.Logger.Error("邮件处理异常",
                logger.String("email_id", email.ID.String()),
                logger.Any("panic", r))
            cli.failed++
        }
    }()

    if err := cli.SearchService.GenerateAndSaveEmbedding(ctx, email); err != nil {
        cli.Logger.Warn("邮件重建失败",
            logger.String("email_id", email.ID.String()),
            logger.Error(err))
        cli.failed++
        return err
    }

    cli.success++
    return nil
}
```

##### 2.3.2.3 配置参数说明

| 参数 | 默认值 | 说明 | 影响 |
|------|--------|------|------|
| **ChunkSize** | 1000 | 文本分块大小(字符) | API调用稳定性 |
| **MaxChunks** | 10 | 最大分块数量 | 处理效果和成本 |
| **BatchSize** | 50 | 批量处理大小(预留) | 未来性能优化 |
| **LogLevel** | info | 日志记录级别 | 调试便利性 |

##### 2.3.2.4 数据完整性保障

- **事务处理**: 每封邮件的嵌入更新使用独立事务
- **向量验证**: 检查向量维度和数值范围
- **索引维护**: 自动更新pgvector索引
- **备份保护**: 删除旧嵌入前保存备份

---

## 三、技术架构组件

### 3.1 向量嵌入处理流程

```go
// 标准嵌入生成流程
func (s *EmbeddingService) GenerateEmbedding(content string) ([]float32, error) {
    // 1. 文本预处理和分块
    chunks := s.chunkText(content, MaxChunkSize)

    // 2. 批量生成嵌入
    embeddings := make([][]float32, len(chunks))
    for i, chunk := range chunks {
        embeddings[i] = s.aiProvider.Embed(chunk)
    }

    // 3. 聚合多块嵌入
    finalEmbedding := s.aggregateEmbeddings(embeddings)

    // 4. 维度验证和标准化
    return s.validateAndNormalizeVector(finalEmbedding)
}
```

### 3.2 AI Provider架构

```mermaid
flowchart LR
    subgraph "AI Provider 层"
        A[统一接口] --> B[OpenAI Provider]
        A --> C[Gemini Provider]
        A --> D[SiliconFlow Provider]
        A --> E[Ollama Provider]
        A --> F[Mock Provider]
    end

    subgraph "配置管理层"
        G[动态配置] --> H[模型选择]
        G --> I[维度配置]
        G --> J[API密钥管理]
    end

    subgraph "缓存层"
        K[嵌入缓存] --> L[TTL管理]
        K --> M[LRU淘汰策略]
    end

    A --> G
    A --> K
```

### 3.3 数据库模式设计

```sql
-- 邮件主表
CREATE TABLE emails (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    subject TEXT,
    body_text TEXT,
    sender TEXT,
    date TIMESTAMP WITH TIME ZONE,
    message_id TEXT UNIQUE,
    ai_analysis JSONB,  -- AI分析结果
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 向量嵌入表
CREATE TABLE email_embeddings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email_id UUID NOT NULL REFERENCES emails(id) ON DELETE CASCADE,
    content TEXT NOT NULL,        -- 用于嵌入的文本内容
    vector vector(1536),          -- pgvector向量
    dimensions INTEGER NOT NULL,  -- 实际维度
    model_version TEXT,           -- 嵌入模型版本
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 向量索引
CREATE INDEX idx_email_embeddings_vector ON email_embeddings
USING hnsw (vector vector_cosine_ops);

-- 用户索引
CREATE INDEX idx_emails_user_id ON emails(user_id);
CREATE INDEX idx_email_embeddings_email_id ON email_embeddings(email_id);
```

### 3.4 事件驱动架构

```mermaid
flowchart TD
    A[邮件同步完成] --> B[事件总线]
    B --> C[AI分析任务]
    B --> D[通知任务]
    B --> E[统计任务]

    C --> F[嵌入生成任务]
    F --> G[向量索引更新]

    H[邮件更新] --> B
    I[邮件删除] --> B

    B --> J[任务队列]
    J --> K[后台处理器]
```

---

## 四、系统增强方案

### 4.1 当前搜索流程分析

基于代码分析，当前搜索流程相对简单，仅包含向量相似度搜索，缺少后处理环节。以下是建议的增强流程：

### 4.2 增强搜索流程时序图

```mermaid
sequenceDiagram
    participant Client as 客户端
    participant SearchHandler as 搜索处理器
    participant SearchService as 搜索服务
    participant AIService as AI服务
    participant SummaryService as 摘要服务
    participant Cache as 缓存层

    Note over Client,Cache: 基础向量搜索完成(参考上方时序图 1-13)

    %% 14. 搜索结果后处理 (新增)
    activate SearchService
    SearchService->>AIService: 分析搜索结果集合
    activate AIService
    AIService->>AIService: 结果聚类分析
    Note right of AIService: 按主题、发件人、时间聚类

    AIService->>AIService: 检测结果模式
    Note right of AIService: 识别紧急邮件、重要发件人、趋势主题

    AIService->>AIService: 个性化排序
    Note right of AIService: 基于用户偏好调整排序

    AIService-->>SearchService: 返回增强结果
    deactivate AIService

    %% 15. 智能摘要生成 (新增)
    SearchService->>SummaryService: 生成搜索结果摘要
    activate SummaryService
    SummaryService->>SummaryService: 提取关键信息
    Note right of SummaryService: - 识别3-5个核心主题<br/>- 统计邮件数量分布<br/>- 提取重要联系人

    SummaryService->>SummaryService: 生成自然语言摘要
    Note right of SummaryService: "找到15封相关邮件，主要来自项目团队，<br/>包含3个讨论主题，其中有2封紧急邮件需要关注"

    SummaryService-->>SearchService: 返回智能摘要
    deactivate SummaryService

    %% 16. 结果缓存 (新增)
    SearchService->>Cache: 缓存搜索结果
    activate Cache
    Cache->>Cache: 生成缓存键 (query_hash + user_id)
    Cache->>Cache: 存储结果和摘要 (TTL: 30分钟)
    Cache-->>SearchService: 缓存完成
    deactivate Cache

    %% 17. 构建增强响应
    SearchService->>SearchService: 构建最终响应
    Note right of SearchService: SearchResultResponse:<br/>- Results: enhanced_results[]<br/>- Summary: ai_summary<br/>- Insights: search_insights<br/>- TotalCount: total_count

    %% 18. 返回增强响应
    SearchService->>SearchHandler: 返回增强搜索结果
    deactivate SearchService

    activate SearchHandler
    SearchHandler->>Client: HTTP 200 OK + 增强数据
    deactivate SearchHandler

    Note over Client,Cache: 增强搜索流程完成
```

### 4.3 建议实施的搜索增强功能

#### 4.3.1 结果聚类分析
- **主题聚类**: 按邮件内容相似性分组
- **发件人聚类**: 按发件人/部门分组显示
- **时间聚类**: 按时间周期(今天/本周/本月)分组

#### 4.3.2 智能摘要服务

```go
type SearchSummary struct {
    TotalCount       int                    `json:"total_count"`
    KeyTopics        []string              `json:"key_topics"`
    ImportantPeople  []PersonSummary       `json:"important_people"`
    UrgentCount      int                   `json:"urgent_count"`
    TimeDistribution map[string]int       `json:"time_distribution"`
    NaturalSummary   string                `json:"natural_summary"`
}

type PersonSummary struct {
    Name   string `json:"name"`
    Email  string `json:"email"`
    Count  int    `json:"count"`
    Urgent int    `json:"urgent"`
}
```

#### 4.3.3 个性化搜索
- **用户偏好学习**: 记录用户点击和查看模式
- **发件人权重**: 重要联系人优先显示
- **时间权重**: 近期邮件适当提权
- **上下文相关**: 基于用户当前工作内容调整

#### 4.3.4 性能优化
- **搜索结果缓存**: 相同查询30分钟内返回缓存
- **预取相关数据**: 提前加载邮件完整内容
- **分页优化**: 实现高效的游标分页
- **搜索建议**: 基于历史记录提供搜索建议

---

## 五、监控与运维

### 5.1 关键性能指标 (KPIs)

| 指标类型 | 指标名称 | 正常范围 | 告警阈值 |
|----------|----------|----------|----------|
| **性能** | 搜索响应时间 | < 1秒 | > 2秒警告, > 5秒严重 |
| **性能** | 嵌入生成延迟 | < 500ms | > 1秒警告, > 2秒严重 |
| **质量** | 搜索准确率 | > 85% | < 80%警告 |
| **可靠性** | API调用成功率 | > 95% | < 90%警告 |
| **可靠性** | 系统可用性 | > 99.5% | < 99%警告 |

### 5.2 错误处理流程

```mermaid
flowchart TD
    A[错误发生] --> B{错误类型}
    B -->|维度不匹配| C[记录参数错误日志]
    B -->|API调用失败| D[触发重试机制]
    B -->|数据库错误| E[数据库连接检查]
    B -->|内存不足| F[触发资源清理]

    C --> G[返回400错误]
    D --> H{重试次数}
    H -->|< 3次| I[指数退避重试]
    H -->|≥ 3次| J[记录错误日志]

    E --> K[切换备用数据库]
    F --> L[重启相关服务]

    I --> M{重试结果}
    M -->|成功| N[恢复服务]
    M -->|失败| J

    J --> O[发送告警通知]
    K --> N
    L --> N
    G --> P[返回客户端错误]
    O --> P
```

### 5.3 运维最佳实践

#### 5.3.1 定期维护任务
- **向量索引重建**: 每周一次索引优化
- **缓存清理**: 每日清理过期缓存
- **日志归档**: 每月归档和压缩日志
- **数据备份**: 每日增量备份，每周全量备份

#### 5.3.2 监控告警设置
- **系统资源**: CPU、内存、磁盘使用率
- **数据库性能**: 查询响应时间、连接数
- **AI服务**: API调用成功率、响应延迟
- **业务指标**: 搜索量、同步成功率、用户活跃度

通过这些增强，EchoMind将从一个基础的邮件搜索系统升级为智能化的信息发现和处理平台。