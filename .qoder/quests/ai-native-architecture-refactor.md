# AI Native 架构重构终版设计方案

## 项目背景

EchoMind 当前前端基于 Next.js 16.0.3 和 React 19.2.0 构建,现需要进行全面架构升级,转向 AI Native 架构模式。本次重构旨在充分利用 Next.js 16 和 React 19 的最新特性,对接后端 AI 服务,并全面采用现代化的 UI 组件体系和状态管理方案。

### 当前架构状态

- **框架版本**: Next.js 16.0.3, React 19.2.0 (已是最新版本)
- **UI 组件**: 使用 Radix UI + Tailwind CSS,部分采用 shadcn/ui 模式
- **状态管理**: Zustand v5 + TanStack Query v5
- **国际化**: 自定义 LanguageContext 实现
- **主题系统**: next-themes v0.4.6
- **构建工具**: 默认使用 Turbopack
- **AI 能力**: 后端已提供完整 AI API (草稿生成/聊天/搜索)

### 重构目标

1. **架构升级**: 全面迁移至 AI Native 架构,对接后端 AI 服务
2. **UI 标准化**: 完全采用 shadcn/ui 组件体系,替换自定义实现
3. **框架规范**: 实现 Next.js 16 异步 API 规范 (params/searchParams)
4. **React 19**: 充分利用新特性 (Server Actions, useActionState, use hook)
5. **国际化升级**: 采用 Next.js 原生支持的 Dictionary Pattern
6. **流式处理**: 实现 SSE 客户端对接后端流式 AI 响应
7. **性能优化**: 提升首屏加载速度和 AI 交互体验

## 关键技术决策

### 决策 1: AI 服务架构 - 采用后端 API 模式

**决策内容**: 前端不集成 Vercel AI SDK,直接调用后端已有的 AI API

**决策依据**:
1. 后端已实现完整 AI 能力 (草稿生成/聊天/搜索)
2. 后端支持 SSE 流式响应
3. AI 成本和配额由后端统一管理
4. Prompt 模板集中在后端,方便统一优化

**优势**:
- 前后端职责清晰分离
- 避免前端暴露 AI API 密钥
- 降低前端包体积
- 后端可统一实施限流和缓存策略

**权衡**:
- 前端需要实现 SSE 客户端 (增加少量代码)
- 依赖后端接口稳定性

### 决策 2: 国际化方案 - 保留现状

**决策内容**: 继续使用现有的 LanguageContext 模式,不进行架构重构

**决策依据**:
1. 现有国际化方案已稳定运行
2. 迁移到 Dictionary Pattern 成本较高
3. 对核心业务功能影响较小
4. 团队熟悉现有实现方式

**优势**:
- 零迁移成本
- 避免路由结构大幅调整
- 降低重构风险
- 团队无需学习新模式

**未来改进**:
- 可在后续版本中评估迁移到 Dictionary Pattern
- 当前专注于 AI Native 架构核心功能

### 决策 3: UI 组件库 - 全面采用 shadcn/ui

**决策内容**: 完全替换现有自定义 UI 组件为 shadcn/ui 标准实现

**决策依据**:
1. shadcn/ui 是当前最佳 UI 组件方案
2. 与 Radix UI 和 Tailwind CSS 完美集成
3. 社区活跃,组件丰富
4. 支持主题系统和无障碍访问

**优势**:
- 代码质量和可维护性更高
- 组件一致性更好
- 社区支持和文档完善
- 持续更新和优化

**权衡**:
- 需要调整现有组件的样式和行为
- 迁移工作量约 2 天

### 决策 4: 状态管理 - 保持 Zustand + TanStack Query

**决策内容**: 继续使用现有的 Zustand v5 和 TanStack Query v5

**决策依据**:
1. 现有方案已经很成熟
2. Zustand v5 与 React 19 完全兼容
3. TanStack Query v5 针对 React 19 优化
4. 不需要额外的学习成本

**优势**:
- 零迁移成本
- 团队熟悉现有方案
- 性能优秀

### 决策 5: Server Actions 场景 - 表单提交为主

**决策内容**: 仅在表单提交场景使用 Server Actions,非表单操作继续使用 API 调用

**决策依据**:
1. Server Actions 最适合表单处理
2. 与 useActionState 完美配合
3. 避免过度使用导致架构混乱

**应用场景**:
- 登录/注册表单
- 组织创建/更新表单
- 邮件同步表单

## 架构设计

### 核心技术栈

| 技术领域 | 选型方案 | 版本 | 备注 |
|---------|---------|------|------|
| 框架 | Next.js (App Router) | 16.0.3 | 已满足要求 |
| UI 库 | React | 19.2.0 | 已满足要求 |
| 样式方案 | Tailwind CSS v4 | ^4 | 已升级至 v4 |
| 组件库 | shadcn/ui | latest | 需全面采用 |
| 工具库 | clsx + tailwind-merge | latest | 已集成 |
| AI 服务 | 后端 AI API | - | 使用现有后端接口 |
| 流式处理 | EventSource (SSE) | - | 对接后端 SSE 流 |
| 服务端状态 | TanStack Query | v5.90.10 | 已满足要求 |
| 客户端状态 | Zustand | v5.0.8 | 已满足要求 |
| 国际化 | LanguageContext | - | 保持现状 |
| 构建工具 | Turbopack | - | Next.js 16 默认 |

### 目录结构设计

```
frontend/
├── src/
│   ├── app/                              # Next.js 16 App Router
│   │   ├── (auth)/                       # 认证路由组
│   │   │   ├── auth/
│   │   │   │   ├── page.tsx              # 登录/注册页面
│   │   │   │   └── layout.tsx
│   │   │   └── onboarding/
│   │   │       └── page.tsx
│   │   ├── (dashboard)/                  # 主应用路由组
│   │   │   ├── layout.tsx                # Dashboard 布局
│   │   │   ├── page.tsx                  # 首页 (重定向到 inbox)
│   │   │   ├── inbox/
│   │   │   │   └── page.tsx
│   │   │   ├── search/
│   │   │   │   └── page.tsx              # 搜索页面 (async searchParams)
│   │   │   ├── insights/
│   │   │   │   └── page.tsx
│   │   │   ├── tasks/
│   │   │   │   └── page.tsx
│   │   │   ├── opportunities/
│   │   │   │   ├── page.tsx
│   │   │   │   └── [id]/
│   │   │   │       └── page.tsx          # 动态路由 (async params)
│   │   │   └── settings/
│   │   │       └── page.tsx
│   │   ├── layout.tsx                    # 根布局
│   │   ├── globals.css                   # 全局样式
│   │   └── not-found.tsx
│   │
│   ├── components/                       # 组件系统
│   │   ├── ui/                           # shadcn/ui 原子组件
│   │   │   ├── button.tsx
│   │   │   ├── input.tsx
│   │   │   ├── card.tsx
│   │   │   ├── dialog.tsx
│   │   │   ├── sheet.tsx
│   │   │   ├── dropdown-menu.tsx
│   │   │   ├── tabs.tsx
│   │   │   ├── label.tsx
│   │   │   ├── checkbox.tsx
│   │   │   ├── toast.tsx
│   │   │   ├── skeleton.tsx
│   │   │   ├── avatar.tsx
│   │   │   ├── badge.tsx
│   │   │   └── ... (其他 shadcn 组件)
│   │   ├── layout/                       # 布局组件
│   │   │   ├── header.tsx
│   │   │   ├── sidebar.tsx
│   │   │   ├── mobile-nav.tsx
│   │   │   └── dashboard-shell.tsx
│   │   ├── features/                     # 功能性复合组件
│   │   │   ├── copilot/
│   │   │   │   ├── copilot-panel.tsx
│   │   │   │   ├── chat-interface.tsx
│   │   │   │   ├── search-results.tsx
│   │   │   │   └── message-stream.tsx    # AI 流式渲染
│   │   │   ├── email/
│   │   │   │   ├── email-list.tsx
│   │   │   │   └── email-viewer.tsx
│   │   │   ├── insights/
│   │   │   │   ├── insight-dashboard.tsx
│   │   │   │   └── chart-widgets.tsx
│   │   │   ├── opportunities/
│   │   │   │   ├── opportunity-list.tsx
│   │   │   │   └── opportunity-form.tsx
│   │   │   └── auth/
│   │   │       ├── login-form.tsx        # 使用 useActionState
│   │   │       └── register-form.tsx
│   │   ├── providers/                    # Context 提供者
│   │   │   ├── query-provider.tsx
│   │   │   ├── theme-provider.tsx
│   │   │   └── toast-provider.tsx
│   │   └── widgets/                      # 可复用小部件
│   │       ├── language-switcher.tsx
│   │       ├── theme-toggle.tsx
│   │       ├── organization-selector.tsx
│   │       └── ... (其他小部件)
│   │
│   ├── lib/                              # 工具库
│   │   ├── ai/                           # AI 客户端封装
│   │   │   ├── chat-client.ts            # SSE 流式聊天客户端
│   │   │   ├── draft-client.ts           # AI 草稿生成客户端
│   │   │   └── stream-parser.ts          # SSE 流解析器
│   │   ├── api/                          # API 客户端
│   │   │   ├── client.ts                 # Axios 实例配置
│   │   │   ├── endpoints/                # API 端点定义
│   │   │   │   ├── auth.ts
│   │   │   │   ├── email.ts
│   │   │   │   ├── search.ts
│   │   │   │   ├── insights.ts
│   │   │   │   ├── opportunities.ts
│   │   │   │   └── tasks.ts
│   │   │   └── types.ts                  # API 类型定义
│   │   ├── utils/                        # 工具函数
│   │   │   ├── cn.ts                     # clsx + tailwind-merge
│   │   │   ├── date.ts
│   │   │   ├── format.ts
│   │   │   └── validators.ts
│   │   ├── hooks/                        # 自定义 Hooks
│   │   │   ├── use-toast.ts
│   │   │   ├── use-media-query.ts
│   │   │   └── use-debounce.ts
│   │   └── constants/                    # 常量定义
│   │       ├── routes.ts
│   │       └── config.ts
│   │
│   ├── store/                            # Zustand 状态管理
│   │   ├── auth.ts                       # 认证状态
│   │   ├── ui.ts                         # UI 状态 (侧边栏、模态框等)
│   │   ├── copilot.ts                    # Copilot 状态
│   │   ├── organization.ts               # 组织状态
│   │   └── index.ts                      # Store 聚合导出
│   │
│   ├── types/                            # TypeScript 类型
│   │   ├── index.ts                      # 通用类型
│   │   ├── models.ts                     # 数据模型
│   │   └── api.ts                        # API 类型 (从 lib/api 导入)
│   │
│   ├── actions/                          # Server Actions (React 19)
│   │   ├── auth.ts                       # 登录/注册 actions
│   │   ├── email.ts                      # 邮件操作 actions
│   │   ├── organization.ts               # 组织管理 actions
│   │   └── opportunities.ts              # 机会管理 actions
│   │
│   └── middleware.ts                     # Next.js 中间件
│
├── public/                               # 静态资源
│   ├── locales/                          # (如果需要静态字典文件)
│   └── ... (其他静态资源)
│
├── .env.local                            # 环境变量
├── next.config.ts                        # Next.js 配置
├── tailwind.config.ts                    # Tailwind 配置
├── tsconfig.json                         # TypeScript 配置
├── components.json                       # shadcn/ui 配置
└── package.json
```

### 技术实施策略

#### 1. Next.js 16 异步 API 规范

Next.js 16 要求所有页面的 `params` 和 `searchParams` 必须作为 Promise 处理:

**页面组件实现规范**:
- 所有 `page.tsx` 必须声明为异步函数
- 使用 `await` 解包 `params` 和 `searchParams`
- `layout.tsx` 中的 `params` 也需要异步处理

**示例模式**:
```typescript
// 动态路由页面
export default async function OpportunityPage({
  params,
  searchParams,
}: {
  params: Promise<{ id: string }>;
  searchParams: Promise<{ tab?: string }>;
}) {
  const { id } = await params;
  const { tab } = await searchParams;
  // ... 页面逻辑
}
```

#### 2. React 19 新特性应用

**Server Actions 集成**:
- 表单提交优先使用 Server Actions
- 结合 `useActionState` 处理表单状态和错误
- 使用 `useFormStatus` 显示提交状态

**应用场景**:
| 功能模块 | Actions 文件 | 主要操作 |
|---------|-------------|---------|
| 认证 | `actions/auth.ts` | login, register, logout |
| 邮件 | `actions/email.ts` | sync, delete, archive |
| 组织 | `actions/organization.ts` | create, update, invite |
| 机会 | `actions/opportunities.ts` | create, update, delete |

**use Hook 应用**:
- 在客户端组件中解包 Server Components 传递的 Promise
- 处理异步数据流

#### 3. Tailwind CSS v4 迁移

**主要变更点**:
- 已使用 `@tailwindcss/postcss` v4
- 配置文件保持兼容
- CSS 变量驱动的主题系统继续使用

**色彩系统维护**:
- 保持现有 HSL 变量体系
- 继续支持亮色/暗色模式切换

#### 4. shadcn/ui 全面采用

**组件迁移策略**:
- 使用 `npx shadcn@latest add <component>` 添加缺失组件
- 移除现有 `components/ui/` 中的自定义实现
- 重新组织为标准 shadcn 结构

**核心组件清单**:
- button, input, card, dialog, sheet
- dropdown-menu, tabs, label, checkbox
- toast, skeleton, avatar, badge
- select, textarea, switch, radio-group
- popover, tooltip, accordion, alert
- table, form, separator, scroll-area

#### 5. AI Native 架构核心

**后端 AI API 现状**:

后端已提供完善的 AI 能力,前端直接调用:

| API 端点 | 方法 | 功能 | 实现状态 |
|---------|------|------|----------|
| `/api/v1/ai/draft` | POST | 生成邮件草稿 | ✅ 已实现 |
| `/api/v1/ai/reply` | POST | 生成智能回复 | ✅ 已实现 |
| `/api/v1/chat/completions` | POST | 流式 AI 对话 | ✅ 已实现 (SSE) |
| `/api/v1/search` | GET | 语义搜索 | ✅ 已实现 |

**后端 AI 架构优势**:
- 统一的 AI Provider 抽象层 (支持 OpenAI/Gemini)
- 内置 RAG 上下文注入机制
- SSE (Server-Sent Events) 流式响应
- Prompt 模板集中管理
- AI 成本和配额后端控制

**前端集成策略**:
- 使用 EventSource API 接收 SSE 流
- TanStack Query 管理非流式 AI 请求
- 自定义 Hook 封装流式聊天逻辑
- React 19 `useTransition` 优化 UI 更新

**流式 UI 优化**:
- React Memoization 减少重渲染
- Suspense 边界处理加载状态
- 虚拟滚动优化长对话列表

#### 6. 国际化方案 - 保持现状

**现有实现保留**:
- 继续使用 `LanguageContext` 提供国际化功能
- 保持现有的 `LanguageProvider` 和 `useLanguage` hook
- 保持现有的翻译文件组织结构
- 语言切换器继续使用 Context 状态切换

**不改动的原因**:
- 现有方案稳定可靠
- 避免大规模路由重构
- 降低本次架构升级的复杂度和风险
- 专注于 AI Native 核心功能

**未来改进方向** (可选):
- 在后续版本中可评估迁移到 Next.js Dictionary Pattern
- 当前不影响 AI Native 架构的核心目标

#### 7. 缓存策略设计

**Next.js 16 缓存默认行为**:
- `fetch()` 默认缓存 (force-cache)
- Route Handler 默认静态
- 需要显式配置动态行为

**缓存配置策略**:
| 路由类型 | 缓存策略 | 配置方式 |
|---------|---------|---------|
| 静态页面 | 完全缓存 | 默认行为 |
| 动态页面 | 禁用缓存 | `export const dynamic = 'force-dynamic'` |
| API 路由 | 按需配置 | `export const revalidate = 60` |
| AI 端点 | 禁用缓存 | Edge Runtime + 流式响应 |

**TanStack Query 配置**:
- staleTime: 1分钟 (保持现有配置)
- cacheTime: 5分钟
- refetchOnWindowFocus: true
- retry: 1

#### 8. 后端 API 对接规范

**SSE 流式响应处理**:

后端聊天接口使用 SSE (Server-Sent Events) 协议:

- Content-Type: `text/event-stream`
- 数据格式: `data: {JSON}\n\n`
- 结束标记: `data: [DONE]\n\n`

**EventSource 客户端实现**:

前端使用原生 EventSource API 或自定义 fetch + SSE 解析器:

- 支持认证 Token 注入
- 支持组织 ID Header
- 错误重连机制
- 流式解析和状态更新

**非流式 API 调用**:

草稿生成和回复生成使用标准 REST API:

- 通过 axios 客户端调用
- TanStack Query 管理缓存和状态
- 错误处理和重试策略

## 实施路线图

### Step 1: 项目脚手架与依赖升级

**目标**: 完成 Next.js 16 环境验证和新依赖安装

**任务清单**:
1. 安装缺失的 shadcn/ui 组件
   - 执行 `npx shadcn@latest init` (如果尚未初始化)
   - 逐一添加所需组件
2. 验证 Tailwind CSS v4 配置
3. 创建 `components.json` (shadcn 配置)
4. 验证 TypeScript 配置符合 React 19 要求
5. 创建 `lib/ai/` 目录结构 (后端 API 封装)

**验收标准**:
- 所有依赖安装完成且无冲突
- `pnpm dev` 正常启动
- shadcn/ui 组件可正常导入

**预估工作量**: 0.5 天

---

### Step 2: shadcn/ui 组件迁移

**目标**: 完全替换现有 UI 组件为标准 shadcn/ui 实现

**任务清单**:
1. 备份现有 `components/ui/` 目录为 `components/ui.backup/`
2. 使用 shadcn CLI 重新生成所有 UI 组件
3. 调整组件导入路径 (统一为小写命名)
4. 迁移自定义样式和变体
5. 删除备份目录

**核心组件迁移**:
| 组件名称 | CLI 命令 | 迁移复杂度 |
|---------|---------|-----------|
| Button | `npx shadcn@latest add button` | 低 |
| Input | `npx shadcn@latest add input` | 低 |
| Dialog | `npx shadcn@latest add dialog` | 中 (需迁移 ConfirmDialog) |
| Sheet | `npx shadcn@latest add sheet` | 中 (侧边栏逻辑) |
| Toast | `npx shadcn@latest add toast` | 高 (ToastContainer 重构) |
| Dropdown | `npx shadcn@latest add dropdown-menu` | 中 |
| Tabs | `npx shadcn@latest add tabs` | 低 |

**验收标准**:
- 所有页面组件正常渲染
- 主题切换功能正常
- 无样式破坏或布局错乱

**预估工作量**: 2 天

---

### Step 3: Next.js 16 异步 API 改造

**目标**: 所有页面符合 Next.js 16 的 async params/searchParams 规范

**任务清单**:
1. 审计所有 `page.tsx` 和 `layout.tsx` 文件
2. 重构动态路由页面 (`[id]/page.tsx`)
   - `opportunities/[id]/page.tsx`
   - 其他动态路由 (如果有)
3. 重构搜索页面 (`search/page.tsx`)
   - 处理 `searchParams` Promise
4. 更新布局文件的 `params` 处理
5. 验证所有异步操作符合规范

**改造模式**:
```
// 改造前
export default function Page({ params }: { params: { id: string } }) {
  const { id } = params;
}

// 改造后
export default async function Page({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
}
```

**验收标准**:
- 无 TypeScript 类型错误
- 无运行时警告或错误
- 所有动态路由正常工作
- 搜索参数正确解析

**预估工作量**: 1 天

---

### Step 4: React 19 Server Actions 集成

**目标**: 实现表单和数据变更的 Server Actions 架构

**任务清单**:
1. 创建 `src/actions/` 目录
2. 实现认证 Actions (`auth.ts`)
   - `loginAction`
   - `registerAction`
   - `logoutAction`
3. 实现邮件 Actions (`email.ts`)
   - `syncEmailsAction`
   - `deleteEmailAction`
4. 实现组织 Actions (`organization.ts`)
   - `createOrganizationAction`
   - `updateOrganizationAction`
5. 重构表单组件使用 `useActionState`
   - `LoginForm`
   - `RegisterForm`
   - `OrganizationForm`
6. 集成 `useFormStatus` 显示加载状态

**Actions 实现模式**:
```
'use server'

export async function loginAction(prevState: any, formData: FormData) {
  const email = formData.get('email') as string;
  const password = formData.get('password') as string;
  
  // 验证逻辑
  // API 调用
  // 返回结果或错误
}
```

**表单组件使用模式**:
```
'use client'

const [state, formAction] = useActionState(loginAction, initialState);
const { pending } = useFormStatus();
```

**验收标准**:
- 登录/注册表单正常工作
- 表单提交显示加载状态
- 错误处理正确显示
- 无需页面刷新完成操作

**预估工作量**: 2 天

---

### Step 5: 后端 AI API 集成与 SSE 流处理

**目标**: 实现后端 AI 接口的前端封装和流式处理

**任务清单**:
1. 创建 `lib/ai/` 工具库
   - `chat-client.ts`: SSE 流式聊天客户端
   - `draft-client.ts`: AI 草稿生成客户端
   - `stream-parser.ts`: SSE 流解析器
   - `types.ts`: AI 相关类型定义
2. 实现 SSE 流式处理
   - 基于 fetch API 的 SSE 客户端
   - 支持认证 Token 和组织 Header
   - 错误重连机制
   - 流式数据解析
3. 封装后端 AI API
   - `/api/v1/ai/draft` 封装
   - `/api/v1/ai/reply` 封装
   - `/api/v1/chat/completions` SSE 流封装
4. 创建自定义 Hooks
   - `useStreamChat`: 流式聊天 Hook
   - `useAIDraft`: AI 草稿生成 Hook
   - 集成 TanStack Query
5. 重构 Copilot 组件
   - `chat-interface.tsx`: 使用 `useStreamChat`
   - `message-stream.tsx`: 流式消息渲染
   - 优化渲染性能 (Memoization)
6. 集成 AI Draft 功能
   - 更新 `AIDraftReplyModal` 组件
   - 使用新的 `useAIDraft` Hook

**SSE 客户端示例**:
```typescript
// lib/ai/chat-client.ts
export async function streamChat(
  messages: Message[],
  contextRefIds: string[],
  onChunk: (chunk: ChatChunk) => void,
  onError: (error: Error) => void,
  signal?: AbortSignal
) {
  const token = useAuthStore.getState().token;
  const orgId = useOrganizationStore.getState().currentOrgId;

  const response = await fetch('/api/v1/chat/completions', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
      'X-Organization-ID': orgId,
    },
    body: JSON.stringify({ messages, context_ref_ids: contextRefIds }),
    signal,
  });

  const reader = response.body?.getReader();
  const decoder = new TextDecoder();

  while (true) {
    const { done, value } = await reader.read();
    if (done) break;

    const chunk = decoder.decode(value);
    const lines = chunk.split('\n\n');

    for (const line of lines) {
      if (line.startsWith('data: ')) {
        const data = line.slice(6);
        if (data === '[DONE]') return;
        
        try {
          const parsed = JSON.parse(data);
          onChunk(parsed);
        } catch (e) {
          onError(new Error('Failed to parse SSE data'));
        }
      }
    }
  }
}
```

**验收标准**:
- AI 聊天功能正常工作
- 流式响应实时显示
- AI 草稿生成功能可用
- 上下文正确注入
- 错误处理完善
- SSE 重连机制生效

**预估工作量**: 2.5 天

---

### Step 6: 性能优化与测试验证

**目标**: 确保架构重构后的性能和稳定性

**任务清单**:
1. 实施 React 19 优化策略
   - 高频更新组件使用 `memo`
   - AI 流式渲染使用 `useTransition`
   - Suspense 边界配置
2. 配置 Next.js 缓存策略
   - 静态页面缓存
   - 动态路由缓存配置
   - API 路由 revalidate 设置
3. TanStack Query 优化
   - 调整 staleTime 和 cacheTime
   - 配置预取策略
   - 乐观更新实现
4. 执行端到端测试
   - 认证流程测试
   - AI 功能测试
   - 国际化切换测试
   - 响应式布局测试
5. 性能基准测试
   - Lighthouse 评分
   - 首屏加载时间
   - AI 响应延迟
6. 代码审查与清理
   - 移除未使用的依赖
   - 删除备份文件
   - 代码格式化

**性能目标**:
| 指标 | 目标值 |
|-----|-------|
| Lighthouse Performance | ≥ 90 |
| First Contentful Paint | ≤ 1.5s |
| Time to Interactive | ≤ 3s |
| AI 首字节响应时间 | ≤ 500ms |
| 语言切换延迟 | ≤ 200ms |

**验收标准**:
- 所有 E2E 测试通过
- 性能指标达标
- 无 TypeScript 错误
- 无 ESLint 警告
- 代码库整洁

**预估工作量**: 2 天

---

## 风险评估与缓解策略

| 风险项 | 影响等级 | 概率 | 缓解措施 |
|-------|---------|------|----------|
| shadcn 组件样式不兼容 | 中 | 中 | 渐进式迁移,保留备份,逐个组件验证 |
| SSE 流处理兼容性 | 中 | 低 | 提前测试流式解析,实现重连机制 |
| Server Actions 性能问题 | 低 | 低 | 监控响应时间,优化后端 API 调用 |
| 现有功能回归 | 高 | 中 | 完善 E2E 测试覆盖,分阶段发布 |
| React 19 兼容性问题 | 中 | 低 | 验证第三方库兼容性,准备降级方案 |

## 总工作量估算

| 步骤 | 预估工作量 | 依赖关系 |
|-----|-----------|---------||
| Step 1: 脚手架与依赖 | 0.5 天 | - |
| Step 2: shadcn/ui 迁移 | 2 天 | Step 1 |
| Step 3: 异步 API 改造 | 1 天 | Step 2 |
| Step 4: Server Actions | 2 天 | Step 3 |
| Step 5: 后端 AI API 集成 | 2.5 天 | Step 4 |
| Step 6: 优化与测试 | 2 天 | Step 5 |
| **总计** | **10 天** | - |

## 成功指标

### 功能完整性
- ✅ 所有现有功能正常工作
- ✅ AI 聊天和草稿生成功能可用
- ✅ 国际化功能保持正常
- ✅ 主题切换正常

### 技术合规性
- ✅ 符合 Next.js 16 异步 API 规范
- ✅ 使用 React 19 新特性 (Actions, useActionState)
- ✅ 完全采用 shadcn/ui 组件
- ✅ 直接调用后端 AI API

### 性能表现
- ✅ Lighthouse 性能评分 ≥ 90
- ✅ AI 响应流畅无卡顿
- ✅ 页面切换延迟 ≤ 200ms

### 代码质量
- ✅ 无 TypeScript 类型错误
- ✅ 无 ESLint 警告
- ✅ 代码覆盖率保持或提升
- ✅ 组件可复用性提高

## 后端 AI 能力改进计划

在前端重构过程中,后端已提供了基础的 AI 能力,但存在以下优化空间,建议按优先级逐步改进:

### 高优先级改进 (P0)

#### 1. AI 聊天流式响应增强

**现状问题**:
- SSE 流缺乏结构化 Chunk 类型区分
- 无法传递 Token 使用统计信息
- 不支持流式中断和重连

**改进方案**:

扩展 `ChatCompletionChunk` 类型:

| 字段 | 类型 | 说明 |
|------|------|------|
| id | string | Chunk 唯一 ID |
| choices | array | 选项列表 |
| usage | object | Token 使用统计 (最后一个 chunk) |
| finish_reason | string | 结束原因 (stop/length/error) |
| metadata | object | 扩展元数据 |

**实现路径**:
- 修改 `backend/pkg/ai/provider.go` 中的 `ChatCompletionChunk` 定义
- 更新 `openai/provider.go` 和 `gemini/provider.go` 的 `StreamChat` 实现
- 前端同步更新类型定义

**预估工作量**: 1 天

---

#### 2. 上下文管理 API 增强

**现状问题**:
- 聊天接口 `context_ref_ids` 只支持邮件 ID
- 无法灵活注入自定义上下文 (Context 实体)
- 缺乏上下文优先级机制

**改进方案**:

扩展 `ChatRequest` 结构:

| 字段 | 类型 | 说明 |
|------|------|------|
| messages | []Message | 对话历史 |
| context_sources | []ContextSource | 统一上下文来源 |
| max_context_tokens | int | 上下文最大 Token 数 |
| stream | bool | 是否流式输出 |

**ContextSource 类型定义**:

| 字段 | 类型 | 可选值 |
|------|------|--------|
| type | string | email / context / document / search |
| id | string | 资源 ID |
| priority | int | 优先级 (1-10) |
| metadata | object | 扩展元数据 |

**后端处理逻辑**:
1. 按 priority 排序上下文源
2. 根据 `max_context_tokens` 截断过长内容
3. 支持多种来源类型的上下文注入

**实现路径**:
- 更新 `backend/internal/handler/chat.go` 的 `ChatRequest`
- 修改 `backend/internal/service/chat.go` 的 `StreamChat` 逻辑
- 增加 `ContextSource` 解析器

**预估工作量**: 1.5 天

---

#### 3. AI 草稿生成批量 API

**现状问题**:
- 只支持单个草稿生成
- 用户需要批量生成回复时效率低

**改进方案**:

新增 API 端点: `POST /api/v1/ai/batch-reply`

**请求结构**:

| 字段 | 类型 | 说明 |
|------|------|------|
| email_ids | []string | 邮件 ID 列表 |
| tone | string | 统一语气 |
| context | string | 统一上下文 |
| parallel | bool | 是否并行生成 |

**响应结构**:

| 字段 | 类型 | 说明 |
|------|------|------|
| results | []DraftResult | 生成结果列表 |
| total | int | 总数 |
| successful | int | 成功数 |
| failed | int | 失败数 |

**实现路径**:
- 新增 `backend/internal/handler/ai_draft.go` 中的 `BatchGenerateReply`
- 使用 goroutine 池并行处理
- 增加 Rate Limiting 防止滥用

**预估工作量**: 1 天

---

### 中优先级改进 (P1)

#### 4. AI 响应缓存机制

**现状问题**:
- 相同请求重复调用 AI 接口
- 增加成本和响应延迟

**改进方案**:

实现基于 Redis 的 AI 响应缓存:

**缓存策略**:

| 场景 | 缓存键 | TTL |
|------|---------|-----|
| 草稿生成 | `ai:draft:{emailID}:{tone}:{context}` | 1 小时 |
| 搜索结果 | `search:{query}:{hash}` | 15 分钟 |
| 聊天回复 | 不缓存 (上下文敏感) | - |

**实现路径**:
- 在 `backend/internal/service/ai_draft.go` 增加缓存层
- 使用 Redis 存储缓存结果
- 支持缓存失效策略

**预估工作量**: 0.5 天

---

#### 5. Prompt 模板管理 UI

**现状问题**:
- Prompt 模板硬编码在配置文件
- 无法动态调整和 A/B 测试

**改进方案**:

新增 Prompt 管理 API:

| 端点 | 方法 | 功能 |
|------|------|------|
| `/api/v1/admin/prompts` | GET | 列表所有 Prompt 模板 |
| `/api/v1/admin/prompts` | POST | 创建新模板 |
| `/api/v1/admin/prompts/:id` | PUT | 更新模板 |
| `/api/v1/admin/prompts/:id/activate` | POST | 激活模板 |

**Prompt 模板数据模型**:

| 字段 | 类型 | 说明 |
|------|------|------|
| id | string | 模板 ID |
| name | string | 模板名称 |
| type | string | 类型 (summary/classify/draft) |
| template | string | Prompt 模板内容 |
| variables | []string | 可用变量 |
| is_active | bool | 是否启用 |
| version | int | 版本号 |

**实现路径**:
- 新增 `backend/internal/model/prompt.go`
- 新增 `backend/internal/handler/prompt.go`
- 前端增加 Prompt 管理页面

**预估工作量**: 2 天

---

### 低优先级改进 (P2)

#### 6. AI 成本监控与报表

**现状问题**:
- 无 AI 使用量统计
- 无法追踪成本

**改进方案**:

新增 AI 使用统计表 `ai_usage_logs`:

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uuid | 记录 ID |
| user_id | uuid | 用户 ID |
| organization_id | uuid | 组织 ID |
| service_type | string | 服务类型 (chat/draft/search) |
| provider | string | AI 提供商 |
| model | string | 模型名称 |
| prompt_tokens | int | Prompt Token 数 |
| completion_tokens | int | 输出 Token 数 |
| total_tokens | int | 总 Token 数 |
| estimated_cost | decimal | 估算成本 |
| created_at | timestamp | 创建时间 |

**统计 API**:

| 端点 | 方法 | 功能 |
|------|------|------|
| `/api/v1/admin/ai-usage/stats` | GET | 获取使用统计 |
| `/api/v1/admin/ai-usage/cost` | GET | 获取成本报表 |

**预估工作量**: 1.5 天

---

#### 7. 多模态 AI 支持

**现状问题**:
- 只支持文本输入
- 无法处理图片/附件

**改进方案**:

支持多模态输入 (GPT-4V/Gemini Pro Vision):

**扩展请求结构**:

| 字段 | 类型 | 说明 |
|------|------|------|
| content | []ContentPart | 多模态内容 |

**ContentPart 类型**:

| 字段 | 类型 | 说明 |
|------|------|------|
| type | string | text / image / file |
| text | string | 文本内容 |
| image_url | string | 图片 URL |
| mime_type | string | MIME 类型 |

**应用场景**:
- 分析邮件中的图片附件
- PDF/文档内容提取
- 截图分析

**预估工作量**: 3 天

---

## 后端改进总结

| 优先级 | 改进项 | 工作量 | 前端依赖 |
|---------|---------|---------|----------|
| P0 | AI 聊天流式响应增强 | 1 天 | 是 |
| P0 | 上下文管理 API 增强 | 1.5 天 | 是 |
| P0 | AI 草稿生成批量 API | 1 天 | 否 |
| P1 | AI 响应缓存机制 | 0.5 天 | 否 |
| P1 | Prompt 模板管理 UI | 2 天 | 是 |
| P2 | AI 成本监控与报表 | 1.5 天 | 是 |
| P2 | 多模态 AI 支持 | 3 天 | 是 |
| **总计** | - | **10.5 天** | - |

**建议实施顺序**:
1. **阶段 1** (P0): 与前端重构同步进行,确保 API 兼容性
2. **阶段 2** (P1): 前端重构完成后进行,提升用户体验
3. **阶段 3** (P2): 根据业务需求逐步实施

---

## 执行检查清单

### 阶段 0: 前期准备 (开始前必须完成)

- [ ] 团队对设计方案达成一致
- [ ] 确认后端 AI API 接口稳定可用
- [ ] 备份现有代码库 (创建分支 `backup/pre-refactor`)
- [ ] 准备测试环境
- [ ] 通知相关利益方 (产品/运维/测试)

### 阶段 1: 项目脚手架与依赖升级 (0.5 天)

- [ ] 创建开发分支 `feat/ai-native-refactor`
- [ ] 初始化 shadcn/ui (`npx shadcn@latest init`)
- [ ] 验证 Tailwind CSS v4 配置
- [ ] 创建 `lib/ai/` 目录结构
- [ ] 更新 TypeScript 配置 (如需)
- [ ] 执行 `pnpm install` 并验证启动
- [ ] 提交代码: `git commit -m "chore: setup scaffolding"`

### 阶段 2: shadcn/ui 组件迁移 (2 天)

**第 1 天**:
- [ ] 备份现有 `components/ui/` 为 `components/ui.backup/`
- [ ] 添加基础组件: button, input, card, label
- [ ] 添加表单组件: checkbox, select, textarea
- [ ] 调整导入路径
- [ ] 测试基础功能正常
- [ ] 提交代码: `git commit -m "feat: migrate basic ui components"`

**第 2 天**:
- [ ] 添加复杂组件: dialog, sheet, dropdown-menu, tabs
- [ ] 重构 Toast 系统 (`toast.tsx` + `use-toast.ts`)
- [ ] 迁移 ConfirmDialog 逻辑
- [ ] 全面测试所有页面
- [ ] 修复样式问题
- [ ] 删除备份目录
- [ ] 提交代码: `git commit -m "feat: complete shadcn ui migration"`

### 阶段 3: Next.js 16 异步 API 改造 (1 天)

- [ ] 审计所有 `page.tsx` 和 `layout.tsx`
- [ ] 重构动态路由 (`opportunities/[id]/page.tsx`)
- [ ] 重构搜索页面 (`search/page.tsx`)
- [ ] 更新布局文件的 params 处理
- [ ] 验证 TypeScript 类型正确
- [ ] 测试所有路由正常工作
- [ ] 提交代码: `git commit -m "feat: adopt nextjs 16 async api"`

### 阶段 4: React 19 Server Actions 集成 (2 天)

**第 1 天**:
- [ ] 创建 `src/actions/` 目录
- [ ] 实现 `actions/auth.ts` (loginAction, registerAction)
- [ ] 重构 LoginForm 使用 useActionState
- [ ] 重构 RegisterForm 使用 useActionState
- [ ] 集成 useFormStatus 显示加载状态
- [ ] 测试认证流程
- [ ] 提交代码: `git commit -m "feat: implement auth server actions"`

**第 2 天**:
- [ ] 实现 `actions/email.ts` (syncEmailsAction, deleteEmailAction)
- [ ] 实现 `actions/organization.ts`
- [ ] 更新相关表单组件
- [ ] 全面测试表单提交
- [ ] 验证错误处理
- [ ] 提交代码: `git commit -m "feat: complete server actions integration"`

### 阶段 5: 后端 AI API 集成与 SSE 流处理 (2.5 天)

**第 1 天**:
- [ ] 创建 `lib/ai/chat-client.ts` (SSE 客户端)
- [ ] 创建 `lib/ai/draft-client.ts`
- [ ] 创建 `lib/ai/stream-parser.ts`
- [ ] 创建 `lib/ai/types.ts`
- [ ] 测试 SSE 连接和解析
- [ ] 提交代码: `git commit -m "feat: implement sse client"`

**第 2 天**:
- [ ] 封装后端 AI API (`/api/v1/ai/draft`, `/api/v1/ai/reply`)
- [ ] 创建 `useStreamChat` Hook
- [ ] 创建 `useAIDraft` Hook
- [ ] 集成 TanStack Query
- [ ] 测试 AI 草稿生成
- [ ] 提交代码: `git commit -m "feat: implement ai hooks"`

**第 3 天 (0.5)**:
- [ ] 重构 Copilot 组件
  - [ ] 更新 `chat-interface.tsx`
  - [ ] 更新 `message-stream.tsx`
  - [ ] 优化渲染性能 (Memoization)
- [ ] 更新 `AIDraftReplyModal` 组件
- [ ] 全面测试 AI 功能
- [ ] 测试 SSE 重连机制
- [ ] 提交代码: `git commit -m "feat: complete ai integration"`

### 阶段 6: 性能优化与测试验证 (2 天)

**第 1 天**:
- [ ] 实施 React 19 优化 (memo, useTransition)
- [ ] 配置 Next.js 缓存策略
- [ ] 优化 TanStack Query 配置
- [ ] 运行 Lighthouse 评分
- [ ] 测量首屏加载时间
- [ ] 测量 AI 响应延迟
- [ ] 提交代码: `git commit -m "perf: apply performance optimizations"`

**第 2 天**:
- [ ] 执行 E2E 测试
  - [ ] 认证流程测试
  - [ ] AI 功能测试
  - [ ] 国际化功能测试
  - [ ] 响应式布局测试
- [ ] 代码审查与清理
- [ ] 移除未使用的依赖
- [ ] 删除备份文件
- [ ] 代码格式化
- [ ] 更新文档
- [ ] 提交代码: `git commit -m "test: complete e2e testing"`

### 阶段 8: 发布上线 (预留)

- [ ] 合并到主分支 (`git merge feat/ai-native-refactor`)
- [ ] 创建版本标签 (`git tag v2.0.0`)
- [ ] 部署到测试环境
- [ ] 执行烟雾测试
- [ ] 通知用户升级
- [ ] 监控系统表现
- [ ] 收集用户反馈
- [ ] 部署到生产环境

---

## 质量保证清单

### 代码质量
- [ ] 无 TypeScript 类型错误
- [ ] 无 ESLint 警告
- [ ] 所有组件有适当的注释
- [ ] 关键逻辑有单元测试

### 功能完整性
- [ ] 所有现有功能正常工作
- [ ] AI 聊天功能可用
- [ ] AI 草稿生成功能可用
- [ ] 国际化功能正常
- [ ] 主题切换正常

### 性能指标
- [ ] Lighthouse Performance ≥ 90
- [ ] First Contentful Paint ≤ 1.5s
- [ ] Time to Interactive ≤ 3s
- [ ] AI 首字节响应时间 ≤ 500ms
- [ ] 语言切换延迟 ≤ 200ms

### 浏览器兼容性
- [ ] Chrome 最新版
- [ ] Firefox 最新版
- [ ] Safari 最新版
- [ ] Edge 最新版
- [ ] 移动端浏览器

### 响应式设计
- [ ] 桌面端 (>= 1024px)
- [ ] 平板 (768px - 1023px)
- [ ] 手机 (< 768px)

---

## 应急预案

### 预案 1: shadcn 组件样式问题

**触发条件**: shadcn 组件样式与现有设计冲突

**应对措施**:
1. 回滚到备份的 UI 组件
2. 修改 shadcn 组件的样式配置
3. 调整 Tailwind 主题变量

### 预案 2: SSE 流处理失败

**触发条件**: SSE 连接不稳定或解析失败

**应对措施**:
1. 启用重连机制 (exponential backoff)
2. 降级为轮询模式
3. 显示错误提示给用户

### 预案 3: 性能回退

**触发条件**: 重构后性能指标未达标

**应对措施**:
1. 分析性能瓶颈
2. 启用代码分割 (Code Splitting)
3. 优化图片和资源加载
4. 实施预加载策略
