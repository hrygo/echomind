# 🎨 Design Spec: EchoMind v0.7.0 "Neural Interface"

> **Status**: Draft (Planned for Post-v0.6.0)
> **Inspiration**: Google NotebookLM, Notion, Perplexity
> **Goal**: Transform EchoMind from an "Email Client" to an "Executive Intelligence OS".

---

## 1. Vision & Philosophy

### 1.1 The Paradigm Shift
*   **Old World (v0.5)**: "Inbox Zero". Users manually process lists of emails.
*   **Current World (v0.6 RAG)**: "Search". Users query a database.
*   **New World (v0.7)**: **"Contextual Intelligence"**. The UI is a fluid canvas where AI proactively synthesizes information and reacts to user intent.

### 1.2 Core Principles
1.  **AI First, Not Mail First**: The primary interface is *generated content* (Briefings, Answers), not raw data (Email Lists). Raw emails are only shown as *citations*.
2.  **Context is King**: Users define "Context Windows" (e.g., "Project Alpha", "Key Clients"). The AI adapts its answers based on the active context.
3.  **Generative UI**: The interface adapts to the content. A query about dates renders a Calendar widget; a query about relationships renders a Network Graph.

---

## 2. The "Fluid Canvas" Layout

We abandon the traditional 3-column email layout for a **Source-Canvas-Detail** architecture.

### 2.1 Zone A: The Context Manager (Left Sidebar)
*   **Concept**: Instead of static folders, this area manages **"Attention Scopes"**.
*   **Sections**:
    *   **Auto-Clusters**: "Priority Inbox", "Needs Reply", "Follow-ups".
    *   **Smart Contexts** (User Defined): "Q4 Fundraising", "Legal Issues", "Person: John Doe".
    *   **Source Data**: Attached documents (PDF/Word), Spreadsheets extracted from emails.
*   **Interaction**:
    *   Clicking a Context *filters* the RAG engine for the main view.
    *   Checkbox selection allows multi-context synthesis ("Compare *Q4 Budget* with *Q3 Spend*").

### 2.2 Zone B: The Intelligence Canvas (Center Stage)
This is the main workspace. It has two distinct modes:

#### Mode 1: The Briefing (Passive / "Morning Coffee" Mode)
*   **Analogy**: A personalized newspaper or Notion dashboard.
*   **Content**:
    *   **The Lead Story**: "You have 3 critical decisions today." (Summary of top risks/requests).
    *   **Deal Watch**: "Acme Corp contract is stalling." (Dealmaker insight).
    *   **Pulse**: A mini-heatmap of today's communication volume.
*   **Actionability**: Each card has quick actions (`[Draft Reply]`, `[Snooze]`, `[Delegate]`).

#### Mode 2: The Studio (Active / "Deep Work" Mode)
*   **Analogy**: NotebookLM / ChatGPT / Perplexity.
*   **Trigger**: Clicking a Briefing card or typing in the Omni-Bar.
*   **Features**:
    *   **Streaming Q&A**: Answers appear in real-time.
    *   **Rich Citations**: Every claim has a `[1]` footnote. Hovering shows a snippet; clicking opens the source.
    *   **Embedded Widgets**:
        *   "Show me the timeline" -> Renders a Timeline Component *inside* the chat.
        *   "Who is involved?" -> Renders a Contact Card.

### 2.3 Zone C: The Grounding Panel (Right Slide-over / Overlay)
*   **Purpose**: Truth verification.
*   **Content**: The raw email thread or document.
*   **Behavior**:
    *   Hidden by default.
    *   Slides in when a citation `[x]` is clicked.
    *   Highlights the exact paragraph referenced by the AI.

---

## 3. Generative UI Components (The "Widgets")

The Chat interface won't just output text. It will render React Components based on the **Intent Classification** from the backend.

| Intent | UI Component | Description |
| :--- | :--- | :--- |
| `intent:scheduling` | `<CalendarWidget />` | Interactive calendar slot picker. |
| `intent:relationship` | `<NetworkGraph />` | (Existing) Force-directed graph of connections. |
| `intent:finance` | `<DataTabel />` | Extracted numbers in a clean table. |
| `intent:decision` | `<ApprovalCard />` | Big "Approve / Reject" buttons with risk analysis. |
| `intent:draft` | `<Editor />` | A rich-text editor pre-filled with the draft. |

---

## 4. User Experience Flows

### 4.1 The Executive "Morning Flow"
1.  **Open EchoMind**: Sees "The Briefing".
2.  **Scan**: Reads "The Lead Story" (3 critical emails).
3.  **Act**: Clicks `[Approve]` on Item 1. Clicks `[Draft Reply]` on Item 2.
4.  **Deep Dive**: Item 3 is complex ("Project Delay").
5.  **Transition**: Clicks Item 3. UI shifts to "Studio Mode".
6.  **Query**: "Why is it delayed?" -> AI answers citing 5 emails from the engineering lead.
7.  **Resolve**: Types "Draft an email to Engineering asking for a recovery plan."

### 4.2 The Dealmaker "Hunter Flow"
1.  **Select Context**: Checks "Active Deals" in the Left Sidebar.
2.  **Prompt**: "Who hasn't replied in 7 days?"
3.  **Result**: AI lists 3 people. Renders a `<FollowUpList />` widget.
4.  **Action**: Clicks `[Nudge All]`. AI generates 3 personalized follow-up drafts.

---

## 5. Technical Requirements for v0.6.0 (Preparation)

To enable this UI in v0.7.0, the v0.6.0 Backend/RAG must support:

1.  **Citation Metadata**: Search results **must** return exact `chunk_id` and `text_offset` to allow the UI to highlight the source.
2.  **Structured Threads**: The API must be able to return a full conversation thread structure to populate the "Grounding Panel".
3.  **Streaming JSON**: The API should support streaming structured data (e.g., `text` chunks followed by `widget_data` JSON) to render UI components on the fly.
4.  **Dynamic Filters**: The Vector Store query params must support complex filtering by `sender`, `date`, `thread_id` to support the "Context Manager".

---

## 6. Migration Strategy (The "Burn the Ships" Approach)

*   **Direct Replacement**: v0.7.0 will **overwrite** the existing Dashboard. There will be no toggle to switch back.
*   **Onboarding is Critical**: Since the paradigm shift is massive, the first login must trigger a high-quality "Neural Interface" tutorial (e.g., "Here is your Briefing", "Ask your first question").
*   **Legacy Code Removal**: Delete the old `dashboard/` page components immediately to prevent technical debt accumulation.
