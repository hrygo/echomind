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
