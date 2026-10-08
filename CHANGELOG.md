# Changelog

## 2026-10-09

- Added a Windows browser-window dashboard for the `codex-orchestrator` profile, with live Codex JSONL events, configured agent/provider cards, task submission, and task interruption.
- Added unit and Playwright end-to-end coverage for the dashboard’s core flows.
- Switched the dashboard runner to the supported `codex exec --profile … --json` stream and updated the profile for the installed CLI's named-agent schema.
