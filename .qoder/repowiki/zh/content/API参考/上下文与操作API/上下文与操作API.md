# 上下文与操作API

<cite>
**本文档中引用的文件**
- [backend/internal/handler/context.go](file://backend/internal/handler/context.go)
- [backend/internal/model/context.go](file://backend/internal/model/context.go)
- [backend/internal/model/context_input.go](file://backend/internal/model/context_input.go)
- [backend/internal/handler/action.go](file://backend/internal/handler/action.go)
- [backend/internal/service/action.go](file://backend/internal/service/action.go)
- [backend/internal/service/context.go](file://backend/internal/service/context.go)
- [backend/internal/router/routes.go](file://backend/internal/router/routes.go)
- [backend/pkg/event/bus/bus.go](file://backend/pkg/event/bus/bus.go)
- [frontend/src/lib/api/contexts.ts](file://frontend/src/lib/api/contexts.ts)
- [frontend/src/lib/api/actions.ts](file://frontend/src/lib/api/actions.ts)
- [backend/internal/listener/email_listeners.go](file://backend/internal/listener/email_listeners.go)
- [backend/internal/tasks/analyze.go](file://backend/internal/tasks/analyze.go)
- [backend/cmd/backfill_contexts/main.go](file://backend/cmd/backfill_contexts/main.go)
</cite>

## 目录
1. [简介](#简介)
2. [项目结构概览](#项目结构概览)
3. [上下文API详解](#上下文api详解)
4. [操作API详解](#操作api详解)
5. [系统架构分析](#系统架构分析)
6. [事件总线集成](#事件总线集成)
7. [数据模型设计](#数据模型设计)
8. [最佳实践指南](#最佳实践指南)
9. [故障排除](#故障排除)
10. [总结](#总结)

## 简介

EchoMind是一个智能邮件管理系统，提供了强大的上下文管理和用户操作API。该系统允许用户创建自定义的邮件分组（上下文），并通过智能算法自动将邮件分配到相应的上下文中。同时，系统提供了三种核心操作：批准（approve）、稍后处理（snooze）和忽略（dismiss），这些操作反映了用户对智能建议的反馈，并触发后续处理流程。

### 核心概念

**上下文（Context）**：用户自定义的邮件分组，用于组织和管理邮件。每个上下文包含名称、颜色标识、关键词规则和利益相关者列表。

**操作（Action）**：用户对智能建议的反馈机制，包括：
- **批准（Approve）**：标记邮件为已完成，通常会归档邮件
- **稍后处理（Snooze）**：延迟处理邮件，设置特定时间恢复
- **忽略（Dismiss）**：从智能推荐中移除邮件

## 项目结构概览

系统采用分层架构设计，主要分为以下几个层次：

```mermaid
graph TB
subgraph "前端层"
FE[前端应用]
API_CLIENT[API客户端]
end
subgraph "HTTP层"
ROUTER[路由器]
HANDLER[处理器]
end
subgraph "业务逻辑层"
SERVICE[服务层]
CONTEXT_SERVICE[上下文服务]
ACTION_SERVICE[操作服务]
end
subgraph "数据访问层"
REPOSITORY[仓库层]
MODEL[数据模型]
end
subgraph "基础设施层"
EVENT_BUS[事件总线]
TASK_QUEUE[任务队列]
end
FE --> API_CLIENT
API_CLIENT --> ROUTER
ROUTER --> HANDLER
HANDLER --> SERVICE
SERVICE --> REPOSITORY
REPOSITORY --> MODEL
SERVICE --> EVENT_BUS
SERVICE --> TASK_QUEUE
```

**图表来源**
- [backend/internal/router/routes.go](file://backend/internal/router/routes.go#L26-L98)
- [backend/internal/handler/context.go](file://backend/internal/handler/context.go#L12-L18)
- [backend/internal/handler/action.go](file://backend/internal/handler/action.go#L12-L18)

## 上下文API详解

### 创建上下文组 - POST /api/v1/contexts

创建新的邮件分组上下文，支持关键词匹配和利益相关者识别。

#### 请求格式

```typescript
interface ContextInput {
  name: string;           // 上下文名称（必填，最大100字符）
  color: string;          // 颜色标识（可选，最大20字符）
  keywords: string[];     // 关键词列表
  stakeholders: string[]; // 利益相关者邮箱列表
}
```

#### 示例请求

```json
{
  "name": "重要项目",
  "color": "blue",
  "keywords": ["Project Alpha", "Q4 Budget", "财务报告"],
  "stakeholders": ["boss@example.com", "team@company.com"]
}
```

#### 响应格式

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "重要项目",
  "color": "blue",
  "keywords": ["Project Alpha", "Q4 Budget", "财务报告"],
  "stakeholders": ["boss@example.com", "team@company.com"],
  "createdAt": "2024-01-15T10:30:00Z",
  "updatedAt": "2024-01-15T10:30:00Z"
}
```

#### 处理流程

```mermaid
sequenceDiagram
participant Client as 客户端
participant Handler as 上下文处理器
participant Service as 上下文服务
participant DB as 数据库
participant EventBus as 事件总线
Client->>Handler : POST /api/v1/contexts
Handler->>Handler : 验证用户身份
Handler->>Handler : 解析请求体
Handler->>Service : CreateContext(userID, input)
Service->>Service : 序列化关键词和利益相关者
Service->>DB : 创建上下文记录
DB-->>Service : 返回新上下文
Service-->>Handler : 返回上下文对象
Handler->>EventBus : 发布上下文创建事件
EventBus-->>Handler : 确认发布
Handler-->>Client : 返回201 Created
```

**图表来源**
- [backend/internal/handler/context.go](file://backend/internal/handler/context.go#L21-L36)
- [backend/internal/service/context.go](file://backend/internal/service/context.go#L22-L46)

**节来源**
- [backend/internal/handler/context.go](file://backend/internal/handler/context.go#L21-L36)
- [backend/internal/model/context_input.go](file://backend/internal/model/context_input.go#L4-L9)

### 列出上下文 - GET /api/v1/contexts

获取当前用户的全部上下文列表，按创建时间降序排列。

#### 响应格式

```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "重要项目",
    "color": "blue",
    "keywords": ["Project Alpha", "Q4 Budget"],
    "stakeholders": ["boss@example.com"],
    "createdAt": "2024-01-15T10:30:00Z",
    "updatedAt": "2024-01-15T10:30:00Z"
  },
  {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "name": "日常事务",
    "color": "green",
    "keywords": ["会议", "日程安排"],
    "stakeholders": [],
    "createdAt": "2024-01-14T09:15:00Z",
    "updatedAt": "2024-01-14T09:15:00Z"
  }
]
```

#### 处理流程

```mermaid
flowchart TD
Start([开始请求]) --> ValidateAuth["验证用户身份"]
ValidateAuth --> FetchContexts["从数据库获取上下文"]
FetchContexts --> SortByDate["按创建时间排序"]
SortByDate --> SerializeResponse["序列化响应数据"]
SerializeResponse --> SendResponse["返回200 OK"]
SendResponse --> End([结束])
```

**图表来源**
- [backend/internal/handler/context.go](file://backend/internal/handler/context.go#L39-L49)
- [backend/internal/service/context.go](file://backend/internal/service/context.go#L49-L55)

**节来源**
- [backend/internal/handler/context.go](file://backend/internal/handler/context.go#L39-L49)
- [backend/internal/service/context.go](file://backend/internal/service/context.go#L49-L55)

### 更新上下文 - PATCH /api/v1/contexts/:id

修改现有上下文的属性，支持部分更新。

#### 请求参数

- **路径参数**：`:id` - 上下文唯一标识符
- **请求体**：同创建请求的结构

#### 处理流程

```mermaid
flowchart TD
Start([开始更新]) --> ValidateID["验证上下文ID格式"]
ValidateID --> CheckOwnership["检查用户所有权"]
CheckOwnership --> ParseInput["解析更新输入"]
ParseInput --> ValidateInput["验证输入数据"]
ValidateInput --> UpdateContext["更新数据库记录"]
UpdateContext --> SerializeOutput["序列化输出"]
SerializeOutput --> SendResponse["返回200 OK"]
SendResponse --> End([结束])
```

**图表来源**
- [backend/internal/handler/context.go](file://backend/internal/handler/context.go#L52-L74)
- [backend/internal/service/context.go](file://backend/internal/service/context.go#L67-L92)

**节来源**
- [backend/internal/handler/context.go](file://backend/internal/handler/context.go#L52-L74)
- [backend/internal/service/context.go](file://backend/internal/service/context.go#L67-L92)

### 删除上下文 - DELETE /api/v1/contexts/:id

删除指定的上下文及其关联关系。

#### 处理流程

```mermaid
flowchart TD
Start([开始删除]) --> ValidateID["验证上下文ID"]
ValidateID --> CheckOwnership["检查用户所有权"]
CheckOwnership --> DeleteContext["删除上下文记录"]
DeleteContext --> CleanupRelations["清理关联关系"]
CleanupRelations --> SendResponse["返回204 No Content"]
SendResponse --> End([结束])
```

**图表来源**
- [backend/internal/handler/context.go](file://backend/internal/handler/context.go#L77-L92)
- [backend/internal/service/context.go](file://backend/internal/service/context.go#L95-L104)

**节来源**
- [backend/internal/handler/context.go](file://backend/internal/handler/context.go#L77-L92)
- [backend/internal/service/context.go](file://backend/internal/service/context.go#L95-L104)

## 操作API详解

### 批准邮件 - POST /api/v1/actions/approve

标记邮件为已处理完成，通常会归档邮件。

#### 请求格式

```typescript
interface ApproveRequest {
  email_id: string; // 邮件唯一标识符（必填）
}
```

#### 示例请求

```json
{
  "email_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

#### 响应格式

```json
{
  "status": "approved"
}
```

#### 处理流程

```mermaid
sequenceDiagram
participant Client as 客户端
participant Handler as 操作处理器
participant Service as 操作服务
participant DB as 数据库
participant EventBus as 事件总线
Client->>Handler : POST /api/v1/actions/approve
Handler->>Handler : 验证用户身份
Handler->>Handler : 解析邮件ID
Handler->>Service : ApproveEmail(userID, emailID)
Service->>DB : 验证邮件所有权
DB-->>Service : 验证结果
Service->>DB : 软删除邮件记录
DB-->>Service : 操作结果
Service->>EventBus : 发布邮件处理事件
EventBus-->>Service : 确认发布
Service-->>Handler : 返回成功
Handler-->>Client : 返回200 OK
```

**图表来源**
- [backend/internal/handler/action.go](file://backend/internal/handler/action.go#L33-L53)
- [backend/internal/service/action.go](file://backend/internal/service/action.go#L20-L48)

**节来源**
- [backend/internal/handler/action.go](file://backend/internal/handler/action.go#L33-L53)
- [backend/internal/service/action.go](file://backend/internal/service/action.go#L20-L48)

### 稍后处理 - POST /api/v1/actions/snooze

延迟处理邮件，设置特定时间恢复显示。

#### 请求格式

```typescript
interface SnoozeRequest {
  email_id: string;     // 邮件唯一标识符（必填）
  duration?: string;    // 延迟时长："4h", "tomorrow", "next_week" 或 ISO 时间戳
}
```

#### 支持的持续时间格式

| 格式 | 描述 | 示例 |
|------|------|------|
| `"4h"` | 4小时后 | 默认延迟时间 |
| `"tomorrow"` | 明天早上9点 | 自动调整到工作日 |
| `"next_week"` | 下周一早上9点 | 自动调整到工作日 |
| `"ISO时间戳"` | 具体时间 | `2024-01-15T09:00:00Z` |

#### 示例请求

```json
{
  "email_id": "550e8400-e29b-41d4-a716-446655440000",
  "duration": "tomorrow"
}
```

#### 响应格式

```json
{
  "status": "snoozed",
  "until": "2024-01-16T09:00:00Z"
}
```

#### 处理流程

```mermaid
flowchart TD
Start([开始稍后处理]) --> ParseDuration["解析延迟时长"]
ParseDuration --> ValidateEmail["验证邮件ID"]
ValidateEmail --> DetermineTime["确定具体时间"]
DetermineTime --> UpdateDB["更新数据库"]
UpdateDB --> ScheduleTask["调度后续任务"]
ScheduleTask --> SendResponse["返回响应"]
SendResponse --> End([结束])
```

**图表来源**
- [backend/internal/handler/action.go](file://backend/internal/handler/action.go#L56-L100)
- [backend/internal/service/action.go](file://backend/internal/service/action.go#L51-L63)

**节来源**
- [backend/internal/handler/action.go](file://backend/internal/handler/action.go#L56-L100)
- [backend/internal/service/action.go](file://backend/internal/service/action.go#L51-L63)

### 忽略邮件 - POST /api/v1/actions/dismiss

从智能推荐中移除邮件，降低其紧急程度。

#### 请求格式

```typescript
interface DismissRequest {
  email_id: string; // 邮件唯一标识符（必填）
}
```

#### 示例请求

```json
{
  "email_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

#### 响应格式

```json
{
  "status": "dismissed"
}
```

#### 处理流程

```mermaid
flowchart TD
Start([开始忽略]) --> ValidateEmail["验证邮件ID"]
ValidateEmail --> UpdateUrgency["更新紧急程度为低"]
UpdateUrgency --> RemoveFromFeed["从智能推荐中移除"]
RemoveFromFeed --> LogAction["记录操作日志"]
LogAction --> SendResponse["返回200 OK"]
SendResponse --> End([结束])
```

**图表来源**
- [backend/internal/handler/action.go](file://backend/internal/handler/action.go#L103-L123)
- [backend/internal/service/action.go](file://backend/internal/service/action.go#L66-L78)

**节来源**
- [backend/internal/handler/action.go](file://backend/internal/handler/action.go#L103-L123)
- [backend/internal/service/action.go](file://backend/internal/service/action.go#L66-L78)

## 系统架构分析

### 整体架构图

```mermaid
graph TB
subgraph "前端应用层"
WEB_UI[Web界面]
MOBILE_APP[移动应用]
end
subgraph "API网关层"
AUTH_MIDDLEWARE[认证中间件]
ROUTER[路由器]
end
subgraph "控制器层"
CONTEXT_HANDLER[上下文处理器]
ACTION_HANDLER[操作处理器]
end
subgraph "服务层"
CONTEXT_SERVICE[上下文服务]
ACTION_SERVICE[操作服务]
EMAIL_SERVICE[邮件服务]
end
subgraph "任务处理层"
ANALYZE_TASK[邮件分析任务]
CONTEXT_MATCHER[上下文匹配器]
end
subgraph "事件总线层"
EVENT_BUS[事件总线]
LISTENERS[事件监听器]
end
subgraph "数据存储层"
POSTGRES_DB[(PostgreSQL数据库)]
REDIS_CACHE[(Redis缓存)]
end
WEB_UI --> AUTH_MIDDLEWARE
MOBILE_APP --> AUTH_MIDDLEWARE
AUTH_MIDDLEWARE --> ROUTER
ROUTER --> CONTEXT_HANDLER
ROUTER --> ACTION_HANDLER
CONTEXT_HANDLER --> CONTEXT_SERVICE
ACTION_HANDLER --> ACTION_SERVICE
CONTEXT_SERVICE --> POSTGRES_DB
ACTION_SERVICE --> POSTGRES_DB
EMAIL_SERVICE --> POSTGRES_DB
ANALYZE_TASK --> CONTEXT_MATCHER
CONTEXT_MATCHER --> POSTGRES_DB
EVENT_BUS --> LISTENERS
LISTENERS --> POSTGRES_DB
```

**图表来源**
- [backend/internal/router/routes.go](file://backend/internal/router/routes.go#L26-L98)
- [backend/internal/handler/context.go](file://backend/internal/handler/context.go#L12-L18)
- [backend/internal/handler/action.go](file://backend/internal/handler/action.go#L12-L18)

### 上下文匹配算法

系统使用智能算法自动将邮件分配到合适的上下文中：

```mermaid
flowchart TD
Start([邮件同步完成]) --> FetchContexts["获取用户所有上下文"]
FetchContexts --> IterateContexts["遍历每个上下文"]
IterateContexts --> CheckKeywords["检查关键词匹配"]
CheckKeywords --> KeywordsMatch{"关键词匹配？"}
KeywordsMatch --> |是| AddMatch["添加到匹配列表"]
KeywordsMatch --> |否| CheckStakeholders["检查利益相关者"]
CheckStakeholders --> StakeholdersMatch{"利益相关者匹配？"}
StakeholdersMatch --> |是| AddMatch
StakeholdersMatch --> |否| NextContext["下一个上下文"]
AddMatch --> NextContext
NextContext --> MoreContexts{"还有更多上下文？"}
MoreContexts --> |是| IterateContexts
MoreContexts --> |否| AssignContexts["批量分配上下文"]
AssignContexts --> End([完成])
```

**图表来源**
- [backend/internal/service/context.go](file://backend/internal/service/context.go#L107-L151)
- [backend/internal/tasks/analyze.go](file://backend/internal/tasks/analyze.go#L153-L166)

**节来源**
- [backend/internal/service/context.go](file://backend/internal/service/context.go#L107-L151)
- [backend/internal/tasks/analyze.go](file://backend/internal/tasks/analyze.go#L153-L166)

## 事件总线集成

### 事件总线架构

系统使用事件驱动架构，通过事件总线实现组件间的解耦：

```mermaid
classDiagram
class Event {
<<interface>>
+Name() string
}
class Listener {
<<interface>>
+Handle(ctx Context, event Event) error
}
class Bus {
-listeners map[string][]Listener
-mu sync.RWMutex
+Subscribe(eventName string, listener Listener)
+Publish(ctx Context, event Event) error
}
class EmailSyncedEvent {
+UserID uuid.UUID
+Email Email
+Name() string
}
class AnalysisListener {
-asynqClient AsynqClientInterface
-logger CompatibleLogger
+Handle(ctx Context, event Event) error
}
class ContactListener {
-contactService *ContactService
-logger CompatibleLogger
+Handle(ctx Context, event Event) error
}
Event <|-- EmailSyncedEvent
Listener <|-- AnalysisListener
Listener <|-- ContactListener
Bus --> Listener
Bus --> Event
```

**图表来源**
- [backend/pkg/event/bus/bus.go](file://backend/pkg/event/bus/bus.go#L8-L62)
- [backend/internal/listener/email_listeners.go](file://backend/internal/listener/email_listeners.go#L22-L115)

### 事件处理流程

```mermaid
sequenceDiagram
participant Sync as 邮件同步
participant EventBus as 事件总线
participant Analysis as 分析监听器
participant Contact as 联系人监听器
participant TaskQueue as 任务队列
Sync->>EventBus : 发布 EmailSyncedEvent
EventBus->>Analysis : 调用分析监听器
Analysis->>TaskQueue : 创建邮件分析任务
EventBus->>Contact : 调用联系人监听器
Contact->>Contact : 更新联系人统计
TaskQueue-->>Analysis : 任务创建确认
Contact-->>EventBus : 处理完成
Analysis-->>EventBus : 处理完成
EventBus-->>Sync : 所有监听器处理完成
```

**图表来源**
- [backend/internal/listener/email_listeners.go](file://backend/internal/listener/email_listeners.go#L35-L65)
- [backend/internal/listener/email_listeners.go](file://backend/internal/listener/email_listeners.go#L80-L101)

**节来源**
- [backend/pkg/event/bus/bus.go](file://backend/pkg/event/bus/bus.go#L8-L62)
- [backend/internal/listener/email_listeners.go](file://backend/internal/listener/email_listeners.go#L35-L101)

## 数据模型设计

### 上下文数据模型

```mermaid
erDiagram
CONTEXT {
uuid id PK
uuid user_id FK
string name
string color
jsonb keywords
jsonb stakeholders
datetime created_at
datetime updated_at
datetime deleted_at
}
EMAIL_CONTEXT {
uuid email_id PK,FK
uuid context_id PK,FK
}
EMAIL {
uuid id PK
uuid user_id FK
string subject
string snippet
string sender
datetime date
datetime created_at
datetime updated_at
datetime deleted_at
}
CONTEXT ||--o{ EMAIL_CONTEXT : "belongs to"
EMAIL ||--o{ EMAIL_CONTEXT : "assigned to"
```

**图表来源**
- [backend/internal/model/context.go](file://backend/internal/model/context.go#L11-L30)

### 操作状态转换

```mermaid
stateDiagram-v2
[*] --> Inbox : 邮件到达
Inbox --> Approved : 批准操作
Inbox --> Snoozed : 稍后处理
Inbox --> Dismissed : 忽略操作
Approved --> Archived : 自动归档
Snoozed --> Inbox : 到达预定时间
Dismissed --> LowUrgency : 降低紧急度
Archived --> [*]
LowUrgency --> [*]
```

**节来源**
- [backend/internal/model/context.go](file://backend/internal/model/context.go#L11-L30)
- [backend/internal/service/action.go](file://backend/internal/service/action.go#L20-L78)

## 最佳实践指南

### 上下文创建最佳实践

1. **命名规范**：使用清晰、描述性的名称，避免过于宽泛或模糊的术语
2. **关键词策略**：
   - 包含具体的项目名称、产品名称
   - 使用行业术语和专业词汇
   - 避免过于通用的词语
3. **利益相关者管理**：
   - 包含直接负责人的邮箱地址
   - 考虑团队协作成员
   - 定期更新联系人列表

### 操作使用指南

1. **批准操作**：
   - 仅用于真正完成的任务
   - 确保邮件内容已妥善处理
   - 可能触发自动化工作流

2. **稍后处理**：
   - 设置合理的延迟时间
   - 考虑工作日程安排
   - 避免过长的延迟时间

3. **忽略操作**：
   - 用于不重要的邮件
   - 不影响长期邮件管理
   - 可重新激活被忽略的邮件

### 性能优化建议

1. **批量操作**：对于大量邮件的操作，考虑使用批量API
2. **缓存策略**：对频繁访问的上下文信息进行缓存
3. **异步处理**：大型操作使用后台任务处理

## 故障排除

### 常见问题及解决方案

#### 上下文创建失败

**问题症状**：
- HTTP 400 错误
- 输入验证失败

**可能原因**：
- 上下文名称为空或过长
- 关键词列表包含无效字符
- 利益相关者邮箱格式错误

**解决方案**：
- 检查输入数据格式
- 验证字段长度限制
- 使用正确的邮箱格式

#### 邮件分配不准确

**问题症状**：
- 邮件未分配到正确上下文
- 分配过多或过少的邮件

**可能原因**：
- 关键词匹配规则过于宽松或严格
- 利益相关者列表不完整
- 上下文数量过多导致冲突

**解决方案**：
- 调整关键词精确度
- 完善利益相关者列表
- 合理规划上下文结构

#### 操作执行超时

**问题症状**：
- 操作请求长时间无响应
- 服务器超时错误

**可能原因**：
- 数据库连接问题
- 事件总线阻塞
- 后台任务积压

**解决方案**：
- 检查数据库连接状态
- 监控事件总线性能
- 清理积压任务

**节来源**
- [backend/internal/handler/context.go](file://backend/internal/handler/context.go#L24-L28)
- [backend/internal/service/context.go](file://backend/internal/service/context.go#L107-L151)

## 总结

EchoMind的上下文与操作API提供了一个强大而灵活的邮件管理系统框架。通过上下文功能，用户可以创建个性化的邮件分组，实现智能化的邮件分类和管理。操作API则提供了用户反馈机制，使系统能够根据用户的实际行为不断优化智能建议。

### 主要特性

1. **智能上下文匹配**：基于关键词和利益相关者的自动邮件分类
2. **灵活的操作选项**：三种基本操作满足不同场景需求
3. **事件驱动架构**：通过事件总线实现组件解耦和扩展性
4. **实时反馈机制**：用户操作立即反映到系统行为中
5. **可扩展的设计**：支持自定义上下文规则和操作类型

### 技术优势

- **模块化架构**：清晰的分层设计便于维护和扩展
- **异步处理**：大量数据处理使用后台任务队列
- **强一致性**：关键操作保证数据完整性
- **高性能**：合理使用缓存和索引优化查询性能

这套API设计不仅满足了当前的功能需求，还为未来的功能扩展和系统优化奠定了坚实的基础。通过持续的迭代和改进，EchoMind将成为一个更加智能和高效的邮件管理平台。