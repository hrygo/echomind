# EchoMind

**Your external brain for email.**

EchoMind is a SaaS-level Intelligent Email Decision System that acts as an AI cognitive layer on top of traditional email services. It reads, understands, classifies, and summarizes emails, transforming unstructured communication into structured insights and tasks for executives, managers, and dealmakers.

## 🚀 Features (Current Beta)

*   **Multi-Model AI Support**: Switch between DeepSeek, OpenAI, and Gemini via configuration.
*   **Intelligent Summary**: Generates concise 3-sentence summaries of emails.
*   **Sentiment Radar**: Analyzes email sentiment (Positive/Negative/Neutral) and urgency (High/Medium/Low).
*   **Contact Intelligence**: Automatically builds a contact database with interaction stats.
*   **Async Processing**: Robust background job processing with Redis & Asynq.
*   **Modern UI**: Next.js 14 Dashboard with real-time insights.

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

### 3. API Key Configuration

**Recommended Method (Local Development):**

1.  **Create your local config file**: Copy the example configuration and rename it:
    ```bash
    cp backend/configs/config.example.yaml backend/configs/config.yaml
    ```
2.  **Edit `backend/configs/config.yaml`**: Open this newly created file and replace `YOUR_DEEPSEEK_KEY` (or others) with your actual API keys.

    *Note: `backend/configs/config.yaml` is ignored by Git.*

**Environment Variables (CI/CD):**
Override any config using `ECHOMIND_` prefix (e.g., `ECHOMIND_AI_DEEPSEEK_API_KEY`).

### 4. Run the App

**Step 1: Start Infrastructure (DB & Redis)**
```bash
make docker-up
```

**Step 2: Run Backend Server**
```bash
make run-backend
```

**Step 3: Run AI Worker** (In a new terminal)
```bash
make run-worker
```

**Step 4: Run Frontend** (In a new terminal)
```bash
make run-frontend
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
