# EchoMind API Documentation

Base URL: `/api/v1`

## Authentication

All endpoints (except `/auth/*` and `/health`) require a JWT token in the header:
`Authorization: Bearer <token>`

---

## Search API

### `GET /search`

Performs a semantic search over the user's emails.

**Parameters:**

| Name | Type | Required | Default | Description |
|---|---|---|---|---|
| `q` | string | Yes | - | The natural language query string. |
| `limit` | int | No | 10 | Maximum number of results to return (1-100). |
| `sender` | string | No | - | Filter results by sender name or email (partial match). |
| `start_date` | string | No | - | Filter emails on or after this date (Format: `YYYY-MM-DD`). |
| `end_date` | string | No | - | Filter emails on or before this date (Format: `YYYY-MM-DD`). |

**Response:**

```json
{
  "query": "project update",
  "count": 2,
  "results": [
    {
      "email_id": "uuid-string",
      "subject": "Project Alpha Weekly Update",
      "snippet": "Here is the status of...",
      "sender": "alice@example.com",
      "date": "2025-11-20T10:00:00Z",
      "score": 0.89
    },
    ...
  ]
}
```

**Error Responses:**

*   `400 Bad Request`: Missing `q` parameter.
*   `401 Unauthorized`: Missing or invalid token.
*   `500 Internal Server Error`: Search failure (DB or AI provider issue).

---

## Other Endpoints

### `GET /health`
Returns the system health status.

### `POST /auth/login`
Login with email and password.

### `POST /auth/register`
Register a new user account.
