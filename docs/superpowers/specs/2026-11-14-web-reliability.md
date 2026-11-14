# Web reliability fixes

Approved scope: resolve the defects documented in the 2026-11-14 web audit. Preserve the current lesson structure and navigation while fixing reliability.

Acceptance criteria:
- Web answers carry their expected stage; stale submissions cannot advance another stage. Session mutations serialize in-process; completion and reward commit atomically in SQLite and under the JSON store lock.
- An unavailable post-lesson quality model cannot prevent completion. Quality checking remains bounded by a short timeout and is best effort.
- Leaving voice input releases its microphone, including delayed permission resolution. Voice control is a localized native button.
- Duplicate UI operations are suppressed; finishing another operation cannot clear the pending indicator.
- Preparation timing distinguishes explicit retries from refresh and works without browser storage.
- PWA updates wait while there is unsaved input or an active operation. Tutor drafts survive reload.
- Addressed Go regressions and desktop/mobile Playwright pass. Document remaining limits.

Implementation sequence: failing server regressions; atomic completion and guarded submissions; browser regressions; input lifecycle, pending state and draft/update handling; broad checks; documentation and Git sync.

Verification: React production build passed. Full Go suite passed (201.995s); after the final completion/summary refinements, targeted completion and summary regressions passed. Full desktop/mobile Playwright passed 183 cases with 51 platform skips; two additional PWA update regressions passed separately. SQLite trigger injection proved reward failure rolls back session completion; replay grants one reward. Browser tests proved microphone cleanup and tutor draft restoration. Physical microphone, Safari and multiple server processes were not tested.

Compatibility: stage remains optional for older Telegram/API callers. Updated web requests carry it. Post-check is synchronous but limited to five seconds after atomic completion; outage cannot revoke XP or completed status. No persistent job queue or automatic failed-quality retry is introduced. PWA controller changes are ignored while input is dirty; the new version is used on the next deliberate navigation/reload. Draft persistence requires available sessionStorage; blocked storage retains input only while mounted.
