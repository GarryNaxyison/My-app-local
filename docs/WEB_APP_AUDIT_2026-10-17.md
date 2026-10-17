# Web app audit: 2026-10-17

## Scope and evidence

Reviewed the AI Tutor start path, the shared lesson generator/validator, the React tutor and app shell structure, existing responsive audit notes, the existing Playwright suite, the user's 29 browser captures, and the attached bug-report export dated September 7-8, 2026.

## Confirmed AI Tutor issue

`aiTutorEngine.Start` generated one lesson and ran one preflight quality check. A transient AI error, malformed JSON, local validation failure, or rejected preflight immediately failed `/api/ai-tutor/start` with HTTP 500. No retry occurred. The learner's learned-vocabulary count is not a generation prerequisite: every tutor lesson generates its own six words and six learning/recall steps. A three-word account should be allowed to start.

The engine now retries generation and preflight up to three times and varies the topic rotation per attempt. Backend tests cover recovery after an invalid draft and exhaustion without creating a session.

The accepted-answer feedback misattribution is also fixed in the web tutor: feedback is stored against the stage that was submitted and appears when reviewing that completed stage, rather than being shown as feedback for the newly opened question. Desktop and mobile Playwright regression coverage passes for this behavior.

## Visual findings from supplied captures

### Mobile, Android browser

- Large text/display scaling causes headings and control labels to consume far more space than the current default mobile layout expects.
- Several screens retain a large outer frame, an inset panel, and another bordered content card, wasting width and making text feel buried in nested windows.
- Learn Words and conversation results use constrained inner scroll areas; content/audio controls are cut off while the outer screen has substantial unused space.
- Practice/roleplay/translation composers show oversized textareas and send buttons; after submission the button, language controls, attachments, and save-to-notes row can overlap or escape their panel.
- Some pages (notably words/practice/tools) leave large blank vertical regions instead of using the available screen height for the active task.
- Audio sample widgets occupy too much area when the user needs only playback controls.
- The bottom navigation works as a horizontally scrollable rail, but its visible edges crop neighboring items; it needs clearer active-item context and spacing under narrow widths.

### Desktop

- The desktop shell uses a broad top identity/header strip plus an icon ribbon with many functions. The app content starts below these stacked bars, leaving limited vertical space at laptop heights.
- The tutor's 21-step sidebar remains a full-height column while the active task is comparatively simple; a compact progress navigator would give the lesson more width.
- Content panels have layered rounded borders and shadows that create a nested-window effect even where no separate interaction layer exists.

## Functional reports to reproduce and verify

The attached reports describe behavior beyond responsive styling:

- AI Tutor rejected short answers but advanced, with the feedback appearing against the next question.
- Vocabulary answer choices were all graded incorrectly for at least one round.
- Roleplay repeated a correction for a later unrelated answer.
- A lesson appeared complete after one sentence.
- Pronunciation advice mixed English and Russian; mistake-word playback failed.
- The dashboard calendar used English month/day labels and Sunday-first ordering.
- Adding an example to Notes showed success feedback without an observable saved item.
- Image upload and responsive composers collided with send controls in mobile translation flows.

These reports should become reproducible backend/UI regression tests rather than being treated as solved by layout changes.

## Existing automated coverage

- `npm --prefix web-react run build`: passed.
- `npm --prefix web-react run e2e`: 174 passed, 50 skipped by platform-specific conditions; no failures. Existing tests include default mobile flows for tutor, tools, audio, words, and scroll/overlap regressions.
- `go test ./...`: passed.

The default Playwright device emulation does not reproduce the supplied Android text scaling/browser state. Real-device or browser zoom/system-font checks remain necessary.

## Redesign priorities

1. Fix the reported tutor stage/feedback synchronization and validate generated answer choices; preserve the complete lesson stage count.
2. Establish a compact mobile shell: one content surface per task, no nested scrolling unless the content itself is a history/list, safe-area-aware header/actions, and a bottom navigation that does not crop the active task.
3. Reflow composers so textarea, send, voice, upload, and save controls have separate rows at narrow widths and after result rendering.
4. Remove fixed heights/min-heights that create blank space or clip dynamic text; size audio playback as a compact control.
5. Reduce the desktop tutor roadmap to a compact progress control and give the active lesson the workspace width.
6. Add Playwright coverage for 320/360/390px, tablet, laptop height, enlarged text/zoom, translated long labels, and post-submit states; retain all existing learning routes.

## Changes implemented in this pass

- AI Tutor lesson generation now retries generation, validation, and preflight failures; checker errors and malformed output cannot incorrectly accept an answer.
- Learned vocabulary IDs accept supported legacy aliases, and vocabulary/report state is stored and read through SQLite.
- New-word selection uses only validated AI-generated SQLite cards. The web flow requests a new card at the user's CEFR level when the shared pool is exhausted, and duplicate suggestions are rejected and retried.
- The startup migration removes old imported dictionary rows while retaining accepted AI cards. A separate one-time migration resets all users' learning progress and sessions, setting their level to A1, while preserving accounts, credentials, and profile preferences.
- Desktop and mobile layouts were adjusted to reduce nested surfaces, prevent text/audio/composer overlap, and keep send, upload, voice, and note controls in separate rows. Phrasebook success now follows server confirmation; audio playback falls back to browser speech where available.
- The calendar is localized and Monday-first. Pronunciation guidance is localized, and long reports retain enough detail on narrow screens.
- Text generation defaults now use `openai/gpt-6-luna` for chat, translation, vocabulary, and text-only pronunciation feedback. Speech recognition, audio analysis, and synthesis retain their specialist models.
- Corrected the AI Tutor loading status: a lesson is generated by AI and does not require prior vocabulary study. Newly generated and reused approved lessons seed single-word entries into shared SQLite; multiword chunks are excluded, duplicates are ignored, and storage errors do not prevent lesson startup.\n- Pronunciation audio assessment now defaults to `google/gemini-3-flash-preview`, which supports audio and has lower listed input/output token rates than Gemini 2.5 Pro. Perceived pronunciation quality still needs real-user sample evaluation after rollout.

## Verification and remaining limits

- `npm run build` in `web-react`: passed after removing a stale intermediate bundle left by an interrupted build.
- Full Playwright suite: 177 passed, 51 skipped by desktop/mobile test conditions, 0 failed.
- Full `go test ./...`: passed (178.180s); the targeted empty-vocabulary onboarding regression also passed.
- The supplied Android captures show enlarged browser/system text. Playwright covers narrow/mobile and zoom-equivalent viewports, but a real Android browser/device check remains useful.
- Deployed through the permanent `.online` upload service. Public `/healthz`, `/app/`, `/neriva.html`, and upload `/healthz` returned HTTP 200. Main, transition, and upload services are active; post-deploy load average was 0.05, the shared SQLite vocabulary database is active, and the old JSON vocabulary directory is absent. The one-time progress-reset marker is present, all 57 user profiles and all 15 web accounts remain, learned-word and mistake rows are empty; small XP/habit counts reflect activity recorded after the reset.
