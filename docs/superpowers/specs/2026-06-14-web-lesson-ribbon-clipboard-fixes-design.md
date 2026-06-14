# Web Lesson Ribbon Clipboard Fixes Design

## Goal

Fix the three reported web regressions without changing the navigation model: the mobile Lesson "New lesson" action must stay reachable, desktop ribbon arrows must remain visible and clickable when menu chips expand, and Tools image mode must accept images pasted from the clipboard even when paste lands outside the upload box.

## Scope

- Mobile Lesson keeps the new-lesson action inside the viewport and above the fixed bottom navigation.
- Desktop function-ribbon arrow columns reserve enough space so expanded menu chips cannot cover the left/right arrow buttons.
- Tools image mode handles pasted image files from the focused tools area, not only from the small upload control itself.

## Non-Goals

- No backend API changes.
- No new navigation destinations.
- No broad V2 visual redesign.

## Acceptance Checks

- Playwright has focused regressions for the mobile lesson button viewport position, desktop ribbon arrow hit targets during chip expansion, and image-tool clipboard paste from the tools composer area.
- `npm --prefix web-react run build` succeeds.
- Targeted Playwright checks for the three regressions pass on desktop and mobile projects where applicable.
