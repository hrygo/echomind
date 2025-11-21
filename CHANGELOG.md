# Changelog

All notable changes to this project will be documented in this file.

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
