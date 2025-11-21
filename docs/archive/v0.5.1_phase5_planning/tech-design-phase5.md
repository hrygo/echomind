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
