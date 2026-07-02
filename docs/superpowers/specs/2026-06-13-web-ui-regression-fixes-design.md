# Web UI Regression Fixes Design

## Goal

Restore the NERIVA web UI after the recent layout regressions so desktop and mobile modes keep text readable, controls visible, and phrasebook saves obvious.

## Scope

- Roleplay session after scenario selection uses one wide dialogue surface like practice, with the answer composer at the bottom and no unused side column.
- Desktop function ribbon arrows stay visible inside the menu frame and do not overlap or push function buttons off-screen.
- Pronunciation on mobile flows vertically: target text, audio, recording controls, check/result, next phrase. Blocks must not overlap.
- Quick phrase saves update the bookmark visual state immediately after save and remain available across tutor, practice, word results, roleplay, and tools where a saveable phrase exists.
- Translation tool mode buttons are compact and input/output areas use the available text space.
- Roleplay accepts a voice answer using the existing upload/STT path; the transcript becomes the learner answer and the roleplay continues.

## Non-Goals

- No backend schema changes.
- No redesign of the product navigation model.
- No new landing or marketing pages.

## Acceptance Checks

- Playwright verifies desktop roleplay dialogue width, desktop ribbon arrow containment, mobile pronunciation non-overlap, compact tools buttons, quick-save active bookmark color, and roleplay voice submission.
- `npm --prefix web-react run build` succeeds.
- Visual browser check confirms desktop and mobile critical screens remain readable.
