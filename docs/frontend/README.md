# Frontend 文档目录

本目录包含EchoMind前端应用的相关文档。

## 📚 文档索引

### 架构设计

- [AI Native 架构设计](../architecture/AI_NATIVE_ARCHITECTURE.md) - AI Native架构的完整设计文档，包括Next.js 16和React 19的现代化实践

### 测试文档

测试相关文档位于 `testing/` 子目录：

- [端到端测试指南](./testing/E2E_TEST_GUIDE.md) - E2E测试的完整指南
- [测试检查清单](./testing/TEST_CHECKLIST.md) - 功能测试检查清单
- [测试验证摘要](./testing/TEST_VERIFICATION_SUMMARY.md) - 测试验证结果摘要

## 🧪 测试执行

### 快速开始

```bash
# 在项目根目录执行
make test-e2e

# 或直接运行测试脚本
bash scripts/frontend/run-tests.sh
```

### 单独测试

```bash
# 单元测试
cd frontend && pnpm test

# E2E测试
cd frontend && pnpm playwright test

# 类型检查
cd frontend && pnpm type-check
```

## 🏗️ 开发指南

### 技术栈

- **Next.js 16.0.3**: App Router + 异步API
- **React 19.2.0**: Server Components + Server Actions
- **TypeScript**: 严格类型检查
- **Tailwind CSS v4**: 原生CSS + PostCSS
- **shadcn/ui**: 基于Radix UI的组件库
- **TanStack Query v5**: 服务端状态管理
- **Zustand v5**: 客户端状态管理

### 开发规范

遵循 [项目规约](../../GEMINI.md) 中定义的开发规则：

- **TDD**: 测试驱动开发
- **Make-First**: 优先使用Makefile命令
- **类型安全**: 严格的TypeScript检查
- **组件复用**: 使用 `src/components/ui` 标准组件
- **双语支持**: 使用 `t('key')` 国际化

## 📁 目录结构

```
frontend/
├── src/
│   ├── app/              # Next.js 16 App Router 页面
│   ├── components/       # React 组件
│   │   └── ui/          # shadcn/ui 组件库
│   ├── actions/         # React 19 Server Actions
│   ├── hooks/           # 自定义 Hooks（含 AI Hooks）
│   ├── lib/             # 工具库
│   │   └── ai/         # AI 客户端（SSE 流式处理）
│   ├── stores/          # Zustand 状态管理
│   └── types/           # TypeScript 类型定义
├── tests/               # 测试文件
│   └── e2e/            # Playwright E2E 测试
└── public/              # 静态资源
```

## 🔗 相关资源

- [Frontend README](../../frontend/README.md) - 前端项目说明
- [测试脚本](../../scripts/frontend/run-tests.sh) - E2E测试执行脚本
- [项目规约](../../GEMINI.md) - 开发规范和最佳实践
