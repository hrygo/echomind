# Changelog

All notable changes to this project will be documented in this file.

## [v0.9.3] - 2025-11-23

### Added
- **Smart Copilot (Omni-Bar)**: Unified Header Search and AI Chat into a single, context-aware input.
  - Seamless mode switching between "Instant Search" and "AI Chat".
  - Auto-detection of chat prompts (e.g., questions) to activate chat mode.
- **RAG Integration**: Chat service now prioritizes explicit context (from search results) for AI responses.
- **Real-time Streaming**: Implemented Server-Sent Events (SSE) for real-time AI chat responses.

### Changed
- **UI/UX**: Replaced dedicated Search/Chat components with the integrated `CopilotWidget`.
- **Backend API**: Updated chat endpoint (`/api/v1/chat/completions`) to accept `context_ref_ids`.
- **Dependencies**: Refactored `ChatService` to use interfaces for `ContextSearcher` and `EmailRetriever`.

### Fixed
- **Authentication**: Resolved missing JWT token in frontend `fetch` requests for search and chat.
- **Frontend Routing**: Fixed 404 error by adding rewrite rules in `next.config.ts` to proxy API calls to the backend.
- **Frontend Display**: Corrected "Search is not defined" error by re-importing the `Search` icon.
- **Linting**: Addressed all remaining frontend lint warnings (unused variables, imports) and backend unreachable code.


## [v0.9.0] - 2025-11-23

### Added
- **Task Engine (Phase 6.2)**:
  - **Task System**: Full backend support for Tasks (Model, Service, API) with status and priority management.
  - **Dashboard Integration**: New `TaskWidget` on the Manager Dashboard to view and manage active tasks.
  - **Smart Action Integration**: "Create Task" button in email details now directly creates tasks in the system.
  - **Optimistic UI**: Instant feedback when toggling task status.

### Changed
- **Database**: Added `tasks` table.

## [v0.8.0] - 2025-11-23

### Added
- **Smart Actions**: AI now extracts actionable items (Calendar Events, Tasks) from emails and presents them as interactive buttons in the email detail view.
- **AI Chat Copilot**: A contextual AI assistant available via a right-side drawer, supporting streaming responses and Markdown rendering.
- **Responsive Design**: Complete mobile adaptation.
  - **Header**: "Mobile Collapse, Desktop Expand" search bar strategy.
  - **Sidebar**: Implemented as a swipeable Bottom Sheet/Drawer on mobile.
- **Unified UI**: Standardized `Input` components across the entire application for consistent look and feel.
- **Internationalization**: Full i18n support for the new Chat interface, Insights page, and Smart Actions (English/Chinese).

### Changed
- **Header UX**: Removed the standalone Search Filter component in favor of a streamlined search experience.
- **Visuals**: Enhanced transparency effects in UI overlays and optimized font readability in Settings.

### Fixed
- **Stability**: Fixed database connection issues related to configuration mismatches.
- **Bug**: Resolved SSE stream parsing errors in the Chat interface.
- **Accessibility**: Fixed Radix UI accessibility warnings in mobile components.

## [v0.7.2] - 2025-11-22

### Added
- **Organization & Team Models**: Implemented backend models and migration logic for multi-tenancy (Organizations, Teams, Members).
- **Organization API**: Added CRUD endpoints for organizations (`POST /orgs`, `GET /orgs`).
- **Frontend UI**: Added Organization Switcher and Create Organization Modal.
- **State Management**: Implemented `useOrganizationStore` with Zustand for managing organization context.

### Fixed
- **Frontend Build**: Resolved multiple build errors related to UI components (`Dialog`, `DropdownMenu`) and API imports.
- **Type Safety**: Fixed TypeScript errors in `Zustand` persist configuration and component props.

## [v0.7.0-beta] - 2025-11-22

### Added
- **Semantic Search (RAG)**: Full-text search powered by OpenAI embeddings and pgvector for semantic understanding.
- **Search Filters**: Added support for filtering by Sender, Start Date, and End Date in the Search API.
- **Search History**: Frontend implementation of recent search history with suggestions.
- **Relevance Scoring**: Search results now include a relevance score (0-1) to indicate match quality.
- **Performance**: Optimized search latency to <500ms for standard queries.
- **Team Collaboration Design**: Completed architectural design for Phase 6 (Organizations, Teams).

### Changed
- **API**: Updated `/api/v1/search` to accept `sender`, `start_date`, and `end_date` parameters.
- **Backend**: Integrated `pgvector-go` for efficient vector similarity search.
- **UX**: Improved empty states and loading indicators for search results.

## [v0.5.3] - 2025-11-22

### Fixed
- **Network Graph Rendering**: Resolved blurry canvas issues on high-DPI screens by implementing `ResizeObserver` for dynamic sizing.
- **API Stability**: Fixed a 500 Internal Server Error in the Network Graph API caused by incorrect User ID type assertion.
- **Localization**: Completed full Chinese translation for all Dashboard widgets (Smart Feed, Action Center, Intent Radar).

### Changed
- **Visual Identity**: Simplified the app slogan to "Turn noise into decisions" and increased the logo size in the sidebar for better brand presence.
- **Layout Optimization**: Improved the Insights page layout to utilize full vertical screen space, removing unnecessary padding and max-width constraints.

## [v0.5.2] - 2025-11-22

### Added
- **AI Morning Briefing**: A dynamic header in the dashboard that greets the user and provides a high-level summary of risks and tasks.
- **Smart Feed**: A new component in the Executive View that highlights high-priority emails with AI-generated summaries, risk levels, and suggested actions.
- **Action Center**: Enhanced Manager View with an interactive task list (checkbox support) and filtering by priority.
- **Intent Radar**: Visual radar chart in Dealmaker View to show the distribution of email intents (Buying, Partnership, etc.).
- **Opportunity List**: A list of high-confidence opportunities in Dealmaker View.
- **Smart Follow-up**: Improved visual design for tracking pending replies.

### Changed
- **Sidebar Navigation**: Restructured into "Intelligence", "Mailbox", and "Labels" for better logical grouping.
- **Dashboard Layout**: Unified view switcher and improved overall aesthetic with a modern, clean design.
- **Project Identity**: Updated branding to "EchoMind" with the slogan "Your personal Chief of Staff for email. Turn noise into decisions".

### Fixed
- Removed redundant user profile section from the sidebar footer.
- Fixed various linting issues in backend and frontend code.
