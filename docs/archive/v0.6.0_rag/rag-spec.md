# 🛠️ RAG Implementation Specification (v0.6.0)

> **Status**: Draft
> **Tech Stack**: Go, PostgreSQL (pgvector), OpenAI API.

## 1. Database Schema

### 1.1 Extension
Must enable `vector` extension in Postgres.
```sql
CREATE EXTENSION IF NOT EXISTS vector;
```

### 1.2 Table: `email_embeddings`
Stores text chunks and their vector representations.

```go
// internal/model/embedding.go

type EmailEmbedding struct {
    ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    EmailID   uuid.UUID      `gorm:"type:uuid;not null;index"` // FK to emails
    ChunkIndex int           `gorm:"not null"`                 // Order of chunk in email
    Content   string         `gorm:"type:text;not null"`       // The actual text chunk
    Embedding pgvector.Vector `gorm:"type:vector(1536)"`        // 1536 dim for text-embedding-3-small
    CreatedAt time.Time
}
```

**Index Strategy:**
Use HNSW (Hierarchical Navigable Small World) for fast approximate nearest neighbor search.
```sql
CREATE INDEX ON email_embeddings USING hnsw (embedding vector_cosine_ops);
```

## 2. Interfaces & Algorithms

### 2.1 Embedding Provider
```go
// pkg/ai/interface.go

type EmbeddingProvider interface {
    // GetEmbeddings returns a list of vectors for a list of text chunks.
    // Batching is recommended for performance.
    GetEmbeddings(ctx context.Context, texts []string) ([][]float32, error)
}
```

### 2.2 Chunking Strategy (`pkg/utils/chunker.go`)
**Algorithm**: Recursive Character Split or Simple Paragraph Split.
**Parameters**:
*   `MaxTokens`: 512 (approx 2000 chars).
*   `Overlap`: 50 tokens (to maintain context between chunks).

**Logic**:
1.  Strip HTML tags (use `microcosm-cc/bluemonday` or regex for MVP).
2.  Split by double newline `\n\n` (paragraphs).
3.  If a paragraph > MaxTokens, split by sentence `.` .
4.  Merge small paragraphs until MaxTokens is reached.

## 3. Search Service (`internal/service/search.go`)

### 3.1 Query Logic
Input: `query` (string), `topK` (int).

1.  **Embed Query**: Call `aiProvider.GetEmbeddings` for the user query.
2.  **Vector Search**:
    ```sql
    SELECT 
        e.id as email_id, 
        e.subject, 
        e.sender, 
        e.date,
        emb.content as snippet,
        1 - (emb.embedding <=> ?) as score -- Cosine Similarity
    FROM email_embeddings emb
    JOIN emails e ON emb.email_id = e.id
    WHERE 1 - (emb.embedding <=> ?) > 0.7 -- Similarity Threshold
    ORDER BY score DESC
    LIMIT 10;
    ```
3.  **Deduplication**: If multiple chunks from the same email match, group them and return the email only once (using the highest score).

## 4. API Contract

### 4.1 `GET /api/v1/search`

**Request**:
*   `q`: string (required) - The natural language query.
*   `limit`: int (optional, default 20).

**Response**:
```json
{
  "results": [
    {
      "email_id": "uuid...",
      "subject": "Q4 Budget Plan",
      "sender": "alice@corp.com",
      "date": "2025-11-20T10:00:00Z",
      "score": 0.89,
      "snippet": "...the marketing budget for Q4 needs to be increased by 15%..."
    }
  ],
  "latency_ms": 120
}
```

## 5. Dependency Changes

### Go Modules
*   `github.com/pgvector/pgvector-go`: For mapping vector types in GORM.
*   `github.com/sashabaranov/go-openai`: Already used, check if it supports Embeddings API.

### Infrastructure
*   **Docker Image**: `postgres:16-alpine` -> `pgvector/pgvector:pg16`.
*   **Config**: Add `AI.EmbeddingModel` to `config.yaml`.
