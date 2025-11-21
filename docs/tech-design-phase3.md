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
