# Daily Premium Tutor Overflow Fixes Design

## Goal

Fix the reported web UX regressions: daily bonus state, Premium gating for listening/pronunciation, noisy AI Tutor history text, desktop/mobile overflow, and repeated AI Tutor lessons.

## Observed Causes

- Daily bonus UI derives availability from local habit state and marks a day as claimed even when the server returns `claimed=false`. The home block also uses generic "Today" text in multiple places, so the action can look like a status label instead of a claim action.
- AI Tutor is Premium-gated, but standalone listening/shadowing and pronunciation entry points can still be shown as normal actions in Free UI copy.
- AI Tutor completed lesson cards use raw topic/lesson goal fields. English system-style goals such as "To talk about sports activities and routines..." can appear as user-facing history text.
- The AI Tutor panel and feedback text include long unbroken/generated strings. Several containers need `min-width: 0`, wrapping, and horizontal overflow guards on desktop and mobile.
- AI Tutor lesson selection avoids lessons that have a session row, but the local generated lesson ID still comes from `WordLessonCount + LessonCount`. Starting/restarting before completion can reuse the same deterministic sequence.

## Considered Approaches

1. **Frontend-only patch.** Fastest for labels and overflow, but it cannot reliably prevent repeated lesson prompts or invalid daily claims.
2. **Backend-only patch.** Strong for repeated lessons and claim rules, but it would leave confusing UI labels and overflow symptoms.
3. **Targeted full-stack fix.** Add focused backend guards/tests for eligibility and AI Tutor sequence, plus frontend copy/layout guards.

Use approach 3. It keeps the change small while fixing the root causes instead of hiding symptoms.

## Scope

- Daily bonus claim button is enabled only when the daily goal is complete and the 24-hour claim window is open.
- If the server refuses a daily claim, the client must not mark today's habit as claimed or complete.
- Daily bonus copy should read as status plus a clear action: "Daily streak", "Today: done/in progress", and "Claim XP bonus" only when actionable.
- Listening/shadowing and pronunciation are Premium-only in the web UI and plan copy. Free cards route to Premium/paywall instead of starting costly audio/pronunciation flows.
- Completed lesson history should display a clean lesson title/theme, not raw lesson goals or prompt/system text.
- AI Tutor generated feedback, questions, choices, and composer content must wrap within desktop and mobile viewports.
- New AI Tutor lesson creation should advance to an unused generated lesson sequence when the current deterministic sequence already has a user session.

## Backend Design

Add a store helper that counts AI Tutor sessions for a user in the current language/interface/level band. The AI Tutor engine uses that count as an offset when no banked approved lesson is available. This turns local generated lesson IDs into a monotonically advancing per-user sequence even if `LessonCount` has not changed yet.

The existing `ai_tutor_sessions` table remains the source of "used prompts"; no new table is needed. The helper joins sessions to lessons so it can scope by learning language, interface language, and level band.

Daily bonus server behavior already enforces the 24-hour rule. The required backend test coverage should verify that `/api/daily/claim` returns `claimed=false` without awarding XP inside the lock window.

## Frontend Design

Daily bonus:

- Compute `canClaimBonus` from current habit completion, `dailyBonusLocked`, and `busy`.
- Render a disabled button with a status-specific label when blocked.
- Do not mutate local habit state to `claimed=true` when the API returns `claimed=false`; instead refresh the session or leave the current state as-is and show the lock message.

Premium gating:

- Add a reusable `requiresPremium(view/action)` check for `shadowing` and `pronunciation` entry points.
- Update home route steps, quest actions, learning lab buttons, nav activation, and plan feature copy so Free users see Premium routing for listening/pronunciation.

History cleanup:

- Add `cleanTutorHistoryLabel` that rejects system-like goals (`To talk about...`, `Goal:`, `Topic/theme seed`, prompt labels) and prefers lesson title, theme, level, or localized "AI Tutor".
- Use it for completed lesson cards and local completed lesson writes.

Overflow:

- Add targeted CSS guards for `.tutor-workspace`, `.tutor-session-v2`, `.tutor-context-v2`, `.tutor-message-v2`, `.tutor-feedback-v2`, `.tutor-task-copy-v2`, `.tutor-srs`, completed lesson dialog cards, and composer textareas.
- Use `overflow-wrap: anywhere`, `word-break: break-word`, `min-width: 0`, `max-width: 100%`, and document-level `overflow-x: hidden` only as a last guard.

## Testing

- Go tests:
  - AI Tutor start creates a different lesson/session when the same user starts twice without completing.
  - Daily claim refused inside the 24-hour lock does not award a second XP bonus.
- Frontend/source tests:
  - React source keeps Premium-only copy for listening/pronunciation.
  - React source contains history cleanup guards for "To talk about" style labels.
- Build:
  - `npm run test -- --runInBand ...` for focused Go tests.
  - `npm --prefix web-react run build`.
  - `npm run check` if the focused checks pass.
- Browser:
  - Use Playwright/browser if available. Current MCP Chrome launch is blocked by system policy, so fallback is local build plus source/CSS assertions unless the browser becomes available.
