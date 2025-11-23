# 📚 EchoMind UX/UE 重构与开发规格说明书 (v2.1)

> **版本**: v2.1 | **状态**: 待开发 | **目标**: Zero-Friction Onboarding (零摩擦入职)

这份文档旨在为开发团队提供详尽的实施指南，确保初级工程师也能准确实现设计意图。

---

## 1. 基础设施与工具 (Infrastructure)

### 1.1 邮箱服务商配置 (`lib/constants/mail_providers.ts`)

我们需要一个静态配置来驱动智能邮箱连接功能。

```typescript
export interface MailProviderConfig {
  id: string;
  domains: string[]; // 匹配的域名后缀，如 ["gmail.com", "googlemail.com"]
  name: string;      // 显示名称
  imap: { host: string; port: number; secure: boolean };
  smtp: { host: string; port: number; secure: boolean };
  helpLink?: string; // 获取 App Password 的帮助链接
  requiresAppPassword?: boolean; // 是否强制提示使用应用密码
}

export const MAIL_PROVIDERS: MailProviderConfig[] = [
  {
    id: 'gmail',
    domains: ['gmail.com'],
    name: 'Gmail',
    imap: { host: 'imap.gmail.com', port: 993, secure: true },
    smtp: { host: 'smtp.gmail.com', port: 465, secure: true },
    helpLink: 'https://support.google.com/accounts/answer/185833',
    requiresAppPassword: true
  },
  {
    id: 'outlook',
    domains: ['outlook.com', 'hotmail.com', 'live.com'],
    name: 'Outlook',
    imap: { host: 'outlook.office365.com', port: 993, secure: true },
    smtp: { host: 'smtp.office365.com', port: 587, secure: false }, // STARTTLS usually uses 587
    requiresAppPassword: false // Modern Auth usually supported, but for simple IMAP logic might depend
  },
  // Add QQ, 163, etc. later
];

// Helper function
export function detectProvider(email: string): MailProviderConfig | null {
  const domain = email.split('@')[1]?.toLowerCase();
  if (!domain) return null;
  return MAIL_PROVIDERS.find(p => p.domains.includes(domain)) || null;
}
```

### 1.2 国际化字典更新 (`locales/*.json`)

请参考 `docs/design-v2.0.md` 中的 JSON 结构，务必在开发前先将这些 Key 填入 `en.json` 和 `zh.json`。

---

## 2. 模块详细设计 (Module Specifications)

### 2.1 统一认证模块 (`/auth`)

**页面结构**: `app/auth/page.tsx`

*   **URL 参数**: `?mode=login` (默认) 或 `?mode=register`。
*   **组件树**:
    ```text
    AuthPage
    ├── AuthLayout (左右分栏容器)
    │   ├── LeftSide (品牌视觉，仅 Desktop)
    │   └── RightSide (表单容器)
    │       ├── AuthHeader (标题: Welcome back / Create account)
    │       ├── LoginForm
    │       ├── RegisterForm
    │       └── AuthSwitch (底部切换链接)
    ```

#### 组件详述: `LoginForm`
*   **State**:
    *   `email` (string)
    *   `password` (string)
    *   `isLoading` (boolean)
    *   `errors` (Record<string, string>)
*   **Props**: 无
*   **逻辑**:
    *   提交时调用 `useAuthStore.login(email, password)`。
    *   成功 -> 检查 `user.has_configured_account` (需后端支持字段) ? 跳转 `/dashboard` : 跳转 `/onboarding`。
    *   失败 -> 设置 `errors`，显示在 Input 下方。
*   **扩展点 (Future Phase 7)**:
    *   `SocialLogin`: 在 "Sign In" 按钮下方添加 "WeChat Login" (扫码图标)。
    *   逻辑: 点击弹出微信二维码模态框，轮询扫码状态。

#### 组件详述: `RegisterForm`
*   **State**: `email`, `password`, `name` (新增), `isLoading`, `errors`。
*   **逻辑**:
    *   提交调用 `useAuthStore.register(...)`。
    *   成功 -> 自动登录 -> 跳转 `/onboarding`。

### 2.2 新手引导模块 (`/onboarding`)

**页面结构**: `app/onboarding/page.tsx`

*   **状态管理**: 使用局部 Zustand Store `useOnboardingStore` 或 React Context，因为这些状态只在引导期间有用。

```typescript
// store/onboarding.ts
interface OnboardingState {
  step: 1 | 2 | 3 | 4; // Step 4 is optional WeChat Bind
  role: string | null; // 'executive' | 'manager' | 'dealmaker'
  mailbox: {
    email: string;
    password: string;
    providerConfig: MailProviderConfig | null; // 自动匹配的配置
    manualConfig?: { ... }; // 用户手动输入的配置（如果 providerConfig 为空）
  };
  setStep: (step: number) => void;
  setRole: (role: string) => void;
  setMailbox: (data: Partial<OnboardingState['mailbox']>) => void;
}
```

#### Step 1: `RoleSelector`
*   **UI**: 3 个卡片 (`div`)，Flex 布局。
*   **交互**: 点击卡片 -> 调用 `setRole` -> 选中样式 (Ring/Border) -> "Next" 按钮 `disabled={!role}` 解除。

#### Step 2: `SmartMailboxForm` (核心难点)
*   **UI 元素**:
    *   `EmailInput`: `onChange` 时调用 `detectProvider(value)`。
    *   `ConfigPreview`: 如果 `providerConfig` 存在，显示 "Detected {Provider Name}" 和绿色的锁图标。
    *   `ManualConfigToggle`: 文字链接 "Manual Configuration"。点击展开详细表单。
    *   `PasswordInput`: 如果是 Gmail，下方显示 `requiresAppPassword` 提示。
*   **逻辑**:
    *   **Effect**: 监听 `email` 变化 -> 更新 `providerConfig`。
    *   **Submit**:
        1.  构造 payload: 优先使用 `providerConfig`，否则使用手动填写的表单数据。
        2.  调用 `api.post('/settings/account/validate', payload)` (需后端实现，或暂时直接调保存接口)。
        3.  成功 -> `setStep(3)`。
        4.  失败 -> 显示 Toast 或 Inline Error。

#### Step 3: `InitialSync`
*   **UI**: 简单的 Lottie 动画或 CSS Spinner。
*   **逻辑**:
    *   `useEffect` (mount) -> 调用 `api.post('/sync')`。
    *   等待 3 秒 (为了让用户看清动画，建立心理预期)。
    *   `router.push('/dashboard')`。

#### Step 4: `ConnectWeChat` (Optional - Future Phase 7)
*   **触发**: 在 Step 3 成功后，或者作为 Dashboard 的引导卡片。
*   **UI**:
    *   左侧: "Get Instant Alerts on WeChat" 预览图。
    *   右侧: 绑定二维码。
*   **逻辑**: 扫码绑定成功后，开启 "Risk Alert" 和 "Daily Digest" 推送。

### 2.3 设置中心重构 (`/dashboard/settings`)

**页面结构**: `app/dashboard/settings/page.tsx` (Client Component)

*   **布局**: 使用 Radix UI `Tabs` 组件。
    *   `TabsList`: Profile, Connection, Notification, Preferences...
    *   `TabsContent`: `ProfileTab`, `ConnectionTab`, `NotificationTab`.

#### 组件详述: `ConnectionTab`
*   **数据源**: `useAuthStore.user` (需确保包含连接状态) 或 `useSettingsStore` (如果拆分)。
*   **状态显示**:
    *   如果是 `isConnected`: 显示 "Green Dot" + "Connected"。显示 "Last synced at: {time}"。
    *   如果是 `!isConnected`: 显示 "Red Dot" + "Disconnected"。显示 "Reconnect" 按钮。
*   **重连逻辑**:
    *   点击 "Reconnect" -> 弹出 Modal (复用 `SmartMailboxForm`，但 Email 字段只读)。

#### 组件详述: `NotificationTab` (Future Phase 7)
*   **WeChat Binding**:
    *   状态: "Not Connected" (显示绑定按钮) / "Connected as [Nickname]" (显示解绑按钮)。
    *   Toggle Switches: "Daily Digest via WeChat", "Urgent Alerts via WeChat".

---

## 3. 后端接口需求 (Backend Contract)

为了支持上述前端逻辑，后端需要提供或确认以下接口行为：

### 3.1 `POST /api/v1/auth/login` & `register`
*   **Response**:
    ```json
    {
      "token": "jwt...",
      "user": {
        "id": "...",
        "email": "...",
        "name": "...",
        "role": "manager", // 新增: 可能为空
        "has_account": false // 新增: 标识是否已绑定邮箱
      }
    }
    ```
    *注: 如果后端尚未支持 `role` 和 `has_account`，前端可以通过调用 `GET /api/v1/settings/account` 并在 404 时判断为 `false`，但这样会有额外的 RTT。建议后端在 Login 响应中带上。*

### 3.2 `POST /api/v1/settings/account` (Update)
*   **现有逻辑**: 保存并尝试连接。失败返回 400。
*   **需求**: 保持不变。前端通过 try-catch 处理 400 错误，并在 Step 2 显示。

### 3.3 `PATCH /api/v1/users/me` (New/Update)
*   **用途**: Onboarding Step 1 保存用户角色。
*   **Payload**: `{ "role": "executive" }`。

---

## 4. 开发任务清单 (Task List)

### Phase 0: 准备工作
1.  [ ] 创建 `frontend/src/lib/constants/mail_providers.ts`。
2.  [ ] 更新 `en.json` 和 `zh.json` (参考 v2.0 设计文档)。

### Phase 1: 认证页面 (`/auth`)
1.  [ ] 创建 `AuthLayout` 组件 (UI)。
2.  [ ] 创建 `LoginForm` 和 `RegisterForm` 组件 (逻辑 + UI)。
3.  [ ] 在 `app/auth/page.tsx` 中整合。
4.  [ ] 编写 `tests/e2e/auth-flow.spec.ts`。

### Phase 2: 引导流程 (`/onboarding`)
1.  [ ] 创建 `onboarding` store (Zustand)。
2.  [ ] 开发 `RoleSelector` 组件。
3.  [ ] 开发 `SmartMailboxForm` 组件 (集成 `mail_providers.ts`)。
4.  [ ] 开发 `app/onboarding/page.tsx` (步骤控制器)。
5.  [ ] 编写 `tests/e2e/onboarding-flow.spec.ts`。

### Phase 3: 路由守卫与设置
1.  [ ] 修改 `AuthGuard.tsx`: 检查 `user.has_account`，如果为 false 且不在 `/onboarding`，则跳转。
2.  [ ] 重构 `app/dashboard/settings/page.tsx` 使用 Tabs 布局。
