# Mobile Tutor And Premium Gating Design

## Goal

Fix the mobile AI Tutor flow so the current task comes before progress, prevent mobile horizontal drift, and make the Premium screen push users toward paid tiers without exposing high-cost AI Tutor usage to Free users.

## Requirements

- On mobile widths, the AI Tutor task/context block appears before the progress block.
- The progress block remains available directly after the current task and still supports reviewing completed steps.
- Mobile pages must not create document-level horizontal overflow from tutor, premium cards, or long button/card text.
- AI Tutor is presented as a Premium feature. Free remains useful but does not advertise unlimited costly AI work.
- Plan copy should create a clear upgrade ladder:
  - Free: daily habit, starter lessons, basic word training, saved phrasebook, progress overview, no AI Tutor.
  - Premium: AI Tutor, guided AI lessons, voice/pronunciation checks, photo/translator tools, normal daily limits.
  - Platinum: higher limits, priority/intensive practice, roleplay/pronunciation depth, best value for heavy learners.

## Design

Keep the existing React/Vite app structure and avoid backend schema changes. `PremiumView` will derive a local feature list from plan product/title, render it inside each plan card, and mark unavailable Free capabilities as locked text. This gives the product ladder without depending on new API fields.

For AI Tutor layout, keep desktop unchanged: progress stays as a sticky left rail and task stays to the right. At the mobile breakpoint, use CSS ordering so `.tutor-context-v2` renders first and `.tutor-plan-v2` second. This avoids duplicating the progress component and keeps completed-step navigation intact.

For overflow, add targeted mobile guards: containers use `min-width: 0`, `max-width: 100%`, `overflow-wrap: anywhere`, and mobile tutor/premium grids avoid fixed/min-content widths. The document should keep `scrollWidth <= innerWidth` for the AI Tutor mobile view.

## Testing

Add Playwright regression coverage in `web-react/e2e/web-smoke.spec.ts`:

- Mobile AI Tutor: verify `.tutor-context-v2` is visually above `.tutor-plan-v2` and the document has no horizontal overflow.
- Premium: verify plan cards describe AI Tutor as excluded from Free and included in Premium/Platinum.

Run the targeted Playwright tests first to see them fail, then implement, then rerun targeted tests and the web build.
