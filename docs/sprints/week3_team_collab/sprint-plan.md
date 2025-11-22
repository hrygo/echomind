# Week 3 Sprint Plan: Team Collaboration Foundation

> **Sprint**: Phase 6.1 - Organization Core  
> **Version Target**: v0.7.1  
> **Duration**: 5 days  
> **Dependencies**: Completed Phase 5.3 (v0.7.0-beta)

---

## Sprint Goals

### Primary Objectives
1.  **Multi-Tenancy Core**: Implement `Organization` and `Member` models.
2.  **Context Switching**: Enable frontend to switch between organizations.
3.  **Migration**: Seamlessly upgrade existing single-tenant users to multi-tenant structure.

---

## Day 1: Backend Foundation (Phase 6.1)

### Objective
Implement the database schema and GORM models for Organizations.

### Morning: Database Models
- [x] **Create Models**
    - [x] `backend/internal/model/organization.go`
    - [x] `backend/internal/model/member.go`
- [x] **Update User Model**
    - [x] Add associations in `backend/internal/model/user.go`
- [x] **Auto-Migration**
    - [x] Update `backend/cmd/main.go` to migrate new tables.

### Afternoon: Data Migration Logic
- [x] **Migration Script**
    - [x] Implement logic to ensure every existing user has a "Personal Workspace".
    - [x] Create `backend/cmd/migrate/main.go` (or hook into startup).
- [x] **Tests**
    - [x] Unit tests for Organization model constraints.

---

## Day 2: Organization API (Phase 6.2)

### Objective
Endpoints to create, list, and manage organizations.

### Morning: CRUD API
- [x] **Service Layer**
    - [x] `OrganizationService` (Create, Get, Update).
- [x] **Handlers**
    - [x] `POST /orgs`
    - [x] `GET /orgs`
    - [x] `GET /orgs/:id`

### Afternoon: Membership API
- [x] **Member Management**
    - [x] `GET /orgs/:id/members`
    - [x] `POST /orgs/:id/invites` (Mock email sending for now).

---

## Day 3: Frontend Integration (Phase 6.3)

### Objective
UI for Organization switching.

### Morning: State Management
- [x] **Store**
    - [x] `useOrganization` (Zustand).
    - [x] Fetch user's orgs on login.
- [x] **Interceptor**
    - [x] Add `X-Organization-ID` header to all API requests.

### Afternoon: UI Components
- [x] **Org Switcher**
    - [x] Component in Sidebar.
    - [x] "Create New Organization" Modal.

---

## Day 4: Teams & Shared Resources (Phase 6.4)

### Objective
Allow resources to be owned by Teams.

### Morning: Team Models
- [x] **Team Schema**
    - [x] `Team` model.
    - [x] `TeamMember` model.
- [x] **Resource Update**
    - [x] Update `EmailAccount` schema (nullable UserID, add TeamID).

### Afternoon: Logic Update
- [x] **Service Layer**
    - [x] Update `AccountService` to handle Team ownership.
    - [x] Update `SyncService` permissions.

---

## Day 5: Polish & Integration

### Objective
End-to-end testing of the multi-tenant flow.

- [x] **E2E Tests**
    - [x] User A invites User B (Mocked).
    - [x] User B accepts and switches context (Switcher implemented).
    - [x] Integration tests for Organization API.
- [x] **Documentation**
    - [x] Update API docs with new endpoints.
    - [x] Document multi-tenant architecture.

