# 🏗️ Technical Design - Phase 5.2: RAG & Semantic Search

> **Target Version**: v0.6.0
> **Focus**: Vector Database, Embeddings, Semantic Search API.

## 1. Background
EchoMind v0.5.3 introduced the AI Command Center. To make the "Smart Feed" and "Intent Radar" truly powerful, we need to move beyond simple keyword matching and database queries. We need **Semantic Understanding** of the entire mailbox history.

## 2. Architecture Changes

### 2.1 Vector Database (pgvector)
We will enable `pgvector` extension on our existing PostgreSQL instance to store embeddings. This avoids managing a separate infrastructure like Chroma/Milvus for now (Keep it Simple).

**Schema Extension:**
```sql
CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE email_embeddings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email_id UUID NOT NULL REFERENCES emails(id) ON DELETE CASCADE,
    chunk_index INT NOT NULL DEFAULT 0,
    content_chunk TEXT NOT NULL,
    embedding vector(1536), -- For OpenAI text-embedding-3-small
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX ON email_embeddings USING hnsw (embedding vector_l2_ops);
```

### 2.2 AI Pipeline Update
The `analyze_task` worker will be updated:
1.  **Extract**: Clean text from HTML body.
2.  **Chunk**: Split long emails into manageable chunks (e.g., 512 tokens).
3.  **Embed**: Call AI Provider (DeepSeek/OpenAI) to get vectors.
4.  **Store**: Save to `email_embeddings`.

### 2.3 Search API (`/api/v1/search`)
New endpoint that accepts a natural language query:
*   **Query**: "Show me budget approvals from last week"
*   **Process**:
    1.  Embed the query.
    2.  Perform vector similarity search in Postgres.
    3.  Filter by metadata (Date range, Sender).
    4.  Rerank (Optional for future).
*   **Response**: List of relevant emails with "Reasoning" (why it matched).

## 3. Implementation Plan

### Step 1: Infrastructure
*   [ ] Update `docker-compose.yml` to use a Postgres image with `pgvector`.
*   [ ] Add `pgvector` migration in Go.

### Step 2: Backend Service
*   [ ] Implement `EmbeddingService` (interface for AI providers).
*   [ ] Update `EmailService` to trigger embedding generation.
*   [ ] Create `SearchService` for vector search logic.

### Step 3: Frontend Integration
*   [ ] Update Top Bar Search to use the new `/search` API.
*   [ ] Display search results with AI relevance snippets.

## 4. Risks & Mitigation
*   **Cost**: Embedding every email can be costly.
    *   *Mitigation*: Only embed "Important" emails first, or use a cheaper model (text-embedding-3-small).
*   **Performance**: Hybrid search (Keyword + Vector) is hard to tune.
    *   *Mitigation*: Start with pure Vector search for "Concept" queries, fallback to SQL `ILIKE` for exact matches.

## 5. Requirements for v0.7.0 Interoperability (Frontend "Neural Interface" Pre-requisites)

To enable the radical UI/UX refactor planned for v0.7.0 ("Neural Interface"), the v0.6.0 RAG backend **must** provide the following capabilities:

1.  **Citation Metadata in Search Results**:
    *   **Requirement**: Search results (from `/api/v1/search`) must return not just the relevant email, but also precise metadata indicating *which part* of the email was used to generate the answer. This includes:
        *   `email_id`: The ID of the source email.
        *   `chunk_id`: (Optional, if chunking is used) The ID of the specific chunk.
        *   `text_offset_start`, `text_offset_end`: Byte/character offsets within the *original email content* (or within the `content_chunk` if chunk_id is provided) to highlight the exact referenced text.
    *   **Purpose**: Allows the v0.7.0 frontend to render accurate citations (`[1]`, `[2]`) and, upon clicking, scroll/highlight the exact text in the "Grounding Panel" (right sidebar).

2.  **Structured Thread Data**:
    *   **Requirement**: A dedicated API (or an extension to existing `Email` API) to retrieve a full, structured conversation thread given an `email_id` or `thread_id`. The response should clearly indicate the chronological order and parent-child relationships within the thread.
    *   **Purpose**: Populates the "Grounding Panel" with the full context of a conversation when a user deep-dives into a citation.

3.  **Streaming JSON Responses**:
    *   **Requirement**: The `/api/v1/search` endpoint (or a new `/api/v1/chat` endpoint for RAG Q&A) should support **streaming responses**. This means the backend sends data chunks (e.g., individual words or sentences, or even JSON objects representing UI widgets) as they are generated, rather than waiting for the complete response.
    *   **Purpose**: Enables a real-time, "typing effect" user experience in the v0.7.0 "Studio Mode" chat interface and allows for dynamic rendering of UI components (widgets) as AI determines intent.

4.  **Dynamic Vector Store Filters**:
    *   **Requirement**: The search API must support advanced filtering on the vector store queries. This includes filtering by:
        *   `sender_email`
        *   `receiver_email`
        *   `date_range`
        *   `thread_id`
        *   Custom tags/labels (if implemented in v0.6.x for email classification).
    *   **Purpose**: Powers the "Context Manager" (left sidebar) in v0.7.0, allowing users to define specific "Attention Scopes" and restrict AI answers to those contexts.
