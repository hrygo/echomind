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
