# 日志优化方案

## 问题描述

当前系统日志存在两个主要问题：
1. 日志中显示源码绝对文件路径，暴露系统信息
2. 部分日志可能包含原始邮件内容，存在隐私风险

## 解决方案

### 1. 优化日志配置

**文件**：`backend/pkg/logger/logger.go`

**修改内容**：
```go
// 使用短文件名，不显示完整路径
encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
```

**效果**：
- 之前：`/Users/huangzhonghui/aicoding/echomind/backend/internal/service/sync.go:191`
- 之后：`internal/service/sync.go:191`

### 2. 替换所有 `log.Printf` 为 zap logger

#### 2.1 Sync Service (`backend/internal/service/sync.go`)

**替换前**：
```go
log.Printf("Error fetching email account for user %s, team %v, org %v: %v", userID, teamID, organizationID, err)
log.Printf("Error dialing/logging into IMAP server %s for account %s: %v", addr, account.ID, err)
log.Printf("Error fetching emails for account %s: %v", account.ID, err)
log.Printf("Failed to create email for user %s: %v", userID, err)
log.Printf("DB error for user %s: %v", userID, result.Error)
log.Printf("Failed to create task for email %s (user %s): %v", email.ID, userID, err)
log.Printf("Failed to enqueue task for email %s (user %s): %v", email.ID, userID, err)
log.Printf("Enqueued analysis task for email %s (user %s)", email.ID, userID)
log.Printf("Failed to update contact for user %s, email %s: %v", userID, senderEmail, err)
```

**替换后**：
```go
s.logger.Errorw("Failed to fetch email account",
    "user_id", userID,
    "team_id", teamID,
    "org_id", organizationID,
    "error", err)
s.logger.Errorw("IMAP connection failed",
    "address", addr,
    "account_id", account.ID,
    "error", err)
// ... 其他类似替换
```

#### 2.2 Analyze Task (`backend/internal/tasks/analyze.go`)

**替换前**：
```go
log.Printf("[Task Started] Email Analysis - EmailID: %s, UserID: %s, TaskType: %s", p.EmailID, p.UserID, t.Type())
log.Printf("[Task Completed] Email Analysis - EmailID: %s, UserID: %s, Duration: %.2fs", p.EmailID, p.UserID, duration.Seconds())
log.Printf("Email %s identified as spam: %s", p.EmailID, spamReason)
log.Printf("[Email Analyzed] EmailID: %s, UserID: %s, Category: %s, Sentiment: %s, Urgency: %s", p.EmailID, p.UserID, email.Category, email.Sentiment, email.Urgency)
log.Printf("Warning: Failed to update contact stats for sender %s: %v", email.Sender, err)
log.Printf("Warning: Failed to assign contexts to email %s: %v", email.ID, err)
log.Printf("Warning: Failed to match contexts for email %s: %v", email.ID, err)
log.Printf("Warning: Failed to process embedding for email %s: %v", p.EmailID, err)
```

**替换后**：
```go
logger.Infow("[Task Started] Email Analysis",
    "email_id", p.EmailID,
    "user_id", p.UserID,
    "task_type", t.Type())
logger.Infow("[Task Completed] Email Analysis",
    "email_id", p.EmailID,
    "user_id", p.UserID,
    "duration_seconds", duration.Seconds())
// ... 其他类似替换
```

#### 2.3 Sync Task (`backend/internal/tasks/sync.go`)

**替换前**：
```go
log.Printf("Starting email sync for user: %s", p.UserID)
log.Printf("Email sync failed for user %s: %v", p.UserID, err)
log.Printf("Email sync completed for user: %s", p.UserID)
```

**替换后**：
```go
logger.Infow("Starting email sync",
    "user_id", p.UserID)
logger.Errorw("Email sync failed",
    "user_id", p.UserID,
    "error", err)
logger.Infow("Email sync completed",
    "user_id", p.UserID)
```

### 3. 安全性增强

#### 3.1 移除敏感信息
所有日志记录都避免包含以下敏感信息：
- 邮件正文内容
- 邮件附件内容
- 用户密码
- 个人身份信息

#### 3.2 使用结构化日志
使用 zap 的 `Infow`/`Errorw` 方法记录结构化日志：
```go
logger.Errorw("Failed to fetch email account",
    "user_id", userID,
    "team_id", teamID,
    "org_id", organizationID,
    "error", err)
```

### 4. 接口适配

为确保 SyncService 正确实现 EmailSyncer 接口，在 `backend/internal/service/sync.go` 中添加：
```go
// Ensure SyncService implements the EmailSyncer interface
var _ tasks.EmailSyncer = (*SyncService)(nil)
```

同时修改 `EmailSyncer` 接口以匹配 `SyncEmails` 方法签名：
```go
type EmailSyncer interface {
    SyncEmails(ctx context.Context, userID uuid.UUID, teamID *uuid.UUID, organizationID *uuid.UUID) error
}
```

## 验证结果

### 1. 编译检查
所有修改后的文件通过了编译检查，无语法错误。

### 2. 接口兼容性
SyncService 正确实现了 EmailSyncer 接口，worker 可以正常调用。

### 3. 日志格式改进
- 文件路径使用相对路径显示
- 避免记录邮件正文等敏感内容
- 使用结构化日志便于后续分析

## 技术细节

### 日志级别使用规范
- **Debug**：调试信息，开发环境使用
- **Info**：常规操作信息，如任务开始/完成
- **Warn**：警告信息，不影响主要流程但需要注意
- **Error**：错误信息，影响功能但可恢复
- **Fatal**：致命错误，程序无法继续运行

### 结构化日志字段命名
采用下划线命名法，保持一致性：
- `user_id`
- `email_id`
- `task_type`
- `duration_seconds`
- `error`

## 后续优化建议

1. **日志轮转配置**：
   - 当前已配置日志文件轮转（100MB，保留3个备份，7天过期）
   - 可根据实际需求调整轮转策略

2. **日志级别动态调整**：
   - 可通过配置文件或环境变量动态调整日志级别
   - 生产环境建议使用 Info 级别

3. **日志聚合**：
   - 可集成 ELK Stack 或类似方案进行日志聚合分析
   - 便于问题追踪和系统监控

4. **性能监控**：
   - 可在关键路径添加性能监控日志
   - 记录方法执行时间，便于性能优化

## 总结

本次优化成功解决了日志中的两个主要问题：
1. ✅ 移除了源码绝对文件路径，使用相对路径显示
2. ✅ 替换了所有 `log.Printf` 为结构化 zap logger
3. ✅ 确保不记录邮件正文等敏感信息
4. ✅ 保持了接口兼容性，系统可正常运行

优化后的日志更加安全、规范且易于维护。