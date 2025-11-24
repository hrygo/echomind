# Toast 通知系统 - UX 改进方案

## 问题描述

之前使用原生 `alert()` 和 `confirm()` 对话框，存在以下问题：
1. **阻断式交互**：必须点击确定才能继续操作
2. **样式无法自定义**：浏览器原生样式，无法与应用风格统一
3. **用户体验差**：弹窗会打断用户流程

## 解决方案

### 1. 创建通用 Toast 通知系统

**文件**：`/frontend/src/lib/hooks/useToast.tsx`

**功能**：
- 支持 4 种类型：success, error, warning, info
- 自动消失（默认 4 秒）
- 支持自定义标题和时长
- 使用 Zustand 管理状态

**使用示例**：
```typescript
const toast = useToast();

// 成功通知
toast.success('同步任务已启动');

// 错误通知
toast.error('同步失败，请重试');

// 带标题的通知
toast.info('系统消息', '邮箱已配置成功', 5000);
```

### 2. 创建确认对话框组件

**文件**：`/frontend/src/components/ui/ConfirmDialog.tsx`

**功能**：
- 非阻断式确认对话框
- 可自定义标题、按钮文本
- 支持确认和取消回调

**使用示例**：
```typescript
const confirm = useConfirm();

confirm(
  '您确定要断开邮箱连接吗？',
  async () => {
    // 用户点击确认后执行
    await disconnectMailbox();
  },
  {
    title: '断开邮箱连接',
    confirmText: '确认',
    cancelText: '取消'
  }
);
```

### 3. Toast 容器组件

**文件**：`/frontend/src/components/ui/ToastContainer.tsx`

**功能**：
- 渲染所有 Toast 通知
- 固定在右上角
- 支持多个通知堆叠显示
- 带关闭按钮

### 4. 集成到根布局

**修改文件**：`/frontend/src/app/layout.tsx`

添加了：
```tsx
<ToastContainer />
<ConfirmDialog />
```

## 修改的文件

### 1. ConnectionTab.tsx

**替换内容**：
- ❌ `alert(t('settings.connection.syncStarted'))`
- ✅ `toast.success(t('settings.connection.syncStarted'))`

- ❌ `alert(error.message)`
- ✅ `toast.error(error.message)`

- ❌ `if (confirm('确定要断开吗？')) { ... }`
- ✅ `confirm('确定要断开吗？', () => { ... })`

### 2. inbox/page.tsx

**替换内容**：
- ❌ `alert('同步任务已启动')`
- ✅ `toast.success(t('inbox.syncStarted'))`

- ❌ `if (confirm('是否前往设置？')) { router.push(...) }`
- ✅ `confirm('是否前往设置？', () => { router.push(...) })`

## i18n 翻译添加

### 中文 (zh.json)

```json
{
  "common": {
    "confirm": "确认",
    "cancel": "取消",
    "goToSettings": "前往设置"
  },
  "settings": {
    "connection": {
      "syncNow": "立即同步",
      "syncing": "同步中...",
      "disconnectTitle": "断开邮箱连接",
      "disconnectSuccess": "邮箱连接已断开",
      "mailboxConfigDesc": "配置您的邮箱以开始同步"
    }
  },
  "inbox": {
    "syncStarted": "同步任务已启动",
    "syncFailed": "同步失败",
    "syncFailedUnknown": "同步失败：发生未知错误",
    "noAccountConfigured": "您尚未配置邮箱账户。是否立即前往设置页面进行配置？",
    "configureAccount": "配置邮箱账户"
  }
}
```

### 英文 (en.json)

对应的英文翻译已添加。

## 视觉效果

### Toast 通知
- 右上角滑入动画
- 带图标（✓ 成功、✗ 错误、ℹ 信息、⚠ 警告）
- 4 秒后自动消失
- 可手动关闭
- 多个通知垂直堆叠

### 确认对话框
- 居中模态对话框
- 带警告图标
- 两个按钮：取消（outline）和确认（primary）
- 点击外部或 ESC 关闭

## 技术细节

### 动画
使用 Tailwind 的 `animate-in slide-in-from-right` 实现滑入效果

### 状态管理
- Toast：Zustand store
- ConfirmDialog：Zustand store

### SSR 兼容
- ToastContainer 使用 `useState` 和 `useEffect` 确保只在客户端渲染
- 避免 hydration 错误

## 测试建议

1. **Toast 通知测试**：
   - 点击"立即同步"按钮
   - 应该看到绿色成功 Toast，4 秒后消失
   - 测试错误场景（未配置邮箱），应显示红色错误 Toast

2. **确认对话框测试**：
   - 点击"断开连接"按钮
   - 应弹出确认对话框（非原生）
   - 点击取消，对话框关闭，无操作
   - 点击确认，执行断开操作并显示成功 Toast

3. **多个通知测试**：
   - 快速触发多个操作
   - 应看到多个 Toast 垂直堆叠
   - 每个 Toast 独立计时消失

## 优势对比

| 特性 | 原生 alert/confirm | 新 Toast 系统 |
|------|-------------------|---------------|
| 阻断性 | ✗ 阻断所有操作 | ✓ 非阻断 |
| 自定义样式 | ✗ 无法自定义 | ✓ 完全自定义 |
| 自动消失 | ✗ 必须手动关闭 | ✓ 自动消失 |
| 多条通知 | ✗ 只能一条 | ✓ 支持堆叠 |
| 动画效果 | ✗ 无动画 | ✓ 滑入动画 |
| 品牌一致性 | ✗ 浏览器样式 | ✓ 应用风格统一 |

## 后续优化建议

1. **添加音效**：成功/错误时播放提示音
2. **撤销功能**：某些操作支持 undo（如删除邮件）
3. **位置配置**：支持左上、左下、右下等位置
4. **持久化 Toast**：重要通知可设置不自动消失
5. **进度 Toast**：长时间操作显示进度条
