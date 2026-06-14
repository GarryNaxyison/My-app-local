# AI Tutor Lesson UI Cleanup Design

## Status

Approved by the user on 2026-06-14 as one shared spec for the AI Tutor and Lesson UI cleanup.

## Goal

Make the AI Tutor and normal Lesson flows finish cleanly, show the word-report control where learners expect it, remove duplicated note/error blocks, and keep system UI text localized instead of leaking English or generic fallback labels.

## Scope

In scope:

- AI Tutor word-learning steps show a visible report button for the active `word_learn_N` stage.
- AI Tutor completion shows the completion text once and does not render stray fallback labels such as `Раздел`.
- AI Tutor system labels, option labels, buttons, and stage instructions use local copy keys for all supported app languages.
- AI Tutor note actions render as one compact horizontal, scrollable sentence/list instead of duplicated vertical sections.
- The AI Tutor `Ошибки` block renders only when there is real error content, not when only note actions exist.
- The normal Lesson chat ends after the user's answer: the chat remains readable, the text input stays inactive, and the CTA says `Следующий урок`.

Out of scope:

- Rebuilding the AI Tutor backend lesson generation model.
- Changing Telegram word-report moderation.
- Changing vocabulary storage or report approval logic.
- Redesigning the whole tutor layout.

## Current Problems

The frontend already has parts of the AI Tutor word-report UI, but the visible button is gated too narrowly and can disappear on actual `word_learn_N` stages. The completion view renders the same body copy through both the general instruction area and the completion summary, which duplicates text. The localization fallback layer can map unknown or fallback-like view text to generic words such as `Раздел`, which is not meaningful in the lesson-complete state.

The note/action extraction currently treats "save to notes" helpers as general error/note content in more than one place. This creates duplicated `Сохранить в заметки` blocks and even an `Ошибки: Сохранить в заметки` line when there is no real learner error to display.

In the regular Lesson view, after an answer is submitted the user should not continue writing in the same lesson. The next action is starting the next lesson, so the CTA must say that explicitly.

## User Experience

AI Tutor word cards:

- The learner sees a small report control on the active word card.
- The report control uses localized text, for Russian: `Сообщить об ошибке`.
- It appears for real `word_learn_N` stages even when the backend `kind` value is missing or stale.

AI Tutor notes and errors:

- Note actions are shown once as a single compact row/sentence with horizontal scroll when needed.
- Duplicate note blocks below the main card are removed.
- `Ошибки` appears only for actual error/explanation text.

AI Tutor completion:

- The completion screen shows the chain-complete message once.
- There is no unexplained `Раздел` label.
- Any XP/status line remains readable and localized.

Normal Lesson:

- After the learner submits an answer, the lesson is finished.
- Previous messages and AI feedback remain visible.
- The input is disabled/hidden for that completed lesson state.
- The main CTA starts another lesson and is labeled `Следующий урок`.

## Implementation Approach

Keep the change frontend-focused and use existing copy/helper patterns. Do not introduce a new design system or shared state layer. The fix should be scoped to `web-react/src/App.tsx`, `web-react/src/lib/i18n.ts`, styles needed for the compact note row, and existing source-level shell tests that validate generated web assets.

Stage detection should prefer explicit stage names for UI affordances:

```ts
const isWordLearnStage = /^word_learn_\d+$/.test(stage);
```

Known AI Tutor stages should render localized stage labels/instructions from `copy(...)` rather than backend English interface text. Generated learner content, such as target-language sentences, words, and questions, remains unchanged.

Note extraction should separate note candidates from error text before rendering:

- note/save helpers go to the compact notes row;
- real mistakes/explanations go to `Ошибки`;
- empty sections are not rendered.

## Testing

Add or update source-level frontend regression tests in the Go test suite for:

- word-report UI wiring and visibility based on `word_learn_N`;
- no duplicate completion body rendering;
- no generic `Раздел` fallback for tutor completion/view labels;
- Lesson completed CTA uses `Следующий урок`;
- AI Tutor note/error rendering keeps note actions out of the errors block.

Verification commands:

```powershell
go test ./... -run 'TestWebApp.*Tutor|TestWebApp.*Lesson'
npm --prefix web-react run build
```

After implementation, verify the local UI through Browser/Playwright when the app can be served locally.
