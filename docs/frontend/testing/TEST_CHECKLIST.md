# 端到端功能测试验证清单

## ✅ 已完成的工作

### 📋 测试文件创建
- [x] **E2E 测试套件**
  - [x] `tests/e2e/server-actions.spec.ts` - Server Actions 测试 (260 行)
  - [x] `tests/e2e/ai-streaming.spec.ts` - AI 流式聊天测试 (294 行)
  - [x] `tests/e2e/ai-draft.spec.ts` - AI 草稿生成测试 (409 行)

- [x] **单元测试**
  - [x] `src/hooks/useAI.test.tsx` - AI Hooks 单元测试 (280 行)

- [x] **测试文档**
  - [x] `E2E_TEST_GUIDE.md` - 完整测试指南 (364 行)
  - [x] `TEST_VERIFICATION_SUMMARY.md` - 验证总结文档 (399 行)
  - [x] `run-tests.sh` - 自动化测试脚本 (212 行)

### 🏗️ 构建验证
- [x] 项目构建成功（`pnpm build`）
- [x] TypeScript 编译通过
- [x] 无类型错误
- [x] 所有路由正确生成

### 📊 测试覆盖
- [x] **37 个 E2E 测试用例**
  - [x] 12 个 Server Actions 测试
  - [x] 10 个 AI 流式聊天测试
  - [x] 15 个 AI 草稿生成测试

- [x] **19 个单元测试用例**
  - [x] 7 个 useStreamChat 测试
  - [x] 6 个 useAIDraft 测试
  - [x] 5 个 useAIReply 测试
  - [x] 1 个集成测试

---

## 📝 后续执行清单

### 🔴 立即执行（实际运行测试）

#### 1. 环境准备
- [ ] 确保后端服务运行 (`http://localhost:8080`)
- [ ] 配置环境变量 (`.env.local`)
  ```env
  NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
  NEXT_PUBLIC_AI_STREAM_ENDPOINT=/api/v1/ai/chat/stream
  ```
- [ ] 创建测试账号
  ```
  Email: test@example.com
  Password: password123
  ```

#### 2. 添加测试选择器
在相关组件中添加 `data-testid` 属性：

- [ ] **认证组件** (`src/components/auth/*`)
  ```tsx
  <input data-testid="email-input" name="email" />
  <input data-testid="password-input" name="password" />
  <button data-testid="login-button" type="submit">登录</button>
  ```

- [ ] **聊天组件** (`src/components/chat/*` 或 `src/app/copilot/*`)
  ```tsx
  <div data-testid="chat-input" />
  <button data-testid="send-button">发送</button>
  <div data-testid="user-message" />
  <div data-testid="ai-message" />
  <div data-testid="streaming-indicator" />
  <button data-testid="cancel-streaming-button">取消</button>
  ```

- [ ] **邮件组件** (`src/components/email/*` 或 `src/app/dashboard/inbox/*`)
  ```tsx
  <button data-testid="sync-emails-button">同步</button>
  <div data-testid="email-item" />
  <button data-testid="delete-email-button">删除</button>
  <button data-testid="archive-email-button">归档</button>
  <div data-testid="email-detail" />
  <div data-testid="email-subject" />
  ```

- [ ] **草稿组件** (`src/components/email/AIDraftReplyModal.tsx` 等)
  ```tsx
  <button data-testid="compose-email-button">撰写</button>
  <div data-testid="compose-dialog" />
  <input data-testid="email-to" name="to" />
  <input data-testid="email-subject" name="subject" />
  <input data-testid="draft-prompt" />
  <button data-testid="generate-draft-button">生成草稿</button>
  <div data-testid="draft-content" />
  <button data-testid="save-draft-button">保存</button>
  <select data-testid="tone-select" />
  ```

- [ ] **回复组件**
  ```tsx
  <button data-testid="reply-button">回复</button>
  <div data-testid="reply-dialog" />
  <button data-testid="ai-reply-button">AI 回复</button>
  <textarea data-testid="reply-content" />
  <input data-testid="reply-instructions" />
  <select data-testid="reply-tone-select" />
  <button data-testid="send-reply-button">发送</button>
  ```

#### 3. 运行测试

- [ ] **运行单元测试**
  ```bash
  cd frontend
  pnpm test
  ```
  预期：所有单元测试通过

- [ ] **运行 E2E 测试（分步）**
  ```bash
  # Server Actions 测试
  pnpm playwright test tests/e2e/server-actions.spec.ts
  
  # AI 流式聊天测试
  pnpm playwright test tests/e2e/ai-streaming.spec.ts
  
  # AI 草稿生成测试
  pnpm playwright test tests/e2e/ai-draft.spec.ts
  ```

- [ ] **运行所有测试**
  ```bash
  ./run-tests.sh
  ```

#### 4. 查看测试报告

- [ ] 查看 Playwright 报告
  ```bash
  pnpm playwright show-report
  ```

- [ ] 查看 Jest 覆盖率
  ```bash
  pnpm test -- --coverage
  ```

- [ ] 记录失败的测试用例

---

### 🟡 中期改进（1-2 周内）

#### 测试修复和优化
- [ ] 修复失败的测试用例
- [ ] 优化慢速测试
- [ ] 增加测试超时配置
- [ ] 处理测试中的竞态条件

#### 测试数据管理
- [ ] 创建测试数据夹具（fixtures）
- [ ] 实现测试数据清理脚本
- [ ] 添加测试数据工厂（factories）

#### 测试增强
- [ ] 添加更多边界情况测试
- [ ] 增加网络错误模拟测试
- [ ] 添加并发操作测试
- [ ] 实现 API Mock（使用 MSW）

---

### 🟢 长期优化（持续进行）

#### CI/CD 集成
- [ ] 配置 GitHub Actions
  ```yaml
  name: E2E Tests
  on: [push, pull_request]
  jobs:
    test:
      runs-on: ubuntu-latest
      steps:
        - name: Checkout
        - name: Setup Node.js
        - name: Install dependencies
        - name: Run tests
        - name: Upload reports
  ```

- [ ] 添加 pre-commit hooks
- [ ] 配置自动化部署前测试

#### 测试监控
- [ ] 设置测试失败告警
- [ ] 监控测试执行时间
- [ ] 跟踪测试覆盖率趋势
- [ ] 建立性能基准

#### 文档和培训
- [ ] 编写测试最佳实践文档
- [ ] 创建测试编写指南
- [ ] 团队测试培训
- [ ] 测试代码审查标准

#### 持续改进
- [ ] 定期审查测试有效性
- [ ] 移除冗余测试
- [ ] 优化测试执行速度
- [ ] 更新测试以适应新功能

---

## 🎯 测试执行优先级

### P0 - 核心功能（必须通过）
1. ✅ 项目构建
2. ⏳ Server Actions - 认证功能
3. ⏳ Server Actions - 邮件基本操作
4. ⏳ AI 流式聊天 - 基础功能

### P1 - 重要功能（应该通过）
5. ⏳ AI 草稿生成 - 基础功能
6. ⏳ AI 回复生成 - 基础功能
7. ⏳ useStreamChat Hook
8. ⏳ useAIDraft Hook

### P2 - 增强功能（最好通过）
9. ⏳ AI 流式聊天 - 上下文感知
10. ⏳ AI 草稿生成 - 高级特性
11. ⏳ Server Actions - 组织管理
12. ⏳ useAIReply Hook

---

## 📈 成功标准

### 最低标准（MVP）
- ✅ 项目成功构建
- ⏳ P0 核心功能测试全部通过
- ⏳ 单元测试覆盖率 > 60%
- ⏳ 关键路径 E2E 测试通过

### 良好标准
- ⏳ P0 + P1 测试全部通过
- ⏳ 单元测试覆盖率 > 75%
- ⏳ 所有 E2E 测试通过率 > 90%
- ⏳ 无关键性能问题

### 优秀标准
- ⏳ 所有测试全部通过
- ⏳ 单元测试覆盖率 > 85%
- ⏳ E2E 测试通过率 100%
- ⏳ CI/CD 集成完成
- ⏳ 性能指标达标

---

## 🐛 问题追踪

### 已知问题
- [ ] 无（待测试执行后填写）

### 待解决问题
- [ ] 无（待测试执行后填写）

---

## 📞 支持和资源

### 文档链接
- [Playwright 官方文档](https://playwright.dev/)
- [Jest 官方文档](https://jestjs.io/)
- [Testing Library](https://testing-library.com/)
- [TanStack Query Testing](https://tanstack.com/query/latest/docs/framework/react/guides/testing)

### 项目文档
- `E2E_TEST_GUIDE.md` - 详细测试指南
- `TEST_VERIFICATION_SUMMARY.md` - 验证总结
- `AI_NATIVE_ARCHITECTURE.md` - 架构文档

### 快速命令参考
```bash
# 安装依赖
pnpm install

# 类型检查
pnpm type-check

# 单元测试
pnpm test
pnpm test -- --watch

# E2E 测试
pnpm playwright test
pnpm playwright test --ui
pnpm playwright test --debug

# 查看报告
pnpm playwright show-report

# 完整测试流程
./run-tests.sh
```

---

## ✅ 验证签字

- **测试框架搭建**: ✅ 已完成
- **测试文件创建**: ✅ 已完成  
- **测试文档编写**: ✅ 已完成
- **构建验证**: ✅ 已通过

**下一步**: 执行实际测试，验证功能正确性

---

_最后更新: 2025-11-27_
