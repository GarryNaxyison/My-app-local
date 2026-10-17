# AI Tutor preparation and next lesson

## Outcome

Make lesson preparation understandable during slow model responses and keep the next-lesson action usable after a completed lesson.

## Scope and constraints

- Keep the existing validated lesson generation and quality check flow.
- Give the browser request enough time for bounded model retries.
- Show that preparation is underway and a live elapsed timer; retain elapsed time across a page refresh in the same tab.
- Clear the timer when a lesson is returned or when a new preparation attempt is more than five minutes old.
- Preserve the completed lesson's explicit “start another lesson” action and verify it starts a distinct session.

## Acceptance criteria

1. While a lesson is preparing, the learner sees a clear preparing message and an elapsed `MM:SS` timer.
2. Refreshing the page does not reset the timer during the current preparation attempt.
3. The client does not abort a valid generation attempt at 30 seconds; it allows up to three minutes.
4. After completing a lesson, “Start another tutor lesson” starts a distinct lesson/session instead of returning the completed session.

## Verification

- Playwright covers the visible timer, refresh persistence, and next-lesson action.
- Go AI Tutor engine regression covers completion followed by a distinct generated lesson.
