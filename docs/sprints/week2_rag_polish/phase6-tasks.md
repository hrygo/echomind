# Phase 6: Team Collaboration - Task Breakdown

> **Parent**: [Phase 6 Design](./phase6-design.md)
> **Status**: Planning
> **Est. Duration**: 2 Weeks

---

## 🚀 Phase 6.1: Foundation (Backend)

**Goal**: Introduce `Organization` concept without breaking existing "Personal" flow.

- [ ] **Database Migrations**
    - [ ] Create `organizations` table (ID, Name, Slug, OwnerID).
    - [ ] Create `organization_members` table (OrgID, UserID, Role).
    - [ ] Add `organization_id` to `users` (optional, for default context).
- **Data Migration Script**
    - [ ] For every existing User:
        - [ ] Create an Organization named "{User}'s Workspace".
        - [ ] Add User as "Owner" of this Org.
- **Backend Core**
    - [ ] Update `AuthMiddleware` to extract `X-Organization-ID`.
    - [ ] Implement `RequireOrgRole(role)` middleware helper.
    - [ ] Update `GORM` scopes to automatically filter by `OrganizationID`.

## 🛠 Phase 6.2: API Implementation

**Goal**: Allow users to manage organizations and members.

- [ ] **Organization API**
    - [ ] `POST /api/v1/orgs`: Create new organization.
    - [ ] `GET /api/v1/orgs`: List user's organizations.
    - [ ] `GET /api/v1/orgs/:id`: Get details (if member).
    - [ ] `PUT /api/v1/orgs/:id`: Update settings (Admin+).
- [ ] **Member API**
    - [ ] `GET /api/v1/orgs/:id/members`: List members.
    - [ ] `POST /api/v1/orgs/:id/invites`: Invite user (by email).
    - [ ] `DELETE /api/v1/orgs/:id/members/:uid`: Remove member.

## 🎨 Phase 6.3: Frontend Core

**Goal**: UI for switching contexts and managing teams.

- [ ] **State Management**
    - [ ] Create `useOrganization` store (Zustand).
    - [ ] Persist selected Org ID in `localStorage`.
- [ ] **Components**
    - [ ] `OrgSwitcher`: Dropdown in sidebar header.
    - [ ] `CreateOrgModal`: Form to start a new workspace.
    - [ ] `MemberSettings`: List with "Invite" button.
- [ ] **Router Interception**
    - [ ] Redirect to `/login` if no Org selected (or default to Personal).
    - [ ] Ensure all API calls include `X-Organization-ID`.

## 📦 Phase 6.4: Resource Sharing (Deep Dive)

**Goal**: Allow Email Accounts to be shared.

- [ ] **Schema Update**
    - [ ] `EmailAccount` add `TeamID` and `OrganizationID`.
    - [ ] Remove unique constraint on `UserID` (allow null).
- [ ] **Logic Update**
    - [ ] Update `SyncWorker` to handle shared accounts.
    - [ ] Update `SearchService` to search across team resources.
