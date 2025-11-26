# 🛡️ EchoMind 规约

**愿景**: 个人智能神经中枢 | **版本**: v1.1.0 (Enterprise Release)
**技术栈**: Go(Gin/GORM/Asynq) + Next.js + Postgres(pgvector) + Redis

---

## 🚀 版本发布规约

### 版本检查清单
- `frontend/package.json`
- `backend/pkg/logger/config.go`
- `Makefile` (VERSION 变量)
- `docs/openapi.yaml`
- `backend/configs/logger*.yaml`
- `README*.md` (路线图)
- `CHANGELOG.md`
- `docs/product-roadmap.md`
- `docs/logger/README.md`

### 发布流程
```bash
git add .
git commit -m "feat: v{version} - description"
git tag -a v{version} -m "release notes"
```

### 版本策略
- **语义化**: `v{MAJOR}.{MINOR}.{PATCH}`
- **企业级**: v1.0+ 标志生产就绪
- **配置同步**: 所有版本引用文件保持一致

---

## ⚡ 核心开发规约

### 质量保证
- **提交前**: `make test` + `make build` + `make test-fe` + `make build-fe`
- **测试优先**: Mock 外部依赖 (AI, DB)
- **构建验证**: 确保编译无错误

### 架构原则
- **数据库**: GORM 模型变更后编译验证
- **前端**: 优先复用 `src/components/ui` 组件
- **重构**: `grep` 全局搜索，保留旧API过渡
- **国际化**: 强制双语 `t('key')`

### 工具使用
```bash
# 优先使用 Make 命令
make test        # > go test ./...
make build       # > go build ./cmd/main.go
make run-backend # > cd backend && go run cmd/main.go
make stop        # 清理所有进程
```

---

## 📋 AI 代理操作标准

### 工作环境
- **目录**: 必须为 `~/aicoding/echomind`
- **验证**: 命令执行前确认工作目录

### 开发操作
- **文件操作**: 最小化上下文，失败时优先状态检查
- **提交规范**: `feat:` `fix:` `docs:` `refactor:` 前缀
- **原子提交**: 频繁、小粒度提交
- **版本发布**: 按照清单逐项检查

### 指导原则
- **Make 优先**: 使用统一的 Make 命令接口
- **状态检查**: 操作失败时使用 `read_file` 验证
- **全局搜索**: 重构前使用 `grep` 查找所有引用
- **渐进式**: 保留旧API，逐步迁移