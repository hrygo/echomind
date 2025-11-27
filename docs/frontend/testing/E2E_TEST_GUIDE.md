# 端到端功能测试指南

## 测试概述

本文档提供了AI Native架构重构后的端到端功能测试指南。测试覆盖以下核心功能模块：

1. **Server Actions** - React 19 Server Actions功能
2. **AI流式聊天** - SSE流式聊天功能
3. **AI草稿生成** - 邮件草稿和回复生成功能
4. **自定义Hooks** - AI相关自定义Hooks

## 测试文件结构

```
tests/e2e/
├── server-actions.spec.ts    # Server Actions 测试
├── ai-streaming.spec.ts      # AI 流式聊天测试
└── ai-draft.spec.ts          # AI 草稿生成测试

src/hooks/
└── useAI.test.tsx            # AI Hooks 单元测试
```

## 前置条件

### 1. 安装依赖

```bash
cd frontend
pnpm install
```

### 2. 环境配置

确保以下环境变量已配置（`.env.local`）：

```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
NEXT_PUBLIC_AI_STREAM_ENDPOINT=/api/v1/ai/chat/stream
```

### 3. 启动后端服务

确保后端服务运行在 `http://localhost:8080`：

```bash
# 在后端目录
cd backend
go run main.go
```

## 运行测试

### 运行所有 E2E 测试

```bash
pnpm playwright test
```

### 运行特定测试套件

```bash
# Server Actions 测试
pnpm playwright test tests/e2e/server-actions.spec.ts

# AI 流式聊天测试
pnpm playwright test tests/e2e/ai-streaming.spec.ts

# AI 草稿生成测试
pnpm playwright test tests/e2e/ai-draft.spec.ts
```

### 运行单元测试

```bash
# 运行所有单元测试
pnpm test

# 运行 AI Hooks 测试
pnpm test useAI.test.tsx

# 监视模式
pnpm test -- --watch
```

### 调试模式

```bash
# UI 模式运行测试
pnpm playwright test --ui

# 调试特定测试
pnpm playwright test --debug tests/e2e/server-actions.spec.ts
```

## 测试详细说明

### 1. Server Actions 测试 (`server-actions.spec.ts`)

测试 React 19 Server Actions 集成：

#### 认证功能
- ✅ 登录 Action
- ✅ 注册 Action
- ✅ 登出 Action
- ✅ 表单验证错误处理

#### 邮件操作
- ✅ 同步邮件
- ✅ 删除邮件
- ✅ 归档邮件

#### 组织管理
- ✅ 创建组织
- ✅ 更新组织
- ✅ 邀请成员

**关键验证点：**
- Server Actions 正确触发
- 表单数据正确提交
- 错误状态正确处理
- 缓存重新验证（revalidatePath）
- Cookie 管理（认证令牌）

### 2. AI 流式聊天测试 (`ai-streaming.spec.ts`)

测试 SSE 流式聊天功能：

#### 基础流式功能
- ✅ 流式响应接收
- ✅ 流式指示器显示
- ✅ 多条消息流式传输
- ✅ 流式错误处理
- ✅ 取消流式传输

#### 上下文感知
- ✅ 对话上下文保持
- ✅ 系统上下文集成

#### 性能测试
- ✅ 首字节延迟（TTFB）
- ✅ 快速连续消息处理

**关键验证点：**
- SSE 连接正确建立
- 流式数据正确解析
- 消息增量更新
- 错误重连机制
- AbortController 取消机制

### 3. AI 草稿生成测试 (`ai-draft.spec.ts`)

测试 AI 草稿和回复生成功能：

#### 草稿生成
- ✅ 新邮件草稿生成
- ✅ 不同语气风格
- ✅ 重新生成草稿
- ✅ 生成进度显示
- ✅ 错误处理
- ✅ 保存草稿

#### 回复生成
- ✅ 邮件回复生成
- ✅ 自定义指令
- ✅ 不同回复语气
- ✅ 上下文包含
- ✅ 发送AI生成回复
- ✅ 编辑后发送

#### 高级特性
- ✅ 附件上下文
- ✅ 多语言生成
- ✅ 草稿建议

**关键验证点：**
- API 请求正确发送
- 生成内容质量验证
- 加载状态管理
- 用户交互流程
- 错误恢复机制

### 4. AI Hooks 单元测试 (`useAI.test.tsx`)

测试自定义 AI Hooks：

#### useStreamChat
- ✅ 初始化状态
- ✅ 发送消息
- ✅ 流式响应处理
- ✅ 错误处理
- ✅ 取消功能
- ✅ 消息历史

#### useAIDraft
- ✅ 草稿生成
- ✅ 不同语气
- ✅ 错误处理
- ✅ 回调函数
- ✅ 取消操作

#### useAIReply
- ✅ 回复生成
- ✅ 自定义上下文
- ✅ 错误处理
- ✅ AbortSignal 支持

**关键验证点：**
- Hook 状态管理
- TanStack Query 集成
- 错误边界处理
- 异步操作管理

## 测试数据

### 测试账号

```
Email: test@example.com
Password: password123
```

### 测试选择器 (data-testid)

确保在组件中添加以下测试选择器：

```tsx
// 认证
<input data-testid="email-input" name="email" />
<input data-testid="password-input" name="password" />
<button data-testid="login-button" type="submit">登录</button>

// 聊天
<div data-testid="chat-input" />
<button data-testid="send-button">发送</button>
<div data-testid="user-message" />
<div data-testid="ai-message" />
<div data-testid="streaming-indicator" />

// 邮件
<button data-testid="sync-emails-button">同步</button>
<div data-testid="email-item" />
<button data-testid="delete-email-button">删除</button>
<button data-testid="archive-email-button">归档</button>

// 草稿
<button data-testid="compose-email-button">撰写</button>
<input data-testid="draft-prompt" />
<button data-testid="generate-draft-button">生成草稿</button>
<div data-testid="draft-content" />
```

## 常见问题

### 1. 测试超时

如果测试经常超时，增加超时时间：

```typescript
test('my test', async ({ page }) => {
  // 增加特定操作的超时时间
  await page.waitForSelector('[data-testid="element"]', { timeout: 15000 });
});
```

或在配置中全局设置：

```typescript
// playwright.config.ts
export default defineConfig({
  timeout: 30000, // 全局测试超时 30 秒
});
```

### 2. 后端服务未运行

确保后端服务正在运行并可访问：

```bash
curl http://localhost:8080/health
```

### 3. 环境变量未设置

检查 `.env.local` 文件是否存在并包含必要的环境变量。

### 4. 测试数据清理

某些测试可能需要清理测试数据。考虑使用测试夹具：

```typescript
test.beforeEach(async () => {
  // 清理测试数据
  await cleanupTestData();
});
```

## 持续集成 (CI)

### GitHub Actions 示例

```yaml
name: E2E Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: 20
      
      - name: Install dependencies
        run: cd frontend && pnpm install
      
      - name: Install Playwright
        run: cd frontend && pnpm exec playwright install --with-deps
      
      - name: Start backend service
        run: |
          cd backend
          go run main.go &
          sleep 5
      
      - name: Run E2E tests
        run: cd frontend && pnpm playwright test
      
      - uses: actions/upload-artifact@v3
        if: always()
        with:
          name: playwright-report
          path: frontend/playwright-report/
```

## 测试报告

测试完成后，查看测试报告：

```bash
# 查看 Playwright 报告
pnpm playwright show-report

# Jest 覆盖率报告
pnpm test -- --coverage
```

## 下一步

1. **集成到 CI/CD**：将测试集成到持续集成流程中
2. **增加覆盖率**：为更多功能添加测试
3. **性能基准**：建立性能基准测试
4. **视觉回归**：添加视觉回归测试
5. **API 模拟**：使用 MSW 进行 API 模拟测试

## 参考资源

- [Playwright 文档](https://playwright.dev/)
- [Jest 文档](https://jestjs.io/)
- [Testing Library](https://testing-library.com/)
- [TanStack Query Testing](https://tanstack.com/query/latest/docs/framework/react/guides/testing)
