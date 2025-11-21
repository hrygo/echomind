# 📐 Technical Design: Phase 4 - Deep Insight & Relationship Intelligence

**Version**: v1.0
**Status**: Planned
**Target Release**: v0.5.0

---

## 1. Overview
Phase 4 shifts the focus from "Message Management" to "Relationship Intelligence". By aggregating data from synchronized emails, we will build a dynamic profile of the user's network and provide AI-assisted workflows.

## 2. Core Features

### 2.1 Contact Intelligence
*   **Automatic Extraction**: Contacts are automatically created/updated when emails are synced.
*   **Stats Aggregation**: Track `InteractionCount`, `LastInteractedAt`, and `SentimentScore` per contact.
*   **Goal**: Answer "Who are my most important connections?"

### 2.2 Relationship Graph
*   **Visual Network**: A force-directed graph showing connections.
*   **Logic**: 
    *   **Nodes**: Contacts.
    *   **Weight**: Number of interactions.
    *   **Links**: (Future) Derived from Cc/To patterns. For v0.5.0, we focus on User-Contact star topology first.

### 2.3 Smart Reply (AI Agent)
*   **Contextual Drafting**: Generate email replies based on the original thread and user intent.
*   **Tone Adjustment**: Allow selecting "Professional", "Casual", or "Firm".

---

## 3. Data Model Changes

### 3.1 Table: `contacts` (Enhancement)
Existing table needs expansion or a companion stats table. We will expand the existing table for simplicity.

```go
type Contact struct {
    ID               uuid.UUID `gorm:"type:uuid;primary_key"`
    UserID           uuid.UUID `gorm:"index;not null"`
    Email            string    `gorm:"index;not null"`
    Name             string
    
    // New Metrics
    InteractionCount int       `gorm:"default:0"`
    LastInteractedAt time.Time
    AvgSentiment     float64   `gorm:"type:decimal(3,2);default:0"` // Range: -1.0 to 1.0
    
    CreatedAt        time.Time
    UpdatedAt        time.Time
}
```

### 3.2 Logic: Aggregation
*   **Trigger**: During `HandleEmailAnalyzeTask` (Async Worker).
*   **Process**:
    1.  Extract Sender/Receiver email.
    2.  `Upsert` Contact record.
    3.  Increment `InteractionCount`.
    4.  Update `LastInteractedAt`.
    5.  Update rolling average of `AvgSentiment`.

---

## 4. API Specifications

### 4.1 Relationship Data
*   **GET /api/v1/insights/network**
    *   **Response**:
        ```json
        {
          "nodes": [
            {"id": "email1", "name": "Alice", "val": 50}, // val = interaction count
            {"id": "email2", "name": "Bob", "val": 20}
          ],
          "links": [] // Future: Connections between Alice and Bob
        }
        ```

### 4.2 AI Draft
*   **POST /api/v1/ai/draft**
    *   **Payload**:
        ```json
        {
          "email_id": "uuid",
          "instruction": "Decline politely but keep door open",
          "tone": "Professional"
        }
        ```
    *   **Response**: `{"draft": "Dear Alice, ..."}`

---

## 5. Implementation Steps

### Step 1: Backend Core (Go)
1.  Update `Contact` model.
2.  Modify `internal/tasks/analyze.go` to update contact stats during email analysis.
3.  Create a migration script to backfill stats for existing emails.

### Step 2: Insight APIs
1.  Create `internal/handler/insight.go`.
2.  Implement `GetNetworkGraph` using optimized SQL aggregation.

### Step 3: AI Draft Module
1.  Update `pkg/ai` provider interface to support `GenerateDraft(context, thread, instruction)`.
2.  Implement logic in `DeepSeek`/`Gemini` providers.

### Step 4: Frontend (Next.js)
1.  Add `react-force-graph-2d` or `recharts`.
2.  Create `NetworkGraph` component on Dashboard.
3.  Add "Reply with AI" button in Email Detail view.

---

## 6. Testing Strategy
*   **Unit**: Test `GenerateDraft` prompt construction. Test Contact stats aggregation logic.
*   **Integration**: Sync a thread -> Verify Contact stats incremented -> Verify Graph API returns new node.
