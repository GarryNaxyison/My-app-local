# Bug Report Dialog Overflow Design

## Goal

The bug report dialog must stay usable on mobile when a user types a long report or attaches a screenshot with a long filename. The dialog must not create horizontal page overflow or shift off-screen.

## Root Cause

The screenshot upload control uses `width: fit-content`. A long unbroken filename can make the control wider than the dialog, and the dialog content then overflows the mobile viewport.

## Design

Keep the existing bug report flow and copy. Constrain the dialog and upload control to the viewport, let long text wrap or truncate inside the control, and keep the footer buttons within the dialog width.

## Acceptance Checks

- Playwright opens the bug report dialog on `/app/?view=home`, pastes an image with a very long filename, and verifies there is no horizontal viewport overflow.
- The long filename remains visible in the upload control without pushing the dialog outside the screen.
- `npm --prefix web-react run build` succeeds.
