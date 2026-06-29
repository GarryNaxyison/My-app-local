# Bot And Web QA Fixes Design

Date: 2026-06-29

## Goals

- Delete previously sent Telegram audio examples once the learner submits a real answer.
- Keep audio examples useful by showing readable text next to the audio controls.
- Keep pronunciation coaching in the learner interface language and show more than one weak spot in the map.
- Separate manually saved notes from notes saved through learning sections, and avoid duplicating translation as note text.
- Improve web chat layout on mobile and desktop so answers are readable, start near the top, and use available space.
- Prevent password copy from leaking into navigation buttons.

## Decisions

- Reuse the existing Telegram `pronunciationMessages` tracking and `deletePreviousWordPronunciation` cleanup path instead of adding new message state.
- Treat notes saved from word/tools flows as section-generated notes, not manual notes.
- Keep React layout fixes in the existing chat, audio, phrasebook, and pronunciation components rather than introducing a new shell.
- Normalize common English pronunciation advice at render time as a fallback, while keeping the model prompt constraints in place.

## Non Goals

- No broad redesign of the app navigation.
- No database migration for notes; the source field remains text-compatible.
- No deployment in this spec unless explicitly requested after the fix.
