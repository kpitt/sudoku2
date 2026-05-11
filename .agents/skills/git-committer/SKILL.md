---
name: git-committer
description: Controls the formatting of Git commit messages. Use this skill whenever creating, editing, or proposing a Git commit message to ensure it follows the project's specific style guidelines (max 50 char subject, no semantic prefixes, present tense body, and mandatory footer).
---

# Git Committer

This skill enforces a specific, high-signal Git commit message format. It prioritizes clarity, consistency, and a clean history by avoiding semantic prefixes and requiring a specific footer.

## Guidelines

### 1. Subject Line

- **Length:** Maximum **50 characters**.
- **Capitalization:** Start with a capital letter.
- **Punctuation:** Do not end with a period.
- **Content:** Be concise but descriptive.
- **Constraints:** **DO NOT** use semantic prefixes (e.g., `feat:`, `fix:`, `chore:`, `refactor:`).

### 2. Message Body

- **Tense:** Always use the **present tense** (e.g., "Fix", "Add", "Change", "Implements", "Resolves").
- **Phrasing:** Avoid filler phrases like "This update", "This change", or "This commit". Start directly with the action (e.g., "Resolves an off-by-one error" instead of "This change resolves...").
- **Wrapping:** Wrap lines at **72 characters**.
- **Content:** Focus on the "why" and the "impact" of the change rather than just what was changed. Use a blank line between the subject and the body.

### 3. Mandatory Footer

- **Format:** Every commit message MUST end with an "Assisted-by: [agent]" footer identifying the AI coding agent responsible for the commit. The "[agent]" string must follow the standard `Name <email>` structure. Examples for common coding agents:
  - **Gemini CLI:** `Assisted-by: Gemini CLI Agent <google-cli-agent@google.com>`
  - **Claude Code:** `Assisted-by: Claude Code <noreply@anthropic.com>`
  - **GitHub Copilot:** `Assisted-by: GitHub Copilot <copilot@github.com>`
- **Placement:** Separate the footer from the body (or subject, if no body) with a blank line.

## Examples

For detailed examples of good and bad commit messages, see [references/message-examples.md](references/message-examples.md).
