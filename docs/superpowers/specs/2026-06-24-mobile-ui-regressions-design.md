# Mobile UI Regressions Design

## Goal

Fix the reported mobile web app regressions without changing the app's visual identity: keep the current dark, compact learning UI, but make text, audio controls, lesson scroll, level reentry, mistakes, and tools usable on narrow screens.

## Requirements

- Free/Premium messaging must wrap inside cards and must not use misleading copy that says free users cannot request phrase audio. Free users may request generated phrase audio; pronunciation/listening checks remain Premium.
- Starting or reopening a lesson must show the beginning of the lesson content. Chat output must not auto-scroll to the bottom when a lesson or tab is mounted.
- Audio buttons and upload controls must stay inside their cards on mobile, including long error labels and file names.
- The level tab must remain usable after leaving and returning. Returning to the level tab should not leave an empty or broken state.
- Mistakes must show an unobtrusive hint in the header: "Нажмите на ошибку, чтобы проработать её"; the header must stay compact.
- Users must be able to delete one mistake without clearing all mistakes.
- Clicking a mistake remains the fastest path into practice, but the UI must expose that behavior with the hint and per-card controls.
- "Потренировать похожие" must start a different same-category mistake when one exists, not restart the same item.
- Tools must avoid overlapping textareas, language selectors, file controls, and the send button on mobile.

## Root Cause Summary

- `ChatComponent` scrolls to `scrollHeight` whenever messages mount, causing lessons to appear at the end.
- `MistakesView` only exposes "clear all"; per-card practice is hidden behind an unlabeled card click, and similar practice chooses the first item in the category.
- `ToolsView` uses an absolutely positioned submit button in a tight textarea shell and mobile grid rows that can overlap the language selector and attachments.
- Audio controls use minimum widths and nowrap labels that can overflow on 390px screens.
- Level activation always starts a fresh request, and the visible level panel has no recovery path if reentry lands without a question.

## Design

- Replace chat auto-scroll with top-on-mount behavior and only preserve manual scrolling inside the chat panel.
- Add `deleteMistake(index)` in React and a server route `/api/mistakes/delete` that removes a single mistake by stable index.
- Add per-card delete icon buttons with accessible labels and stop click propagation so deletion does not open practice.
- Track the active mistake index in React. `trainSimilar` will choose the next same-category mistake excluding the active index, falling back to the next visible mistake.
- Adjust mobile CSS so tools use normal document flow: textarea, send button, language row/file controls, and quick save stack without absolute overlap.
- Tighten audio and upload CSS: `min-width: 0`, wrapped labels, bounded waveforms, and full-width controls inside narrow cards.
- Keep `activateView("level")` from forcibly restarting a usable level session on every return; only start when there is no question/result.

## Test Strategy

- Add focused Playwright regression tests to `web-react/e2e/web-smoke.spec.ts` for lesson scroll, level reentry, mistakes hint/delete/similar, audio containment, and tools overlap.
- Add a Go API regression in `web_api_feature_test.go` for single mistake deletion.
- Run focused tests first, then build and broader checks.
