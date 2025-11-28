# 👻 幽灵故障复盘报告：macOS 大小写敏感性陷阱

## 1. 事故摘要
在 `feat/ai-native-refactor` 分支开发期间，程序运行正常。但在合并至 `main` 分支并执行 Reset 操作后，构建系统突然报错 `Module not found`，提示找不到 UI 组件。即使切回原分支，错误依然存在。

**根本原因**：macOS 文件系统（APFS）的**大小写不敏感**特性，掩盖了代码引用（`import`）与实际文件名（File System）不一致的问题。Git 操作和缓存失效撕开了这层“保护网”。

---

## 2. 事故时间脉络 (Timeline)

### 阶段一：潜伏期 (The Masking)
*   **状态**：开发环境 (macOS)。
*   **操作**：
    1.  创建文件 `Button.tsx` (PascalCase)。
    2.  代码中引用 `import ... from './button'` (kebab-case)。
    3.  运行 `npm run dev`。
*   **现象**：**运行成功**。
*   **原因**：macOS 忽略大小写，Next.js 询问系统“有 `button` 吗？”，系统回答“有（其实是 `Button`），给你”。Next.js 将此结果**缓存**。

### 阶段二：触发期 (The Trigger)
*   **状态**：Git 操作。
*   **操作**：
    1.  `git checkout main`
    2.  `git reset --hard ...`
*   **后果**：
    1.  Git 根据记录，在硬盘上重新写入文件 `Button.tsx`。
    2.  **关键点**：分支切换和 Reset 导致 Next.js/Webpack 的**缓存失效**。

### 阶段三：爆发期 (The Crash)
*   **状态**：构建检查。
*   **操作**：`npm run build`。
*   **现象**：**构建失败** (`Module not found`)。
*   **原因**：
    *   没有了缓存，Webpack 重新扫描依赖。
    *   Webpack (或 CI 环境) 遵循严格模式。它发现代码请求 `button`，但硬盘上是 `Button`。
    *   虽然 macOS 依然说“我有”，但 Webpack 的解析器检测到了大小写不一致（或者 Git 认为文件未变但构建工具认为路径错误），抛出异常。

### 阶段四：迷惑期 (The Confusion)
*   **操作**：切回 `feat/ai-native-refactor`。
*   **现象**：**依然报错**。
*   **用户心理**：“明明刚才在这个分支是好的！”
*   **原因**：
    *   **缓存已死**：之前的“保护伞”没了。
    *   **Git 的惰性**：在 macOS 上，Git 认为 `Button.tsx` 和 `button.tsx` 是同一个文件。如果你在编辑器里改名，Git 可能根本没记录下这个重命名操作。硬盘上依然躺着 `Button.tsx`。

---

## 3. 技术原理

### 核心冲突点
1.  **文件系统 (OS)**：
    *   macOS/Windows: `Button` == `button` (不敏感)
    *   Linux (CI/CD): `Button` != `button` (敏感)
2.  **Git**：
    *   默认配置下，在 macOS 上往往忽略文件名大小写变化 (`core.ignorecase = true`)。
3.  **构建工具**：
    *   现代构建工具（Webpack/Vite）为了跨平台兼容性，通常会尝试模拟大小写敏感检查，或者直接依赖 OS 的返回结果。

---

## 4. 经验教训 (Lessons Learned)

### ✅ 1. 严格的文件命名规范 (The Golden Rule)
*   **强制执行**：所有文件名统一使用 **kebab-case** (全小写，短横线分隔)。
    *   ❌ `UserProfile.tsx`
    *   ✅ `user-profile.tsx`
*   **原因**：全小写可以彻底规避跨平台的大小写敏感性问题。

### ✅ 2. 正确的重命名姿势
在 macOS/Windows 上重命名文件大小写时，**千万不要**直接在 IDE 或 Finder 中修改。
*   **错误做法**：右键重命名 `Button.tsx` -> `button.tsx` (Git 可能检测不到)。
*   **正确做法**：使用 `git mv` 命令。
    ```bash
    git mv Button.tsx button.tsx
    ```
    这会告诉 Git：“这是一个移动/重命名操作”，强制更新索引。

### ✅ 3. 信任 CI，不信任 Local
*   本地开发环境（Local）往往是“宽容”的。
*   CI 环境（Linux）是“严苛”的。
*   **结论**：如果本地能跑但 CI 挂了，99% 是环境差异（大小写问题是头号嫌疑人）。

### ✅ 4. 清理缓存
*   当遇到“灵异现象”（切分支不好使、改代码不生效）时，尝试删除 `.next` 或 `node_modules/.cache` 目录。缓存有时候会欺骗你。

---

**本次修复总结**：
我们通过 `git mv` 将所有 UI 组件强制重命名为小写（`button.tsx` 等），并更新了所有引用路径。这不仅修复了当前的构建错误，也为未来在 Linux 服务器上部署扫清了障碍。
