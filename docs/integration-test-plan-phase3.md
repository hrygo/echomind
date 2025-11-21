# 🧪 Phase 3: Integration Test Plan & Guide

**Objective**: Validate the full "Real-World Sync" lifecycle, ensuring that user credentials can be securely stored, decrypted, and used to perform a dynamic IMAP sync.

**Target Audience**: Junior Developers / QA Engineers

---

## 1. Pre-requisites

Before running the tests, ensure your local environment is set up correctly.

1.  **Infrastructure**:
    *   Docker is running.
    *   Start Postgres & Redis: `make docker-up`
2.  **Configuration**:
    *   Ensure `backend/configs/config.yaml` exists (copy from `config.example.yaml`).
    *   **CRITICAL**: Set a valid encryption key in `backend/configs/config.yaml`.
        ```yaml
        security:
          encryption_key: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" # 64-char hex string
        ```
        *Tip: You can generate one using `openssl rand -hex 32`.*

---

## 2. Test Scenarios

We will execute two main scenarios:
1.  **The "Happy Path" (Real Gmail/Outlook Account)**: Verifies actual connectivity.
2.  **The "Mock Path" (Automated/Local)**: Verifies internal logic without needing real credentials.

---

### 🟢 Scenario A: The "Happy Path" (Manual Verification)

**Goal**: Connect a REAL Gmail account and sync emails.

**Steps**:

1.  **Prepare a Test Gmail Account**:
    *   Enable IMAP in Gmail Settings.
    *   Enable 2-Step Verification.
    *   Generate an **App Password** (Security -> App passwords). *Do not use your main password.*

2.  **Start the Full Stack**:
    *   `make dev` (Starts Backend, Worker, Frontend).
    *   Open `http://localhost:3000`.

3.  **Execute**:
    *   **Register**: Create a new user (e.g., `test@example.com`).
    *   **Navigate**: Go to `Settings`.
    *   **Connect**: Enter your Gmail details:
        *   **Email**: `your.email@gmail.com`
        *   **Server**: `imap.gmail.com`
        *   **Port**: `993`
        *   **Username**: `your.email@gmail.com`
        *   **Password**: `<Your 16-char App Password>`
    *   **Verify UI**: Click "Connect Account". Expect a success message and the UI to switch to the "Connected" state.

4.  **Trigger Sync**:
    *   Go to `Inbox` (Dashboard).
    *   Click "Sync Now".
    *   Wait 5-10 seconds.

5.  **Validation**:
    *   **Frontend**: Do you see your latest emails in the list?
    *   **Backend Logs**: Check terminal for "Enqueued analysis task..." logs.
    *   **Database**: Check `email_accounts` table: `encrypted_password` should NOT be readable text.

---

### 🟡 Scenario B: The "Mock Path" (Automated Integration Test)

**Goal**: Verify the API flow and DB state programmatically.

**Tool**: `curl` or Postman.

**Steps**:

1.  **Start Backend**: `make run-backend`

2.  **Register User**:
    ```bash
    curl -X POST http://localhost:8080/api/v1/auth/register \
      -H "Content-Type: application/json" \
      -d '{"email": "integration@test.com", "password": "password123", "name": "Integration Tester"}'
    ```
    *Save the `token` from the response.*

3.  **Connect Account (Mock)**:
    *   *Note: This will fail connection test unless we mock the `DialTLS` in code, but for now, we test the API response behavior.*
    ```bash
    curl -X POST http://localhost:8080/api/v1/settings/account \
      -H "Authorization: Bearer <YOUR_TOKEN>" \
      -H "Content-Type: application/json" \
      -d '{ 
        "email": "integration@test.com",
        "server_address": "imap.gmail.com",
        "server_port": 993,
        "username": "integration@test.com",
        "password": "invalid-password-for-test"
      }'
    ```
    *   **Expected Result**: `400 Bad Request` with error "IMAP connection test failed...".
    *   **Why**: This confirms the backend is actually trying to connect and validating credentials before saving.

---

## 3. Troubleshooting Guide

*   **Error: "IMAP connection failed: dial tcp..."**:
    *   Check internet connection.
    *   Verify `server_address` and `server_port` (993 for SSL).
*   **Error: "IMAP login failed: [ALERT] Application-specific password required..."**:
    *   You are using your main Gmail password. Generate an **App Password**.
*   **Error: "failed to decrypt password"**:
    *   Did you change the `encryption_key` in `config.yaml` *after* connecting an account? Changing the key invalidates all stored passwords. Re-connect the account.

---

## 4. Code Review Checklist (For Senior Dev)

*   [ ] **No secrets in logs**: Search logs for the App Password used during testing.
*   [ ] **Encryption**: Verify `pkg/utils/crypto.go` uses `AES-256-GCM` (Galois/Counter Mode).
*   [ ] **Connection Cleanup**: Ensure `defer client.Logout()` or `Close()` is called in `SyncService`.
