# App Platform Layout Stability Design

## Scope

This change applies only to the authenticated web application under `/app`.
It excludes the public landing site and its legal or SEO pages.

## Problem

After pronunciation audio is checked, the result report can be hidden by the
recording controls or be difficult to reach at some browser zoom levels. The
pronunciation layout has several competing responsive overrides, including
unscoped `!important` rules. The current placement also keeps the report in
the desktop tools column instead of presenting it as the next content block.

The web application already has independent desktop and mobile platform
selection through `data-platform`. Layout changes must not allow a desktop
rule to alter mobile rendering or vice versa.

## Chosen Design

Use an isolated adaptive pass for `/app`:

1. Keep the existing platform boundary: desktop rules apply only to the
   desktop platform and mobile rules apply only to the mobile platform.
2. Render a completed pronunciation report in normal document flow after the
   recording-control card. It must never be positioned over controls or
   hidden behind a fixed navigation surface.
3. Keep the desktop pronunciation workspace as a two-column work area while
   making the completed result a full-width next block below that workspace.
4. Keep the mobile workspace as a single vertical flow: target phrase,
   recording controls, then completed report. The pronunciation screen owns
   one vertical scroll region and reserves room for the mobile navigation.
5. Consolidate conflicting pronunciation responsive overrides into the final
   platform-scoped rules. Do not use generic viewport-only overrides that can
   accidentally make a zoomed desktop page behave as mobile.

## Component Boundaries

- `web-react/src/App.tsx` owns the pronunciation view hierarchy and moves
  only the completed report boundary needed to make it a post-workspace block.
- `web-react/src/styles/app.css` owns presentation. Platform-specific
  selectors are scoped to the pronunciation view and `data-platform` so they
  cannot affect other app screens.
- `web-react/e2e/web-smoke.spec.ts` owns browser-level regression coverage.

## Responsive Contract

Desktop and mobile are tested independently. For each platform, the
pronunciation result must be visually below its recording controls, fully
inside the viewport width, and reachable by vertical scrolling. No horizontal
overflow or rectangle intersection is permitted.

Coverage includes normal desktop/mobile views plus 80%, 100%, 125%, 150%, and
200% browser-equivalent widths. Tests also use a long pronunciation response
with multiple weak words and tips, since short reports do not expose the
reported overlap.

## Error Handling

No API contract changes are required. While a check is in progress, the
recording controls retain their existing pending state. When the API returns
an error, the existing error message remains in the controls area; no empty
result container is rendered below it.

## Acceptance Criteria

- `/app` desktop and mobile layouts remain independent at their platform
  boundary.
- A completed pronunciation result opens below the recording controls.
- The result does not overlap target text, controls, navigation, or sibling
  panels at supported widths and zoom-equivalent viewports.
- Long reports wrap within the available width and remain vertically
  scrollable.
- Existing pronunciation flows and the rest of `/app` continue to pass their
  smoke tests.

## Verification

Run focused Playwright pronunciation tests across the declared viewport
matrix, then run the full web smoke suite and the production web build.
