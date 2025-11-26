<div align="center">
  <h1 align="center">EchoMind</h1>
  <p align="center">
    <strong>🧠 Your Personal Neural Interface to Navigate the Information Maze 🧠</strong>
  </p>
  <p align="center">
    EchoMind is a smart, context-aware assistant that integrates with your digital life, starting with your emails, to create a searchable and intelligent knowledge base. It helps you stay organized, find information instantly, and gain insights from your communications.
  </p>
  <p align="center">
    <a href="README.zh.md">简体中文</a>
  </p>
  <p align="center">
    <img src="https://img.shields.io/github/actions/workflow/status/hrygo/echomind/ci-cd.yml?branch=main&style=for-the-badge" alt="CI/CD Status">
    <img src="https://img.shields.io/github/v/release/hrygo/echomind?style=for-the-badge" alt="Version">
    <img src="https://img.shields.io/github/license/hrygo/echomind?style=for-the-badge" alt="License">
  </p>
</div>

---



## ✨ Key Features

- **📧 Intelligent Email Sync**: Automatically syncs and processes emails from your IMAP accounts.
- **🧠 Contextual Understanding**: Builds a rich context graph from your communications to surface relevant information.
- **🔍 Advanced Search**: Perform semantic searches across all your synced data. Find not just keywords, but concepts and conversations.
- **🤖 AI-Powered Drafts**: Generate email replies and other text with the help of AI, based on the current context.
- **📈 Insight Generation**: (Coming Soon) Proactively provides summaries and insights from your data.

---

## 🔧 Tech Stack

| Category      | Technology                               |
|---------------|------------------------------------------|
| **Backend**   | Go, Gin, GORM, Asynq                     |
| **Frontend**  | Next.js, TypeScript, Tailwind CSS, Zustand |
| **Database**  | PostgreSQL (with pgvector), Redis        |
| **Container** | Docker                                   |
| **AI**        | OpenAI, Gemini                           |

---

## 🚀 Getting Started

Follow these instructions to get EchoMind up and running on your local machine for development and testing purposes.

### Prerequisites

Make sure you have the following tools installed:
- [Go](https://golang.org/doc/install) (version 1.22+)
- [Node.js](https://nodejs.org/en/download/) (version 18+) with [pnpm](https://pnpm.io/installation)
- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/)
- [make](https://www.gnu.org/software/make/)

### Installation & Setup

1.  **Clone the Repository**
    ```bash
    git clone https://github.com/your-username/echomind.git
    cd echomind
    ```
    *(Note: Remember to replace `your-username` with the actual repository owner's username.)*

2.  **Configure Environment Variables**
    Copy the example configuration files and update them with your credentials (e.g., OpenAI API key, database passwords).
    ```bash
    cp backend/configs/config.example.yaml backend/configs/config.local.yaml
    cp backend/configs/logger.example.yaml backend/configs/logger.local.yaml
    ```
    - Edit `backend/configs/config.local.yaml` to fill in the required secrets.

3.  **Start Backend Services**
    This command starts the required databases (Postgres, Redis) in Docker containers.
    ```bash
    make dev-db
    ```
    Then, run the database migrations:
    ```bash
    make db-init
    ```
    Finally, start the backend server:
    ```bash
    make run-be
    ```
    The backend API will be available at `http://localhost:8080`.

4.  **Start Frontend Application**
    In a new terminal, navigate to the `frontend` directory, install dependencies, and start the development server.
    ```bash
    cd frontend
    pnpm install
    pnpm dev
    ```
    The frontend application will be accessible at `http://localhost:3000`.

---

## 🧪 Running Tests

- **Backend Tests**:
  ```bash
  make test
  ```
- **Frontend Tests**:
  ```bash
  cd frontend
  pnpm test
  ```

---

## 🚢 Deployment

A production-ready setup can be deployed using Docker Compose:
```bash
docker-compose -f deploy/docker-compose.prod.yml up -d
```
This will build and run the frontend and backend containers, along with the required database services.

---

## 🗺️ Roadmap & Documentation

Our development is organized into clear, feature-driven phases. Here are our most recently completed and ongoing phases:

- ✅ **v1.1.0 (Enterprise Foundation)**: 规约化 Vector Architecture, Enterprise-Grade Logging Framework, Multi-process Configuration Management
- ✅ **v0.9.8 (Dashboard Integration Phase 2)**: Complete theme system, opportunity management, enhanced Dashboard API integration
- ✅ **v0.9.6-7 (Dashboard Integration Phase 1)**: SmartFeed AI functionality, task management, basic dashboard components
- ✅ **v0.9.2-4 (Neural Nexus)**: Context Bridge, Omni-Bar, Generative Widget Framework
- 🚧 **v1.2.0+ (Next Generation)**: Advanced AI features, enhanced productivity tools, ecosystem integrations

For detailed architecture and product specifications, please refer to our key documents:

- **[🏗️ Technical Architecture](docs/architecture.md)** - System design, vector search, AI services, and performance optimization
- **[📚 API Documentation](docs/api.md)** - Complete REST API reference with OpenAPI 3.0 specification
- **[📄 Unified Product & Technical Architecture](docs/product-design.md)** - Product design and user experience
- **[🗺️ Product Roadmap](docs/product-roadmap.md)** - Development roadmap and milestones
- **[📋 Product Requirements](docs/prd.md)** - Product requirements document
- **[🔧 Backend Configuration](backend/configs/README.md)** - Development setup and configuration guide

---

## 🤝 Contributing

We welcome contributions! Please read our [CONTRIBUTING.md](CONTRIBUTING.md) to learn about our development process, how to propose bugfixes and improvements, and how to build and test your changes.

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
