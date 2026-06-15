# Web Report, AI Tutor History, and Admin Fix Copy Design

## Status

Approved approach: variant 1, the minimal targeted fix.

## Goal

Fix three small user-facing regressions without redesigning the app shell:

- keep the generic bug report entry point visible on web desktop;
- remove lesson-goal/prompt text from AI Tutor history messages so the history starts with the story itself;
- make the Telegram admin correction prompt explicit about word order.

## Problems

The web app already has a generic bug report dialog and button. Mobile has a separate fixed quick-controls button, but desktop relies on the topbar action row. The topbar action row can run out of horizontal space and clip the right-side controls, so desktop users may not see the report button.

AI Tutor start/restart currently writes a panel history message from `lesson_goal || step.instruction` plus `story.text_target`. That can put prompt-like text such as "Learn vocabulary related to sports..." before the story in history.

The Telegram admin Fix prompt says `word - translation`. The parser already treats the left side as the corrected target-language word and the right side as the interface-language translation, but the copy is ambiguous for admins.

## Scope

In scope:

- Adjust desktop topbar/report-button layout so the bug report button remains visible on desktop widths.
- Keep the existing mobile report button behavior unchanged.
- Change AI Tutor start and restart history messages to use the story text as the history body, not the lesson goal.
- Keep the lesson goal visible in the AI Tutor hero/context where it already belongs.
- Update Telegram admin Fix prompt and validation fallback to say: first the corrected word/phrase in the learning language, then the translation in the interface language.
- Add focused regression tests for these behaviors.

Out of scope:

- New bug-report backend behavior.
- New admin web dashboard.
- Heuristic parsing of reversed Telegram correction order.
- Full topbar redesign.
- Regenerating AI Tutor lesson content.

## Chosen Approach

Use the existing UI surfaces and make them robust.

For desktop bug reports, keep using the existing `.v2-report-button` in `TopBar`, but update CSS and tests so it cannot silently disappear under ordinary desktop width pressure. The likely fix is to let the action cluster shrink/wrap safely or hide the report label earlier while preserving the icon button and accessible label.

For AI Tutor history, introduce a small helper or inline expression that builds history body text from `lesson.story.text_target` first. It can fall back to a clean instruction only when no story exists, but it must not prepend `lesson_goal` to the story history message.

For Telegram Fix, keep the parser as-is. The input format remains positional:

```text
target-language word or phrase - interface-language translation
```

Example for English learning language and Russian interface:

```text
football match - футбольный матч
```

The implementation should use clear Russian admin copy in the bot prompt because the current ops workflow is Russian.

## Acceptance Criteria

- On desktop web, a user can see and open the generic bug report dialog.
- On mobile web, the existing quick report button still opens the same dialog.
- AI Tutor history no longer displays the lesson goal before the story body.
- Telegram admin Fix prompt clearly explains the left side is the learning-language word/phrase and the right side is the interface-language translation.
- Parser behavior stays unchanged: left side maps to `FinalWord`, right side maps to `FinalTranslation`.
- Existing report, AI Tutor, and vocabulary correction flows keep working.

## Tests

Add or update focused tests before implementation:

- Playwright desktop assertion that `.v2-report-button` is visible before opening the bug report dialog.
- Keep or extend the existing mobile assertion for `.mobile-report-button-v2`.
- Frontend/static regression test for AI Tutor history body construction so `lesson_goal` is not joined before `story.text_target`.
- Go unit/regression test for `parseAITutorWordReportFixText` confirming left side is word and right side is translation.
- Telegram prompt test confirming the Fix prompt copy describes the target-language word first and interface-language translation second.

## Risks

The only high-risk path is accepting reversed Telegram corrections. The parser cannot reliably detect languages from two arbitrary free-text fields, so supporting both orders would corrupt lesson payloads and SQLite vocabulary tables. The safe behavior is to keep one explicit format and make the admin prompt unambiguous.
