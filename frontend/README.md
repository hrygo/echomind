# EchoMind Frontend

**Version**: v1.1.0  
**Tech Stack**: Next.js 16.0.3 + React 19.2.0 + TypeScript + Tailwind CSS v4

基于 Next.js 16 和 React 19 的现代化 AI Native 前端应用。

## 🚀 快速开始

### 开发环境

```bash
# 安装依赖
pnpm install

# 启动开发服务器
pnpm dev

# 访问应用
open http://localhost:3000
```

### 构建生产版本

```bash
# 类型检查
pnpm type-check

# 构建
pnpm build

# 启动生产服务器
pnpm start
```

### 测试

```bash
# 运行单元测试
pnpm test

# 运行端到端测试
../scripts/frontend/run-tests.sh
```

## 📁 项目结构

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

## 🎯 核心特性

### AI Native 架构

- **Server Actions**: React 19 原生服务端操作
- **SSE 流式处理**: 实时 AI 响应流
- **自定义 AI Hooks**: `useStreamChat`、`useAIDraft`、`useAIReply`
- **类型安全**: 完整的 TypeScript 类型定义

### 现代化技术栈

- **Next.js 16**: App Router + 异步 API
- **React 19**: Server Components + Actions
- **shadcn/ui**: 基于 Radix UI 的组件库
- **Tailwind CSS v4**: 原生 CSS + PostCSS
- **TanStack Query v5**: 服务端状态管理
- **Zustand v5**: 客户端状态管理

## 📚 文档

- [AI Native 架构设计](../docs/architecture/AI_NATIVE_ARCHITECTURE.md)
- [端到端测试指南](../docs/frontend/testing/E2E_TEST_GUIDE.md)
- [测试检查清单](../docs/frontend/testing/TEST_CHECKLIST.md)
- [测试验证摘要](../docs/frontend/testing/TEST_VERIFICATION_SUMMARY.md)

## 🛠️ 开发规范

遵循 [项目规约](../GEMINI.md) 中的开发规则：

- **TDD**: 测试驱动开发
- **Make-First**: 优先使用 Makefile 命令
- **类型安全**: 严格的 TypeScript 检查
- **组件复用**: 使用 `src/components/ui` 标准组件
- **双语支持**: 使用 `t('key')` 国际化

## 📦 依赖管理

使用 pnpm 进行依赖管理：

```bash
# 添加依赖
pnpm add <package>

# 添加开发依赖
pnpm add -D <package>

# 更新依赖
pnpm update
```

## 🔗 相关链接

- [Next.js 文档](https://nextjs.org/docs)
- [React 19 文档](https://react.dev)
- [shadcn/ui 文档](https://ui.shadcn.com)
- [Tailwind CSS 文档](https://tailwindcss.com/docs)
