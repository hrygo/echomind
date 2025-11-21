# 📅 Weekly Sprint Plan: RAG & Semantic Search (v0.6.0)

> **Goal**: Implement End-to-End Semantic Search for Emails.
> **Timeline**: 5 Days (Mon-Fri).

## 🗓️ Schedule

### Day 1: Infrastructure & Schema (Foundation)
**Objective**: Enable Vector Storage capability.

*   **Morning**: Docker & Environment
    *   [ ] Update `docker-compose.yml` to use `pgvector/pgvector:pg16`.
    *   [ ] Verify extension activation (`CREATE EXTENSION vector;`).
*   **Afternoon**: GORM & Migration
    *   [ ] Install `github.com/pgvector/pgvector-go`.
    *   [ ] Create `internal/model/embedding.go`.
    *   [ ] Implement AutoMigrate logic in `main.go`.
    *   [ ] Create `HNSW` index via raw SQL migration/init script.

### Day 2: Embedding Service (The "Eyes")
**Objective**: Connect to AI models to generate vectors.

*   **Morning**: Provider Interface
    *   [ ] Define `EmbeddingProvider` interface in `pkg/ai`.
    *   [ ] Implement `OpenAI` provider (using `text-embedding-3-small`).
    *   [ ] (Optional) Implement `DeepSeek` provider if compatible/available.
*   **Afternoon**: Text Processing
    *   [ ] Create `pkg/utils/chunker.go`.
    *   [ ] Implement HTML-to-Text stripping.
    *   [ ] Implement "Sliding Window" or "Paragraph-based" chunking strategy.

### Day 3: Ingestion Pipeline (The "Memory")
**Objective**: Automate the flow from Email -> Vector DB.

*   **Morning**: Worker Integration
    *   [ ] Update `internal/tasks/analyze.go`.
    *   [ ] Workflow: `Parse Email` -> `Summary` -> `Chunk` -> `Embed` -> `Save`.
*   **Afternoon**: Backfill & Testing
    *   [ ] Create a CLI command `make reindex` to process existing emails.
    *   [ ] Verify data in `email_embeddings` table.

### Day 4: Search Service (The "Brain")
**Objective**: Query the vector database.

*   **Morning**: Search Logic
    *   [ ] Create `internal/service/search.go`.
    *   [ ] Implement Cosine Distance search using GORM/SQL.
    *   [ ] Join results with `emails` table to get metadata (Subject, Sender).
*   **Afternoon**: API Exposure
    *   [ ] Create `internal/handler/search.go`.
    *   [ ] Register `GET /api/v1/search` endpoint.
    *   [ ] Define Request/Response DTOs (including `relevance_score`).

### Day 5: UI Integration & Polish (The Experience)
**Objective**: Deliver the feature to the user.

*   **Morning**: Frontend
    *   [ ] Update `Header.tsx` Search Bar to handle `Enter` key.
    *   [ ] Create `SearchResults` component (Dropdown or Dedicated Page).
    *   [ ] Show "Semantic Match" highlights (optional).
*   **Afternoon**: QA & Release
    *   [ ] Test queries: "Budget from last week", "Meeting requests from Alice".
    *   [ ] Performance tuning (adjust chunk size or limit).
    *   [ ] Tag v0.6.0-beta.
