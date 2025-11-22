# v0.9.0 Detailed Design: Actionable Intelligence

> **Phase**: 6.2
> **Version Target**: v0.9.0
> **Duration**: 4 Weeks
> **Theme**: From Insight to Action

---

## 1. Core Philosophy: The "Active" Dashboard

Moving from *Reading* to *Doing*. The Dashboard becomes a command center where decisions are made instantly, without navigating away.

---

## 2. Feature Specifications

### 2.1 Actionable Dashboard Cards

**Goal**: Enable "One-Click Decisions" directly from the Briefing view.

#### UI Components
1.  **Pending Decision Card**:
    *   **Content**: AI Summary of the request + Key Entities (Sender, Deadline).
    *   **Actions**:
        *   `[Approve]`: Sends a pre-generated "Approved" reply, archives email.
        *   `[Reply...]`: Opens a mini-editor with AI draft options.
        *   `[Snooze]`: Hides card for 4h/Tomorrow.
    *   **Interaction**:
        *   Click Action -> Card shows loading spinner -> Card fades out (Optimistic UI).
        *   "Undo" toast appears for 5 seconds.

2.  **Risk Warning Card**:
    *   **Actions**:
        *   `[Dismiss]`: Marks risk as handled.
        *   `[Investigate]`: Opens Chat Copilot with context pre-loaded ("Why is this high risk?").

#### API Endpoints
*   `POST /api/v1/actions/approve`: `{ email_id: "uuid" }`
*   `POST /api/v1/actions/snooze`: `{ email_id: "uuid", duration: "4h" }`

---

### 2.2 Smart Contexts (The "Project" View)

**Goal**: Slice the massive inbox into manageable "Attention Scopes".

#### Logic & Rules
A `Context` is a dynamic filter defined by:
1.  **Keywords**: "Project Alpha", "Budget", "Q4".
2.  **Key Stakeholders**: List of email addresses (client@corp.com).
3.  **Timeframe**: "Last 30 days" (rolling) or "Oct 1 - Dec 31" (fixed).

#### Interaction Flow
1.  **Creation**: User clicks "+" in Sidebar -> "Create Smart Context".
    *   Input: Name, Keywords, Key People.
    *   Preview: "Found 42 matching emails".
2.  **Activation**: Clicking a Context in Sidebar (`/dashboard?context=project_alpha`).
    *   **Dashboard**: Re-calculates stats (Risks, Tasks) *only* for this context.
    *   **Search/Chat**: RAG scope is limited to documents/emails in this context.
3.  **Auto-Tagging**:
    *   Backend `AnalyzeTask` checks new emails against active Context rules.
    *   If match: Adds `context_ids` to Email metadata (for fast filtering).

#### API Endpoints
*   `POST /api/v1/contexts`: Create definition.
*   `GET /api/v1/dashboard/stats?context_id=...`: Scoped stats.

---

### 2.3 Task Hub (Internal Task System)

**Goal**: Centralize "To-Dos" from emails and manual entry, preparing for WeChat push.

#### Data Structure (`tasks` table)
*   `source`: "email" | "manual" | "ai_inference"
*   `notify_wechat`: boolean (Default true for High priority).
*   `status`: "todo" | "in_progress" | "done" | "archived"

#### Integration Points
*   **Smart Action**: Clicking "Create Task" in Email Detail -> `POST /api/v1/tasks`.
*   **Dashboard Widget**: A dedicated "Action Items" list replacing the current static list.
    *   Support: Checkbox (Complete), Edit (Rename/Reschedule).

---

## 3. Technical Architecture & Schema

### 3.1 Database Schema (Postgres)

```sql
CREATE TABLE contexts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    name VARCHAR(100) NOT NULL,
    color VARCHAR(20) DEFAULT 'blue',
    keywords TEXT[], -- Array of strings
    stakeholders TEXT[], -- Array of email addresses
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    source_email_id UUID REFERENCES emails(id), -- Optional link to email
    context_id UUID REFERENCES contexts(id),    -- Optional link to context
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(20) DEFAULT 'todo', -- todo, done
    priority VARCHAR(20) DEFAULT 'medium', -- high, medium, low
    due_date TIMESTAMPTZ,
    notify_wechat BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Many-to-Many for Email-Context caching (Optimization)
CREATE TABLE email_contexts (
    email_id UUID REFERENCES emails(id),
    context_id UUID REFERENCES contexts(id),
    PRIMARY KEY (email_id, context_id)
);
```

### 3.2 RAG Pipeline Update
*   **Ingestion**: When indexing an email, run it against `Context` rules. If match, add `context_id` to Vector Metadata.
*   **Query**: `SearchService.Search(query, contextID)` -> Adds filter `metadata['context_id'] == contextID`.

---

## 4. Implementation Roadmap (Weekly)

### Week 1: The Task Engine
*   [BE] `tasks` Table Migration & Model.
*   [BE] `TaskService` (CRUD).
*   [FE] `TaskWidget` component (List, Checkbox, optimistic updates).
*   [Integration] Connect Email Detail "Create Task" button to API.

### Week 2: The Context Brain
*   [BE] `contexts` Table & `ContextService`.
*   [BE] Rule Matcher Logic (Regex/Keyword matching).
*   [FE] Context Sidebar UI & Creator Modal.
*   [BE] Background job: Backfill contexts for existing emails.

### Week 3: Actionable Dashboard
*   [FE] Refactor Dashboard Cards to support "Actions".
*   [BE] `ActionService`: Handle `Approve`, `Snooze` (simple implementation: tag update + archive).
*   [FE] "Undo" Toast mechanism.

### Week 4: Chat & Polish
*   [AI] Prompt Engineering: "Extract tasks as JSON for Task Hub".
*   [FE] Chat Widget: Render `TaskCard` in chat stream.
*   [QA] End-to-end testing of the "Email -> Decision -> Archive" loop.