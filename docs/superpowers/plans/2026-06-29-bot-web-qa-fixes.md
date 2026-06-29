# Bot And Web QA Fixes Plan

Date: 2026-06-29

## Plan

1. Add focused Telegram regression tests proving old tracked audio messages are deleted after shadowing and pronunciation answers.
2. Call the existing audio cleanup helper from the missing Telegram answer handlers.
3. Update web audio/result rendering so example text and translation remain visible near audio controls.
4. Localize pronunciation report fallbacks, show multiple weak spots in compact report/map, and clean mojibake fallbacks in touched strings.
5. Reclassify saved notes by source, avoid translation/note duplication, and display manual vs section-generated notes separately.
6. Fix web chat layout sizing and message order for phone and desktop.
7. Harden navigation copy so password labels cannot appear on Next buttons.
8. Run focused Go tests and a React build, then commit and push only the files changed for this task.
