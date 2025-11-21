# EchoMind

**Your external brain for email.**

EchoMind is a SaaS-level Intelligent Email Decision System that acts as an AI cognitive layer on top of traditional email services. It reads, understands, classifies, and summarizes emails, transforming unstructured communication into structured insights and tasks for executives, managers, and dealmakers.

## 🚀 Features (Current Beta v0.5.0)

*   **Multi-Model AI Support**: Switch between DeepSeek, OpenAI, and Gemini via configuration.
*   **Intelligent Summary**: Generates concise 3-sentence summaries of emails.
*   **Smart Reply**: AI-generated draft replies based on context and user intent.
*   **Relationship Intelligence**: Visual network graph of your email connections and interaction frequency.
*   **Sentiment Radar**: Analyzes email sentiment (Positive/Negative/Neutral) and urgency (High/Medium/Low).
*   **Contact Intelligence**: Automatically builds a contact database with interaction stats.
*   **Async Processing**: Robust background job processing with Redis & Asynq.
*   **Spam Filtering**: Rule-based filtering to skip AI analysis for spam emails, optimizing costs.
*   **Modern UI**: Next.js 16 Dashboard with real-time insights.

## 📚 Documentation

*   **[Product Requirements (PRD)](docs/prd.md)**: Detailed product vision, user personas, and functional requirements.
*   **[System Design](docs/product-design.md)**: High-level system architecture, module breakdown, and UI design.
*   **[Technical Architecture](docs/tech-architecture.md)**: Database schema, API design, and AI pipeline implementation details.
*   **[Roadmap](docs/product-roadmap.md)**: Development phases and future plans.
*   **[Contributing Guide](CONTRIBUTING.md)**: Development workflow, commit standards, and release process.

## 🛠️ Tech Stack

*   **Backend**: Go (Gin, GORM, Viper, Asynq)
*   **Frontend**: Next.js 16 (React, Tailwind CSS)
*   **Database**: PostgreSQL
*   **Cache/Queue**: Redis
*   **AI Integration**: Official SDKs for OpenAI, Gemini, and DeepSeek (via Adapter).

## 🏁 Getting Started

### 1. Prerequisites
*   Go (1.22+)
*   Node.js (20+)
*   pnpm (install with `npm install -g pnpm`)
*   Docker & Docker Compose (for PostgreSQL and Redis)

### 2. Setup

Clone the repository:

```bash
git clone https://github.com/your-username/echomind.git
cd echomind
make init
```

### 3. Configuration

**Local Development (Secure):**

1.  **Create Config File**: Copy the example configuration. This file is git-ignored to prevent secret leaks.
    ```bash
    cp backend/configs/config.example.yaml backend/configs/config.yaml
    ```
2.  **Update Credentials**: Edit `backend/configs/config.yaml` and replace the placeholders with your actual values:
    *   `server.jwt.secret`: Set a secure random string.
    *   `database`: Update `user` and `password` if you changed them in `docker-compose.yml`.
    *   `ai`: Add your API keys (e.g., `YOUR_DEEPSEEK_KEY`).

**Production / CI/CD:**

*   **Environment Variables**: The application supports overriding any config via environment variables with the `ECHOMIND_` prefix (e.g., `ECHOMIND_SERVER_JWT_SECRET`, `ECHOMIND_AI_DEEPSEEK_API_KEY`).
*   **Docker Compose**: The `deploy/docker-compose.prod.yml` is configured to use environment variables for sensitive data.

### 4. Run the App
**One-Click Start:**
```bash
make dev
```
This will start the database, redis, backend server, worker, and frontend in a single terminal. Press `Ctrl+C` to stop.

**Stop All Services:**
```bash
make stop
```

Open `http://localhost:3000` to view the dashboard.

## 🧪 Testing

Run backend unit and integration tests:
```bash
make test
```

Run frontend tests:
```bash
cd frontend && pnpm test
```

## 📂 Project Structure

```text
/echomind
├── backend/           # Go API & Worker
│   ├── cmd/           # Entry points
│   ├── configs/       # Configuration files
│   ├── internal/      # Private application logic
│   └── pkg/           # Public libraries (AI providers, IMAP)
├── frontend/          # Next.js Web App
├── deploy/            # Docker Compose & DevOps
├── docs/              # Project Documentation
└── scripts/           # Utilities
```
