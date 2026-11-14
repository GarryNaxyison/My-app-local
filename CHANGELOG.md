# Changelog

## 2026-11-14

- Fixed stale AI Tutor web answers by sending and checking the expected stage; serialized session submissions within the process.
- Committed SQLite lesson completion, XP, counters and habit progress in one transaction. Added rollback/replay tests; kept JSON completion under its store lock.
- Made post-lesson quality assessment best effort with a five-second deadline after completion; included stored learner answers in its input.
- Released microphone tracks on navigation and late permission resolution, handled recorder setup errors, and used a localized native recording button.
- Preserved tutor drafts in session storage, suppressed duplicate actions, retained pending status for overlapping actions, reset explicit retry timing, and added an in-memory timer fallback.
- Prevented automatic PWA reload during input or pending operations. React build, full Go tests and 183 Playwright cases passed; 51 platform skips. Two extra PWA regression cases passed separately.

- Added a desktop/mobile web audit with code findings, UI recommendations, AI Tutor logic proposals, and explicit verification limits in `WEB_APP_AUDIT_2026-11-14.md`.
- Added the owner's current requests in root `USER_REQUESTS.md`. Installed three React/UI skills in the user Codex directory; application behavior was not changed.

Project change history is maintained in [`docs/tracking/CHANGELOG.md`](docs/tracking/CHANGELOG.md), alongside the user request log and release audit.

## 2026-10-17

- Clarified AI Tutor preparation status, added an elapsed timer that survives page refresh, and extended the lesson-start timeout to three minutes.
- Enabled AI Tutor onboarding with no pre-learned words; successful lessons seed deduplicated single-word entries into the shared SQLite dictionary.
- Replaced misleading dictionary-based lesson loading copy and improved responsive desktop/mobile layouts across learning surfaces.
- Switched text generation to OpenRouter `openai/gpt-6-luna` and pronunciation audio assessment to the audio-capable, lower-token-rate `google/gemini-3-flash-preview` model.
- Reset all users' learning progress once while preserving all 57 user profiles and the 15 web accounts with their login/profile settings; deployed and verified service health and low host load.
