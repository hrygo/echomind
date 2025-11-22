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
