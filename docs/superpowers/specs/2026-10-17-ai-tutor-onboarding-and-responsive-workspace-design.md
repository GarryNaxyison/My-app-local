# AI Tutor Onboarding and Responsive Workspace Redesign

## Goal

Ensure a newly registered learner can start an AI Tutor lesson after learning only a few words, and redesign the authenticated web workspace so its desktop and mobile layouts remain clear, non-overlapping, and usable at supported screen sizes while preserving the product's learning flows.

## Scope

- Trace first-session setup, vocabulary state, lesson generation/validation, and user-facing errors.
- Make lesson generation resilient to transient model failures and invalid/rejected drafts.
- Audit the authenticated shell and high-density learning views at desktop, tablet, and phone widths; record concrete overlap, nesting, and overflow issues.
- Include the user's Android browser screenshots and attached bug reports as first-class reproduction evidence, including large text scaling and post-submit states.
- Simplify the tutor workspace hierarchy and responsive sizing without removing its lesson stages, vocabulary, navigation, or account functions.
- Add Playwright coverage for new-account tutor startup and responsive layout invariants, plus backend tests for generation retries.
- Record the user's UI preferences in `docs/tracking/USER_UI_WISHES.md` and progress in `docs/tracking/CHANGELOG.md`.

## Constraints and assumptions

- Three or four learned words are not a prerequisite for a complete tutor lesson: lessons contain their own six instructional words.
- Preserve the shared lesson engine and the existing web/Telegram parity contracts.
- Reuse the existing React/Vite, Go, CSS, and Playwright architecture; do not add a competing design system.
- Do not create backups in the repository or server.
- Production deployment is a separate operational step after local verification and review of the concrete result.

## Observable behavior

- A model timeout, malformed lesson, local validation failure, or preflight rejection triggers a bounded retry and does not immediately strand the learner on their first attempt.
- After retry exhaustion, the API returns a stable, understandable error and logs the reason without exposing raw provider internals to the UI.
- At 320px, 390px, tablet, and desktop widths, the app shell and tutor view have no horizontal page overflow, duplicate overlays, or overlapping primary controls.
- With Android large-text/display scaling, content cards expand to their text, composer attachments and submit controls remain distinct, and content remains reachable without clipped nested scroll regions.
- The tutor's current task, progress, lesson content, answer composer, and next action have a clear visual order on every viewport.

## Acceptance criteria

1. Generation retries transient and quality/validation failures up to a small fixed cap; successful later attempts create exactly one session and lesson.
2. A new user with an empty or short vocabulary can start a tutor lesson; the tutor generator remains the lesson word source.
3. Backend tests cover success-after-retry, exhausted retries, and no partial session persistence on failure.
4. Playwright covers tutor startup and layout invariants on desktop and phone, including long translated text and narrow viewport.
5. App shell/tutor redesign is checked at desktop, tablet, and phone widths, with no page-level horizontal scrolling or visible overlap.
6. Representative mobile views are checked with enlarged system text and post-submit content, not only default Playwright device emulation.
7. The audit records tested viewports, observed findings, changed areas, and any remaining unverified surfaces.

## Verification

- `go test ./...`
- `npm --prefix web-react run build`
- `npm --prefix web-react run e2e`
- Manual viewport review at 1440x900, 1024x768, 390x844, and 320x700.

## Audit notes

- First-start currently performs one generation call and one preflight call; any error or rejected draft returns directly to `/api/ai-tutor/start` as HTTP 500.
- Lesson validation requires six generated lesson words and derives the six learn/recall tasks when omitted. The learner's existing vocabulary count is not used as a start gate.
- The tutor is rendered inside a large single-file React application and extensive global CSS; existing project audits already identify a monolithic shell and duplicated responsive rules, so the redesign must test the integrated shell rather than isolated cards only.
- User-provided captures show overlapping top controls on the dashboard, clipped/nested output cards in words and conversation tools, submit buttons colliding with attachment controls, unused vertical space, oversized audio sample blocks, and a tutor timeline that consumes most of the desktop content width.
- The attached September 7-8 reports additionally describe tutor feedback appearing after navigation to the next question, incorrect word-choice grading, repeated roleplay corrections, a one-response lesson ending unexpectedly, inconsistent advice language, broken audio playback in mistakes, incorrect calendar locale/week start, and save-example feedback without a saved item.
- Full baseline e2e passes at default emulated viewports, so it does not reproduce the supplied Android browser captures; acceptance must add large-text/zoom and after-submit checks.
