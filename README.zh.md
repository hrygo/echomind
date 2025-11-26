<div align="center">

<h1 align="center">EchoMind</h1>

<p align="center">
  <strong>🧠 您的个人智能神经中枢，轻松驾驭信息迷宫 🧠</strong>
</p>

<p align="center">
  EchoMind 是一款智能的、具备情境感知能力的个人助理。它能深度融入您的数字生活（从电子邮件开始），为您创建一个可搜索的、智能化的知识库。它能帮助您保持条理，即时找到所需信息，并从日常沟通中获得洞见。
</p>

<p align="center">
    <a href="README.md">English</a>
</p>

<p align="center">
    <img src="https://img.shields.io/github/actions/workflow/status/hrygo/echomind/ci-cd.yml?branch=main&style=for-the-badge" alt="CI/CD 状态">
    <img src="https://img.shields.io/github/v/release/hrygo/echomind?style=for-the-badge" alt="版本">
    <img src="https://img.shields.io/github/license/hrygo/echomind?style=for-the-badge" alt="许可证">
</p>
</div>

---



## ✨ 主要功能

- **📧 智能邮件同步**: 自动同步并处理来自您 IMAP 账户的邮件。
- **🧠 情境理解**: 从您的通信内容中构建丰富的上下文关系图，以呈现相关信息。
- **🔍 高级搜索**: 在您所有同步数据中执行语义搜索。不仅能找到关键词，还能发现概念和对话。
- **🤖 AI 驱动的草稿**: 基于当前上下文，借助 AI 生成邮件回复和其他文本。
- **📈 洞察生成**: (即将推出) 主动从您的数据中提供摘要和洞察。

---

## 🔧 技术栈

| 类别       | 技术栈                                   |
|------------|------------------------------------------|
| **后端**   | Go, Gin, GORM, Asynq                     |
| **前端**   | Next.js, TypeScript, Tailwind CSS, Zustand |
| **数据库** | PostgreSQL (with pgvector), Redis        |
| **容器化** | Docker                                   |
| **AI**     | OpenAI, Gemini                           |

---

## 🚀 快速开始

请遵循以下说明在您的本地计算机上启动并运行 EchoMind，以便进行开发和测试。

### 环境准备

请确保您已安装以下工具：
- [Go](https://golang.org/doc/install) (版本 1.22+)
- [Node.js](https://nodejs.org/en/download/) (版本 18+) 及 [pnpm](https://pnpm.io/installation)
- [Docker](https://docs.docker.com/get-docker/) 和 [Docker Compose](https://docs.docker.com/compose/install/)
- [make](https://www.gnu.org/software/make/)

### 安装与设置

1.  **克隆代码仓库**
    ```bash
    git clone https://github.com/your-username/echomind.git
    cd echomind
    ```
    *(注意: 请记得将 `your-username` 替换为实际的代码仓库所有者用户名。)*

2.  **配置环境变量**
    复制示例配置文件，并用您的凭据（例如 OpenAI API 密钥、数据库密码）更新它们。
    ```bash
    cp backend/configs/config.example.yaml backend/configs/config.yaml
    cp backend/configs/logger.example.yaml backend/configs/logger.yaml
    ```
    - 编辑 `backend/configs/config.yaml` 文件，填入所需的安全密钥和配置。

3.  **启动后端服务**
    此命令会在 Docker 容器中启动所需的数据库 (Postgres, Redis)。
    ```bash
    make dev-db
    ```
    然后，运行数据库迁移：
    ```bash
    make db-init
    ```
    最后，启动后端服务器：
    ```bash
    make run-be
    ```
    后端 API 将在 `http://localhost:8080` 上可用。

4.  **启动前端应用**
    在新的终端窗口中，进入 `frontend` 目录，安装依赖并启动开发服务器。
    ```bash
    cd frontend
    pnpm install
    pnpm dev
    ```
    前端应用将可以在 `http://localhost:3000` 访问。

---

## 🧪 运行测试

- **后端测试**:
  ```bash
  make test
  ```
- **前端测试**:
  ```bash
  cd frontend
  pnpm test
  ```

---

## 🔍 CI/CD 监控

EchoMind 包含一个强大的 CI/CD 监控工具，帮助您跟踪构建状态、分析失败原因，并获得可行的洞察。

完整的设置说明、使用示例和高级功能，请参阅 [scripts/CI_README.md](scripts/CI_README.md)。

**快速开始**：
```bash
# 设置日常使用别名
echo 'alias ci="./scripts/ci.sh"' >> ~/.zshrc && source ~/.zshrc

# 基本使用
ci                 # 当前状态
ci watch           # 监控运行
ci history         # 查看历史
ci analyze         # 深度分析
ci interactive     # 交互菜单
```

---

## 🚢 应用部署

可以使用 Docker Compose 部署一个生产就绪的环境：
```bash
docker-compose -f deploy/docker-compose.prod.yml up -d
```
这将构建并运行前端和后端容器，以及所需的数据库服务。

---

## 🗺️ 路线图与文档

我们的开发围绕清晰的、以功能为导向的阶段进行。以下是我们最近完成和正在进行的阶段：

- ✅ **v0.9.8 (Dashboard集成阶段二)**: 完整主题系统、商机管理、增强的Dashboard API集成
- ✅ **v0.9.6-7 (Dashboard集成阶段一)**: SmartFeed AI功能、任务管理、基础Dashboard组件
- ✅ **v0.9.2-4 (智能中枢 / Neural Nexus)**: 上下文桥梁、全能入口、生成式UI组件
- 🚧 **v0.9.9+ (微信操作系统 / WeChat OS)**: 语音指令、一键决策、日历守护、晨间简报

有关详细的架构和产品规格，请参阅我们的主要文档：

- **[📚 统一产品与技术架构](docs/product-design.md)**
- **[🗺️ 产品路线图](docs/product-roadmap.md)**
- **[🔄 EchoMind 邮件处理系统时序图](docs/api_search_sequence_diagram.md)** - 完整的系统流程时序图和架构说明

---

## 🤝 如何贡献

我们欢迎任何形式的贡献！请阅读我们的 [CONTRIBUTING.md](CONTRIBUTING.md) 文件，以了解我们的开发流程、如何提出错误修复和改进建议，以及如何构建和测试您的更改。

---

## 📄 开源许可

本项目基于 MIT 许可证授权 - 详情请参阅 [LICENSE](LICENSE) 文件。
