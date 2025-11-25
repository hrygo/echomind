# Dashboard 组件实现规划

## 📋 项目概述

**项目**: EchoMind Dashboard 组件完善
**版本**: v0.9.6
**状态**: 规划阶段
**目标**: 将 Dashboard 从"中等成熟度"提升到"企业级完整实现"

## 🎯 当前状态分析

### ✅ 完整实现的组件 (3个)

| 组件 | 文件路径 | 功能描述 | API集成状态 |
|------|----------|----------|-------------|
| AIBriefingHeader | `components/dashboard/AIBriefingHeader.tsx` | AI问候头部 | ❌ 硬编码用户数据 |
| TaskWidget | `components/dashboard/TaskWidget.tsx` | 任务管理组件 | ✅ 完整API集成 |
| SmartFeed | `components/dashboard/SmartFeed.tsx` | 智能信息流 | ✅ 完整API集成 |

### ⚠️ 部分实现的组件 (3个)

| 组件 | 文件路径 | 问题描述 | 依赖API |
|------|----------|----------|---------|
| ManagerView | `components/dashboard/ManagerView.tsx` | 统计数据硬编码 | `/insights/manager/stats` |
| ExecutiveView | `components/dashboard/ExecutiveView.tsx` | 完全硬编码数据 | `/insights/executive/*` |
| DealmakerView | `components/dashboard/DealmakerView.tsx` | 商机数据假数据 | `/opportunities/*` |

## 🚀 分阶段实施计划

### 阶段一：完善现有组件API集成 (1-2周)

#### 1.1 ManagerView 统计集成
**优先级**: 🔥 高
**文件**: `components/dashboard/ManagerView.tsx`

**需要实现的API**:
```typescript
// GET /api/v1/insights/manager/stats
{
  "activeTasksCount": 5,
  "overdueTasksCount": 1,
  "completedTasksToday": 2,
  "teamProductivity": 85,
  "urgentEmailsCount": 3
}
```

**实现步骤**:
1. 后端实现统计API端点
2. 前端添加React Query数据获取
3. 替换硬编码数据
4. 添加加载状态和错误处理

#### 1.2 AIBriefingHeader 用户信息集成
**优先级**: 🔥 高
**文件**: `components/dashboard/AIBriefingHeader.tsx`

**需要实现的API**:
```typescript
// GET /api/v1/users/me/profile
{
  "id": "uuid",
  "name": "张三",
  "email": "zhang@example.com",
  "role": "manager",
  "avatar_url": "https://...",
  "preferences": {
    "language": "zh",
    "timezone": "Asia/Shanghai"
  }
}
```

### 阶段二：实现商机管理功能 (2-3周)

#### 2.1 DealmakerView 完整集成
**优先级**: 🔥 高
**文件**: `components/dashboard/DealmakerView.tsx`

**需要实现的API端点**:
```typescript
// GET /api/v1/opportunities
interface Opportunity {
  id: string;
  title: string;
  company: string;
  value: string;
  confidence: number;
  type: "buying" | "partnership" | "renewal";
  createdAt: string;
  contacts: Contact[];
}

// GET /api/v1/insights/dealmaker/radar
interface RadarData {
  category: string;
  value: number;
  fullMark: 150;
}[]
```

#### 2.2 商机状态管理
**新建文件**: `src/store/opportunityStore.ts`

```typescript
interface OpportunityStore {
  opportunities: Opportunity[];
  isLoading: boolean;
  error: string | null;

  fetchOpportunities: () => Promise<void>;
  updateOpportunity: (id: string, updates: Partial<Opportunity>) => Promise<void>;
  deleteOpportunity: (id: string) => Promise<void>;
}
```

### 阶段三：实现高管洞察功能 (2-3周)

#### 3.1 ExecutiveView 统计集成
**优先级**: 🔥 中
**文件**: `components/dashboard/ExecutiveView.tsx`

**需要实现的API端点**:
```typescript
// GET /api/v1/insights/executive/overview
{
  "totalConnections": 1247,
  "activeProjects": 8,
  "teamCollaborationScore": 92,
  "productivityTrend": "upward",
  "criticalAlerts": 2,
  "upcomingDeadlines": 5
}

// GET /api/v1/insights/risks
{
  "highRiskItems": [
    { id: 1, title: "关键客户续约", severity: "high", deadline: "2024-01-15" }
  ],
  "mediumRiskItems": [...],
  "riskTrend": "decreasing"
}

// GET /api/v1/insights/trends
{
  "productivity": [
    { date: "2024-01-01", value: 85 },
    { date: "2024-01-02", value: 88 }
  ],
  "collaboration": [...],
  "communication": [...]
}
```

### 阶段四：AI功能增强 (1-2周)

#### 4.1 SmartFeed AI回复
**优先级**: 🔥 中
**文件**: `components/dashboard/SmartFeed.tsx`

**需要实现的API**:
```typescript
// POST /api/v1/ai/reply
{
  "emailId": "uuid",
  "tone": "professional", // optional
  "context": "brief"      // optional
}

// Response
{
  "reply": "Dear [Name],\n\nThank you for your email...",
  "confidence": 0.95
}
```

### 阶段五：高级功能完善 (3-4周)

#### 5.1 实时数据更新
**优先级**: 🔥 低
- WebSocket集成实现实时数据同步
- 任务状态变更实时推送
- 新邮件实时通知

#### 5.2 高级筛选和搜索
**优先级**: 🔥 低
- 统一的筛选组件
- 多维度搜索功能
- 搜索历史保存

#### 5.3 数据导出功能
**优先级**: 🔥 低
- CSV、Excel、PDF格式导出
- 批量操作功能
- 自定义报告生成

## 📊 技术实施细节

### 后端API实现要求

#### 1. 数据模型定义
```go
// models/opportunity.go
type Opportunity struct {
    ID          string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    Title       string    `gorm:"not null"`
    Company     string    `gorm:"not null"`
    Value       string
    Confidence  int       `gorm:"check:confidence >= 0 AND confidence <= 100"`
    Type        string    `gorm:"type:opportunity_type;default:'buying'"`
    CreatedAt   time.Time `gorm:"autoCreateTime"`
    UpdatedAt   time.Time `gorm:"autoUpdateTime"`
    UserID      string    `gorm:"type:uuid"`
    TeamID      string    `gorm:"type:uuid"`
    OrgID       string    `gorm:"type:uuid"`
}

// models/insight.go
type ManagerStats struct {
    ActiveTasksCount    int `json:"activeTasksCount"`
    OverdueTasksCount   int `json:"overdueTasksCount"`
    CompletedTodayCount int `json:"completedTasksToday"`
    TeamProductivity    int `json:"teamProductivity"`
    UrgentEmailsCount   int `json:"urgentEmailsCount"`
}
```

#### 2. API端点路由
```go
// router/routes.go
insights := api.Group("/insights")
{
    insights.GET("/manager/stats", handler.GetManagerStats)
    insights.GET("/executive/overview", handler.GetExecutiveOverview)
    insights.GET("/executive/risks", handler.GetRisks)
    insights.GET("/executive/trends", handler.GetTrends)
    insights.GET("/dealmaker/radar", handler.GetDealmakerRadar)
}

opportunities := api.Group("/opportunities")
{
    opportunities.GET("", handler.ListOpportunities)
    opportunities.POST("", handler.CreateOpportunity)
    opportunities.GET("/:id", handler.GetOpportunity)
    opportunities.PUT("/:id", handler.UpdateOpportunity)
    opportunities.DELETE("/:id", handler.DeleteOpportunity)
}
```

### 前端实现要求

#### 1. 类型定义
```typescript
// types/api.ts
export interface Opportunity {
  id: string;
  title: string;
  company: string;
  value: string;
  confidence: number;
  type: OpportunityType;
  createdAt: string;
  contacts: Contact[];
  userId: string;
  teamId: string;
  orgId: string;
}

export interface ManagerStats {
  activeTasksCount: number;
  overdueTasksCount: number;
  completedTodayCount: number;
  teamProductivity: number;
  urgentEmailsCount: number;
}

export interface ExecutiveOverview {
  totalConnections: number;
  activeProjects: number;
  teamCollaborationScore: number;
  productivityTrend: 'upward' | 'downward' | 'stable';
  criticalAlerts: number;
  upcomingDeadlines: number;
}
```

#### 2. React Query Hooks
```typescript
// hooks/useOpportunities.ts
export const useOpportunities = () => {
  return useQuery({
    queryKey: ['opportunities'],
    queryFn: () => api.get('/opportunities').then(res => res.data),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};

export const useManagerStats = () => {
  return useQuery({
    queryKey: ['manager-stats'],
    queryFn: () => api.get('/insights/manager/stats').then(res => res.data),
    staleTime: 2 * 60 * 1000, // 2 minutes
  });
};
```

## 🧪 测试策略

### 单元测试要求
- 所有新增组件必须包含单元测试
- API端点必须有完整的测试覆盖
- Store状态管理需要测试验证

### 集成测试要求
- 前后端API集成测试
- 端到端用户流程测试
- 性能测试和负载测试

### 测试工具配置
```bash
# 后端测试
make test

# 前端测试
pnpm test
pnpm test:e2e
```

## 📈 性能要求

### 前端性能
- 页面首次加载时间 < 2秒
- API响应时间 < 500ms
- 大数据量列表需要虚拟滚动

### 后端性能
- API响应时间 < 200ms (不含数据库查询)
- 支持并发用户数 > 100
- 数据库查询优化和索引

## 🔍 监控和日志

### 前端监控
- 错误边界和异常捕获
- 用户行为分析
- 性能指标监控

### 后端监控
- API响应时间监控
- 错误率和异常监控
- 数据库性能监控

## 📚 文档要求

### API文档
- 所有新增API必须在OpenAPI文档中定义
- 包含完整的请求/响应示例
- 错误码和处理说明

### 代码文档
- 复杂业务逻辑必须有注释
- 组件props和state需要类型注解
- 关键算法需要文档说明

## 🚦 部署和发布

### 开发环境
- 本地开发环境配置
- 热重载和调试支持
- 模拟数据和Mock服务

### 测试环境
- 自动化部署到测试环境
- 集成测试和回归测试
- 性能测试和压力测试

### 生产环境
- 蓝绿部署策略
- 回滚机制和监控
- 数据备份和恢复

## 📋 检查清单

### 开发完成标准
- [ ] 所有组件功能完整实现
- [ ] API集成完成且测试通过
- [ ] UI/UX设计符合规范
- [ ] 国际化支持完整
- [ ] 性能指标达标
- [ ] 安全性检查通过
- [ ] 文档完整且最新

### 质量保证标准
- [ ] 代码审查通过
- [ ] 单元测试覆盖率 > 80%
- [ ] 集成测试通过
- [ ] 端到端测试通过
- [ ] 性能测试通过
- [ ] 安全测试通过

## 🔄 迭代计划

### 当前Sprint目标
- 完成阶段一：现有组件API集成
- 实现ManagerView和AIBriefingHeader的完整功能

### 下一个Sprint计划
- 开始阶段二：商机管理功能
- 实现DealmakerView的完整API集成

### 长期目标
- 完成所有5个阶段的开发
- 达到企业级Dashboard功能
- 支持多角色和复杂业务场景

---

*此文档将根据开发进展持续更新，确保与实际实施保持同步。*