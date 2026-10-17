# Changelog

Project change history is maintained in [`docs/tracking/CHANGELOG.md`](docs/tracking/CHANGELOG.md), alongside the user request log and release audit.

## 2026-10-17

- Clarified AI Tutor preparation status, added an elapsed timer that survives page refresh, and extended the lesson-start timeout to three minutes.
- Enabled AI Tutor onboarding with no pre-learned words; successful lessons seed deduplicated single-word entries into the shared SQLite dictionary.
- Replaced misleading dictionary-based lesson loading copy and improved responsive desktop/mobile layouts across learning surfaces.
- Switched text generation to OpenRouter `openai/gpt-6-luna` and pronunciation audio assessment to the audio-capable, lower-token-rate `google/gemini-3-flash-preview` model.
- Reset all users' learning progress once while preserving all 57 user profiles and the 15 web accounts with their login/profile settings; deployed and verified service health and low host load.
