# 🏗️ Technical Designs Archive


# [Source: v0.5.1_phase5_planning/tech-design-phase3.md]
# 📐 Technical Design: Phase 3 - Real-World Sync Integration

**Version**: v1.0
**Status**: Approved
**Target Release**: v0.4.0

---

## 1. Overview
This phase focuses on transitioning EchoMind from a mock sync environment to a real-world IMAP integration. The core goal is to allow users to securely configure their own email accounts (Gmail, Outlook, etc.) and enable the backend to dynamically connect and sync emails on a per-user basis.

## 2. Architecture Changes

### 2.1 Security & Credential Storage
*   **Problem**: Currently, no user credentials are stored. We cannot connect to real user accounts.
*   **Solution**: Store IMAP credentials (Host, Port, Username, Password) in a new database table.
*   **Encryption**: Passwords MUST be encrypted at rest using **AES-256-GCM**. The encryption key will be stored in the server configuration (Env Var), not in the database.

### 2.2 Dynamic Connection Model
*   **Before**: `SyncService` used a single, global `*client.Client` initialized at startup (which failed if no hardcoded creds were present).
*   **After**: `SyncService` will:
    1.  Retrieve the user's encrypted credentials from DB.
    2.  Decrypt the password.
    3.  Establish a **transient** TCP/TLS connection to the user's IMAP server.
    4.  Perform the sync.
    5.  Close the connection immediately.

---

## 3. Data Model

### 3.1 Database Schema (PostgreSQL)
**Table**: `email_accounts`

| Column | Type | Constraints | Description |
| :--- | :--- | :--- | :--- |
| `id` | `UUID` | PK, Default: `gen_random_uuid()` | Unique Account ID |
| `user_id` | `UUID` | FK -> `users(id)`, Unique, Not Null | Owner of the account |
| `email` | `VARCHAR(255)` | Not Null | The email address (display) |
| `server_address` | `VARCHAR(255)` | Not Null | e.g., `imap.gmail.com` |
| `server_port` | `INT` | Not Null, Default: 993 | e.g., 993 |
| `username` | `VARCHAR(255)` | Not Null | IMAP Login Username |
| `encrypted_password` | `TEXT` | Not Null | Base64 encoded ciphertext |
| `is_connected` | `BOOLEAN` | Default: `FALSE` | Status flag |
| `last_sync_at` | `TIMESTAMP` | Nullable | Timestamp of last success |
| `error_message` | `TEXT` | Nullable | Last connection error |
| `created_at` | `TIMESTAMP` | Default: `NOW()` | |
| `updated_at` | `TIMESTAMP` | Default: `NOW()` | |

### 3.2 Go Struct (`internal/model/account.go`)
```go
type EmailAccount struct {
    ID                uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    UserID            uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
    
    Email             string    `gorm:"not null"`
    ServerAddress     string    `gorm:"not null"`
    ServerPort        int       `gorm:"not null;default:993"`
    Username          string    `gorm:"not null"`
    
    EncryptedPassword string    `gorm:"not null"` 
    
    IsConnected       bool      `gorm:"default:false"`
    LastSyncAt        *time.Time
    ErrorMessage      string
    
    CreatedAt         time.Time
    UpdatedAt         time.Time
}
```

---

## 4. API Design

### 4.1 Connect / Update Account
*   **Endpoint**: `POST /api/v1/settings/account`
*   **Auth**: Bearer Token
*   **Payload**:
    ```json
    {
      "email": "user@example.com",
      "server_address": "imap.gmail.com",
      "server_port": 993,
      "username": "user@example.com",
      "password": "app_specific_password"
    }
    ```
*   **Logic**:
    1.  Validate input.
    2.  **Test Connection**: Attempt to Dial & Login to the IMAP server immediately.
    3.  If fail -> Return 400 with error message.
    4.  If success -> Encrypt password -> Upsert into `email_accounts`.
*   **Response**: `200 OK`
    ```json
    { "message": "Account connected successfully" }
    ```

### 4.2 Get Account Status
*   **Endpoint**: `GET /api/v1/settings/account`
*   **Auth**: Bearer Token
*   **Response**: `200 OK`
    ```json
    {
      "has_account": true,
      "email": "user@example.com",
      "server_address": "imap.gmail.com",
      "is_connected": true,
      "last_sync_at": "2023..."
      // NEVER return the password
    }
    ```
    If no account exists, return `{ "has_account": false }`.

---

## 5. Implementation Steps

### Step 1: Security Infrastructure
1.  Add `security.encryption_key` to `configs/config.example.yaml`.
2.  Implement `pkg/utils/crypto.go` with `Encrypt` and `Decrypt` functions using AES-GCM.

### Step 2: Backend Models & Service
1.  Create `internal/model/account.go`.
2.  Create `internal/service/account.go` to handle account storage and connection testing.
3.  Update `internal/service/sync.go` to fetch credentials and use them.

### Step 3: API Layer
1.  Create `internal/handler/account.go`.
2.  Register routes in `cmd/main.go`.

### Step 4: Frontend UI
1.  Create `src/app/dashboard/settings/page.tsx`.
2.  Build a form for IMAP details.
3.  Integrate with the new APIs.

---

## 6. Quality Assurance & Acceptance Criteria

### 6.1 Acceptance Criteria (Definition of Done)
1.  **Security**:
    *   [ ] Database inspection confirms `encrypted_password` is NOT human-readable.
    *   [ ] Decryption fails (returns error) if the wrong key is used (verified via unit test).
2.  **Connectivity**:
    *   [ ] Users can successfully connect a valid Gmail/Outlook account (using App Password).
    *   [ ] Invalid credentials return a user-friendly error message ("Authentication failed", not "500 Internal Server Error").
    *   [ ] Sync successfully fetches emails from the *configured* account, not a hardcoded one.
3.  **Resilience**:
    *   [ ] If the IMAP server is down, the application handles the timeout gracefully and logs the error in `email_accounts.error_message`.
4.  **UI/UX**:
    *   [ ] Settings page correctly reflects the "Connected" state after a successful setup.
    *   [ ] Users can disconnect/delete their account configuration.

### 6.2 Testing Requirements
*   **Unit Test Coverage**: Minimum **80%** coverage for:
    *   `pkg/utils/crypto.go` (Must test encrypt/decrypt roundtrip).
    *   `internal/service/account.go` (Mock database and IMAP client).
*   **Integration Tests**:
    *   **Critical Path**: `TestAccountLifecycle`
        1.  Register a new user.
        2.  POST /settings/account with *mock* valid credentials.
        3.  Verify DB record is created and encrypted.
        4.  Trigger Sync.
        5.  Verify `SyncService` calls the Mock IMAP Client with decrypted credentials.

---

## 7. Code Review & Standards

### 7.1 Review Checklist
*   **No Secrets in Logs**: Ensure passwords (raw or encrypted) are NEVER logged to stdout/Zap.
*   **Error Handling**: All IMAP connection errors must be wrapped (e.g., `fmt.Errorf("imap connection failed: %w", err)`).
*   **Context Usage**: `SyncEmails` must respect `context.Context` for cancellation/timeout.

### 7.2 Branching Strategy
*   Create feature branches from `main`: `feat/phase3-crypto`, `feat/phase3-backend-sync`, `feat/phase3-frontend-settings`.
*   Squash merge into `main` only after CI passes.

---

## 8. Integration Test Plan (Key Path)

Since we cannot use real Gmail passwords in CI, we will use the **Mock IMAP Provider** already present in `pkg/imap`.

1.  **Setup**:
    *   Spin up the Go Backend with an in-memory SQLite DB (or Dockerized Postgres).
    *   Configure `pkg/imap/mock_test.go` to act as the server.
2.  **Execution**:
    *   Call API: `POST /auth/register` (User A).
    *   Call API: `POST /settings/account` (User A, Mock Creds).
    *   Call API: `POST /sync` (User A).
3.  **Verification**:
    *   Check `GET /emails` returns the list of mock emails defined in the Mock Provider.

---

## 9. Future Considerations
*   **OAuth2**: For Gmail/Outlook, we should eventually support OAuth2 instead of App Passwords for better UX and security. This is out of scope for v0.4.0 but noted for future phases.


# [Source: v0.5.1_phase5_planning/tech-design-phase4.md]
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


# [Source: v0.5.1_phase5_planning/tech-design-phase5.md]
# 📐 Technical Design: Phase 5 - Commercialization & Teams

**Version**: v0.1 (Draft)
**Status**: Planned
**Target Release**: v0.6.0

---

## 1. Overview
Phase 5 focuses on transforming EchoMind from a single-user tool into a sustainable SaaS product. This involves implementing payment gateways (Stripe), enforcing usage quotas, and enabling multi-user collaboration within organizations.

## 2. Core Features

### 2.1 Monetization (Stripe Integration)
*   **Subscription Models**: Free Tier (Limited), Pro Tier (Unlimited AI, Priority Support).
*   **Payment Processing**: Securely handle checkout sessions and recurring billing via Stripe.
*   **Webhook Handling**: React to subscription lifecycle events (created, updated, canceled, payment_failed).

### 2.2 Usage Limits & Quotas
*   **Tracking**: Count AI operations (summaries, drafts, analysis) per user/org per billing period.
*   **Enforcement**: Block requests when limits are exceeded with upgrade prompts.
*   **Reset**: Reset quotas automatically at the start of each billing cycle.

### 2.3 Team Collaboration
*   **Organizations**: Users belong to an Organization.
*   **Roles**: Admin (Billing, User Mgmt), Member (Standard access).
*   **Data Isolation**: Ensure strict separation between organizations while allowing sharing within.

---

## 3. Data Model Changes

### 3.1 Table: `organizations`
```go
type Organization struct {
    ID            uuid.UUID `gorm:"type:uuid;primary_key"`
    Name          string
    StripeCustID  string    `gorm:"index"` // Stripe Customer ID
    SubStatus     string    // active, past_due, canceled, etc.
    SubPlan       string    // free, pro
    CurrentPeriodEnd time.Time
}
```

### 3.2 Table: `users` (Update)
*   Add `OrganizationID` (FK).
*   Add `Role` (string).

### 3.3 Table: `usage_records`
```go
type UsageRecord struct {
    ID             uuid.UUID
    OrganizationID uuid.UUID `gorm:"index"`
    Metric         string    // e.g., "ai_summary_count"
    Count          int
    PeriodStart    time.Time
    PeriodEnd      time.Time
}
```

---

## 4. API Specifications

### 4.1 Billing
*   **POST /api/v1/billing/checkout**: Create Stripe Checkout Session.
*   **POST /api/v1/billing/portal**: Create Customer Portal Session.
*   **POST /api/v1/webhooks/stripe**: Handle async events.

### 4.2 Team Management
*   **POST /api/v1/team/invite**: Invite user via email.
*   **GET /api/v1/team/members**: List organization members.

---

## 5. Implementation Steps

### Step 1: Stripe Foundation
1.  Set up Stripe account and keys in `config.yaml`.
2.  Implement `StripeService` wrapper for Go SDK.
3.  Create webhook handler for `invoice.payment_succeeded`, `customer.subscription.updated`.

### Step 2: Organization & Quota Logic
1.  Migrate DB schema (`Organization`, `UsageRecord`).
2.  Update `AuthMiddleware` to inject Org context.
3.  Implement `QuotaService` to check/increment usage before AI calls.

### Step 3: Frontend Billing UI
1.  Create "Subscription" page in Settings.
2.  Display current plan, usage bars, and "Upgrade" buttons.
3.  Handle Stripe redirects.

---

## 6. Security Considerations
*   **Webhook Signature Verification**: strictly verify Stripe signatures.
*   **Idempotency**: Ensure webhook events are processed exactly once.
*   **Access Control**: Strict RBAC for billing actions (Admin only).


# [Source: v0.6.0_rag/tech-design.md]
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


# [Source: v0.7.0_ui_refactor/design-spec.md]
# 🎨 Design Spec: EchoMind v0.7.0 "Neural Interface"

> **Status**: Draft (Planned for Post-v0.6.0)
> **Inspiration**: Google NotebookLM, Notion, Perplexity
> **Goal**: Transform EchoMind from an "Email Client" to an "Executive Intelligence OS".

---

## 1. Vision & Philosophy

### 1.1 The Paradigm Shift
*   **Old World (v0.5)**: "Inbox Zero". Users manually process lists of emails.
*   **Current World (v0.6 RAG)**: "Search". Users query a database.
*   **New World (v0.7)**: **"Contextual Intelligence"**. The UI is a fluid canvas where AI proactively synthesizes information and reacts to user intent.

### 1.2 Core Principles
1.  **AI First, Not Mail First**: The primary interface is *generated content* (Briefings, Answers), not raw data (Email Lists). Raw emails are only shown as *citations*.
2.  **Context is King**: Users define "Context Windows" (e.g., "Project Alpha", "Key Clients"). The AI adapts its answers based on the active context.
3.  **Generative UI**: The interface adapts to the content. A query about dates renders a Calendar widget; a query about relationships renders a Network Graph.

---

## 2. The "Fluid Canvas" Layout

We abandon the traditional 3-column email layout for a **Source-Canvas-Detail** architecture.

### 2.1 Zone A: The Context Manager (Left Sidebar)
*   **Concept**: Instead of static folders, this area manages **"Attention Scopes"**.
*   **Sections**:
    *   **Auto-Clusters**: "Priority Inbox", "Needs Reply", "Follow-ups".
    *   **Smart Contexts** (User Defined): "Q4 Fundraising", "Legal Issues", "Person: John Doe".
    *   **Source Data**: Attached documents (PDF/Word), Spreadsheets extracted from emails.
*   **Interaction**:
    *   Clicking a Context *filters* the RAG engine for the main view.
    *   Checkbox selection allows multi-context synthesis ("Compare *Q4 Budget* with *Q3 Spend*").

### 2.2 Zone B: The Intelligence Canvas (Center Stage)
This is the main workspace. It has two distinct modes:

#### Mode 1: The Briefing (Passive / "Morning Coffee" Mode)
*   **Analogy**: A personalized newspaper or Notion dashboard.
*   **Content**:
    *   **The Lead Story**: "You have 3 critical decisions today." (Summary of top risks/requests).
    *   **Deal Watch**: "Acme Corp contract is stalling." (Dealmaker insight).
    *   **Pulse**: A mini-heatmap of today's communication volume.
*   **Actionability**: Each card has quick actions (`[Draft Reply]`, `[Snooze]`, `[Delegate]`).

#### Mode 2: The Studio (Active / "Deep Work" Mode)
*   **Analogy**: NotebookLM / ChatGPT / Perplexity.
*   **Trigger**: Clicking a Briefing card or typing in the Omni-Bar.
*   **Features**:
    *   **Streaming Q&A**: Answers appear in real-time.
    *   **Rich Citations**: Every claim has a `[1]` footnote. Hovering shows a snippet; clicking opens the source.
    *   **Embedded Widgets**:
        *   "Show me the timeline" -> Renders a Timeline Component *inside* the chat.
        *   "Who is involved?" -> Renders a Contact Card.

### 2.3 Zone C: The Grounding Panel (Right Slide-over / Overlay)
*   **Purpose**: Truth verification.
*   **Content**: The raw email thread or document.
*   **Behavior**:
    *   Hidden by default.
    *   Slides in when a citation `[x]` is clicked.
    *   Highlights the exact paragraph referenced by the AI.

---

## 3. Generative UI Components (The "Widgets")

The Chat interface won't just output text. It will render React Components based on the **Intent Classification** from the backend.

| Intent | UI Component | Description |
| :--- | :--- | :--- |
| `intent:scheduling` | `<CalendarWidget />` | Interactive calendar slot picker. |
| `intent:relationship` | `<NetworkGraph />` | (Existing) Force-directed graph of connections. |
| `intent:finance` | `<DataTabel />` | Extracted numbers in a clean table. |
| `intent:decision` | `<ApprovalCard />` | Big "Approve / Reject" buttons with risk analysis. |
| `intent:draft` | `<Editor />` | A rich-text editor pre-filled with the draft. |

---

## 4. User Experience Flows

### 4.1 The Executive "Morning Flow"
1.  **Open EchoMind**: Sees "The Briefing".
2.  **Scan**: Reads "The Lead Story" (3 critical emails).
3.  **Act**: Clicks `[Approve]` on Item 1. Clicks `[Draft Reply]` on Item 2.
4.  **Deep Dive**: Item 3 is complex ("Project Delay").
5.  **Transition**: Clicks Item 3. UI shifts to "Studio Mode".
6.  **Query**: "Why is it delayed?" -> AI answers citing 5 emails from the engineering lead.
7.  **Resolve**: Types "Draft an email to Engineering asking for a recovery plan."

### 4.2 The Dealmaker "Hunter Flow"
1.  **Select Context**: Checks "Active Deals" in the Left Sidebar.
2.  **Prompt**: "Who hasn't replied in 7 days?"
3.  **Result**: AI lists 3 people. Renders a `<FollowUpList />` widget.
4.  **Action**: Clicks `[Nudge All]`. AI generates 3 personalized follow-up drafts.

---

## 5. Technical Requirements for v0.6.0 (Preparation)

To enable this UI in v0.7.0, the v0.6.0 Backend/RAG must support:

1.  **Citation Metadata**: Search results **must** return exact `chunk_id` and `text_offset` to allow the UI to highlight the source.
2.  **Structured Threads**: The API must be able to return a full conversation thread structure to populate the "Grounding Panel".
3.  **Streaming JSON**: The API should support streaming structured data (e.g., `text` chunks followed by `widget_data` JSON) to render UI components on the fly.
4.  **Dynamic Filters**: The Vector Store query params must support complex filtering by `sender`, `date`, `thread_id` to support the "Context Manager".

---

## 6. Migration Strategy (The "Burn the Ships" Approach)

*   **Direct Replacement**: v0.7.0 will **overwrite** the existing Dashboard. There will be no toggle to switch back.
*   **Onboarding is Critical**: Since the paradigm shift is massive, the first login must trigger a high-quality "Neural Interface" tutorial (e.g., "Here is your Briefing", "Ask your first question").
*   **Legacy Code Removal**: Delete the old `dashboard/` page components immediately to prevent technical debt accumulation.


# [Source: v0.9.0_actionable_intelligence/detailed-spec.md]
# v0.9.0 Technical Design Specification: Task System & Contexts

> **Version**: 1.0
> **Status**: Approved for Development
> **Target Audience**: Junior/Mid-level Engineers
> **Goal**: Provide step-by-step implementation details for the "Actionable Intelligence" sprint.

---

## 1. Overview

This sprint focuses on two new core entities: **Tasks** and **Contexts**.
*   **Tasks**: A centralized to-do list derived from emails or created manually.
*   **Contexts**: A way to group emails, tasks, and contacts by "Project" or "Topic".

---

## 2. Module 1: Task System (Week 1)

### 2.1 Database Schema (GORM Models)

**File**: `backend/internal/model/task.go`

```go
package model

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskStatus string
type TaskPriority string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"

	TaskPriorityHigh   TaskPriority = "high"
	TaskPriorityMedium TaskPriority = "medium"
	TaskPriorityLow    TaskPriority = "low"
)

type Task struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`

	UserID        uuid.UUID      `gorm:"type:uuid;not null;index"`
	
	// Optional: Link to source email
	SourceEmailID *uuid.UUID     `gorm:"type:uuid"` 
	
	// Optional: Link to a project context (Week 2)
	ContextID     *uuid.UUID     `gorm:"type:uuid"`

	Title         string         `gorm:"not null"`
	Description   string         `gorm:"type:text"`
	
	Status        TaskStatus     `gorm:"type:varchar(20);default:'todo'"`
	Priority      TaskPriority   `gorm:"type:varchar(20);default:'medium'"`
	
	DueDate       *time.Time
	NotifyWeChat  bool           `gorm:"default:false"` // For future use
}
```

### 2.2 API Endpoints

**Base URL**: `/api/v1/tasks`

#### A. Create Task
*   **Method**: `POST /`
*   **Request**:
    ```json
    {
      "title": "Review budget",
      "source_email_id": "uuid...", // Optional
      "due_date": "2023-12-01T10:00:00Z"
    }
    ```
*   **Logic**:
    1.  Parse body.
    2.  Set default `UserID` from context.
    3.  Set default `Status` = "todo".
    4.  Save to DB.

#### B. List Tasks
*   **Method**: `GET /`
*   **Query Params**: `status` (optional), `limit` (default 20).
*   **Response**: `[{ "id": "...", "title": "...", "status": "todo" }]`

#### C. Update Task Status
*   **Method**: `PATCH /:id`
*   **Request**: `{"status": "done"}`
*   **Logic**: Validates status enum before saving.

### 2.3 Frontend Components

**File**: `frontend/src/components/dashboard/TaskWidget.tsx`

*   **UI**: A clean list of checkboxes and titles.
*   **Interaction**:
    *   Clicking checkbox -> Calls `PATCH /tasks/:id` with `status: done`.
    *   **Crucial**: Use **Optimistic UI**.
        1.  Set local state to "done" (strike-through text) immediately.
        2.  Send API request.
        3.  If API fails, revert local state and show Toast error.

---

## 3. Module 2: Smart Contexts (Week 2)

### 3.1 Database Schema

**File**: `backend/internal/model/context.go`

```go
package model

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Context struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index"`
	Name      string         `gorm:"not null"`
	Color     string         `gorm:"default:'blue'"` // For UI tagging
	
	// Rules for auto-matching
	Keywords     datatypes.JSON `gorm:"type:jsonb"` // e.g. ["budget", "finance"]
	Stakeholders datatypes.JSON `gorm:"type:jsonb"` // e.g. ["alice@example.com"]
}
```

### 3.2 Service Logic: Context Matcher

**File**: `backend/internal/service/context.go`

*   **Function**: `Match(email *model.Email, contexts []model.Context) []uuid.UUID`
*   **Logic**:
    1.  Iterate through all user contexts.
    2.  Check if `email.Subject` or `email.Body` contains any `Keywords`.
    3.  Check if `email.Sender` or `email.CC` contains any `Stakeholders`.
    4.  Return IDs of matching contexts.

**Integration Point**:
*   Call this `Match` function inside `tasks/analyze.go` (the analysis worker).
*   Update `Email` struct to store `ContextIDs []uuid.UUID` (Postgres array).

---

## 4. Implementation Steps for Junior Devs

### Day 1: Task Backend
1.  Create `backend/internal/model/task.go` with the struct above.
2.  Run `make dev` to trigger GORM AutoMigrate (check logs to confirm table creation).
3.  Create `backend/internal/service/task.go` (CRUD methods).
4.  Create `backend/internal/handler/task.go` (Gin handlers).
5.  Register routes in `backend/cmd/main.go`.

### Day 2: Task Frontend
1.  Create `frontend/src/lib/api/tasks.ts` (API client functions).
2.  Create `frontend/src/components/dashboard/TaskWidget.tsx`.
3.  Implement `useTaskStore` (Zustand) to manage the list state locally.

### Day 3: Integration
1.  Go to `frontend/src/app/dashboard/email/[id]/page.tsx`.
2.  Find the "Smart Actions" button logic.
3.  Change the `onClick` handler from `alert()` to calling `api.createTask()`.
4.  Test the flow: Open Email -> Click "Create Task" -> Go to Dashboard -> Verify Task appears in widget.

---

## 5. Code Standards

*   **Error Handling**: Always return JSON `{"error": "message"}` with appropriate HTTP codes (400 for bad input, 500 for server error).
*   **Logging**: Use `sugar.Infof` for key actions (e.g., "Task created: UUID").
*   **Types**: Ensure Frontend `interface Task` matches Backend JSON response exactly.


# [Source: week2_rag_polish/phase6-design.md]
# Phase 6: Team Collaboration - Technical Design

> **Status**: Draft
> **Author**: EchoMind AI
> **Date**: 2025-11-22
> **Context**: Transitioning from Single-Tenant (Personal) to Multi-Tenant (Team/Organization).

---

## 1. Overview

The goal of Phase 6 is to enable "Team Collaboration". Currently, EchoMind assumes 1 User = 1 World. We need to introduce a hierarchy where Users belong to Organizations and Teams, allowing shared resources (Email Accounts, Contacts, Knowledge Base).

### Core Concepts

*   **Organization (Org)**: The top-level billing and isolation unit (e.g., "Acme Corp").
*   **Member**: A link between a User and an Organization with a Role (Owner, Admin, Member).
*   **Team**: A subgroup within an Organization (e.g., "Sales", "Support").
*   **Resource Ownership**: Resources (like Email Accounts) can now be owned by a User OR a Team.

---

## 2. Database Schema Changes

### 2.1 New Models

#### Organization
```go
type Organization struct {
    ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    Name        string    `gorm:"not null"`
    Slug        string    `gorm:"uniqueIndex;not null"` // For URL routing /orgs/:slug
    OwnerID     uuid.UUID `gorm:"type:uuid;not null"`   // The super admin
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

#### OrganizationMember
```go
type Role string

const (
    RoleOwner  Role = "owner"
    RoleAdmin  Role = "admin"
    RoleMember Role = "member"
)

type OrganizationMember struct {
    OrganizationID uuid.UUID `gorm:"primaryKey"`
    UserID         uuid.UUID `gorm:"primaryKey"`
    Role           Role      `gorm:"default:'member'"`
    JoinedAt       time.Time
}
```

#### Team
```go
type Team struct {
    ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    OrganizationID uuid.UUID `gorm:"type:uuid;not null;index"`
    Name           string    `gorm:"not null"`
    Description    string
    CreatedAt      time.Time
    UpdatedAt      time.Time
}
```

#### TeamMember
```go
type TeamMember struct {
    TeamID    uuid.UUID `gorm:"primaryKey"`
    UserID    uuid.UUID `gorm:"primaryKey"`
    Role      Role      `gorm:"default:'member'"`
}
```

### 2.2 Modified Models

#### EmailAccount (Breaking Change)
Currently, `EmailAccount` has a hard `UserID`. We need to support Team ownership.

*   **Option A (Polymorphic)**: `OwnerID` (UUID) + `OwnerType` ("user" | "team").
*   **Option B (Nullable Columns)**: `UserID` (nullable), `TeamID` (nullable). *Preferred for referential integrity.*

```go
type EmailAccount struct {
    // ... existing fields
    UserID *uuid.UUID `gorm:"type:uuid;index"` // Nullable if owned by team
    TeamID *uuid.UUID `gorm:"type:uuid;index"` // New field
    
    // Constraints:
    // CHECK (user_id IS NOT NULL OR team_id IS NOT NULL)
}
```

#### Contact
Similar to EmailAccount, contacts can be private (User) or shared (Team/Org).

```go
type Contact struct {
    // ... existing fields
    OrganizationID *uuid.UUID `gorm:"type:uuid;index"` // If set, available to whole org
    TeamID         *uuid.UUID `gorm:"type:uuid;index"` // If set, available to team
    IsPrivate      bool       `gorm:"default:true"`
}
```

---

## 3. API Strategy

### 3.1 New Endpoints

*   `POST /api/v1/orgs`: Create organization.
*   `GET /api/v1/orgs`: List my organizations.
*   `POST /api/v1/orgs/:id/members`: Invite user.
*   `POST /api/v1/teams`: Create team.

### 3.2 Authorization Middleware Update

The `AuthMiddleware` currently checks if the `User` exists. It needs to be enhanced to populate a `Context` with:
*   `CurrentOrgID`: From header `X-Org-ID` or query param?
*   `Permissions`: Calculated based on User's role in that Org.

**Strategy**:
1.  User logs in (Global context).
2.  User selects a "Workspace" (Organization context) in UI.
3.  Frontend sends `X-Organization-ID` header with subsequent requests.
4.  Middleware verifies User is a member of that Org.

---

## 4. Migration Plan

1.  **Phase 6.1 (Foundation)**:
    *   Create `organizations` and `members` tables.
    *   Migrate existing users: Create a "Personal Organization" for each existing user.
    *   Link existing data to this Personal Org.
2.  **Phase 6.2 (Team Logic)**:
    *   Implement `teams` table.
    *   Update `EmailAccount` to support `TeamID`.
3.  **Phase 6.3 (UI)**:
    *   Organization Switcher.
    *   Team Management Settings.

---

## 5. Security & Isolation

*   **Row Level Security (RLS)**: Not native in GORM, but we must ensure every query scopes to `OrganizationID`.
*   **Scopes**:
    *   `db.Where("organization_id = ?", ctx.OrgID)` must be enforced in the Repository layer.

## 6. Open Questions

1.  **Billing**: Is billing per User or per Org? (Assumption: Per Org).
2.  **Invite Flow**: Email invitation with token? (Yes, requires Redis for short-lived tokens).


# [Source: week4_personal_intelligence/detailed-design.md]
# Phase 6.0 Detailed Design: Personal Intelligence Deep-Dive

> **Version**: 1.0  
> **Date**: 2025-11-22  
> **Status**: Approved  

---

## 1. Architecture Overview

This sprint focuses on enhancing the **Personal** experience through three pillars:
1.  **Conversational UI**: A streaming AI Chatbot (Copilot) integrated with RAG.
2.  **Mobile-First Web**: Responsive adaptations for mobile browsers.
3.  **Actionable Intelligence**: Structured extraction of tasks/events from emails.

### Tech Stack Additions
*   **Protocol**: Server-Sent Events (SSE) for Chat Streaming.
*   **Frontend**: Radix UI `Sheet` (for Mobile Sidebar & Copilot Drawer).
*   **AI Model**: `gpt-4o-mini` or `gemini-1.5-flash` (Cost-effective, fast).

---

## 2. API Contracts

### 2.1 Chat Completion (Streaming)
**Endpoint**: `POST /api/v1/chat/completions`

**Request**:
```json
{
  "messages": [
    { "role": "system", "content": "You are EchoMind Copilot..." },
    { "role": "user", "content": "What did Alice say about the budget?" }
  ],
  "context_scope": {
    "user_id": "uuid..." // Implicit from Auth Middleware
  }
}
```

**Response (Content-Type: text/event-stream)**:
```text
data: {"id": "...", "choices": [{"delta": {"content": "Alice "}}]}
data: {"id": "...", "choices": [{"delta": {"content": "mentioned "}}]}
...
data: [DONE]
```

### 2.2 Smart Actions (Email Metadata)
**Endpoint**: `GET /api/v1/emails/:id` (Response Extension)

**Response Field**: `smart_actions` (Array)
```json
[
  {
    "type": "calendar_event",
    "label": "Add to Calendar",
    "data": {
      "title": "Budget Review",
      "start": "2025-11-25T10:00:00Z",
      "end": "2025-11-25T11:00:00Z",
      "location": "Room 303"
    }
  },
  {
    "type": "create_task",
    "label": "Create Todo",
    "data": {
      "title": "Reply with updated figures",
      "priority": "high"
    }
  }
]
```

---

## 3. UI/UX Specifications

### 3.1 AI Copilot (Right Drawer)
*   **Trigger**: "Sparkles" Icon in Header (Right aligned).
*   **Behavior**:
    *   **Desktop**: Slides in from right, overlay or push content. Width: 400px.
    *   **Mobile**: Slides in from bottom or right (Full width).
*   **Components**:
    *   `ChatHistory`: Scrollable area.
    *   `MessageInput`: Textarea with auto-resize.
    *   `Citation`: RAG references rendered as interactive links.

### 3.2 Mobile Navigation
*   **Sidebar**:
    *   **Desktop**: Visible (Sticky).
    *   **Mobile**: Hidden. Triggered by "Hamburger" icon in Header (Left aligned). Use `Sheet` component.
*   **Header**:
    *   **Mobile**: Simplified. Logo + Menu Trigger + Copilot Trigger. Search bar might need to be collapsed into an icon or simplified.

---

## 4. Implementation Guide

### 4.1 Backend: SSE Handler (Gin)
```go
c.Stream(func(w io.Writer) bool {
    if msg, ok := <-streamChannel; ok {
        c.SSEvent("message", msg)
        return true
    }
    return false
})
```

### 4.2 Frontend: Streaming Client
Use `fetch` with `ReadableStream` or a library like `eventsource-parser` (recommended for handling SSE format robustness).

### 4.3 AI Prompting (Smart Actions)
System Prompt addition:
> "Analyze the email content. If specific actionable items are found (meetings, tasks), output a JSON object in the `smart_actions` field following this schema..."

---

## 5. Testing Strategy

*   **Unit**: Test `ActionExtractor` regex/json parsing logic.
*   **E2E (Playwright)**:
    *   `mobile.spec.ts`: Set viewport to `iPhone 12`. Verify Sidebar opens/closes via menu button.
    *   `chat.spec.ts`: Mock `/chat/completions` stream. Verify messages appear incrementally.

## 6. Performance Goals

*   **Chat Latency**: Time to First Token < 1.5s.
*   **Mobile Interaction**: Sidebar open animation < 300ms, no jank.
