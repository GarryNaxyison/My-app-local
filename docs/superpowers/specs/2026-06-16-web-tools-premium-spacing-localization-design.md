# Web Tools Premium Spacing Localization Design

## Goal

Tighten the desktop Tools view after the button-size reduction and finish the Russian copy in the Free/Premium plan cards.

## Requirements

- In the Tools view, the work area should start closer to the tool switcher so the previously empty top space is not left behind.
- The compact button sizes in the existing ribbon/tool switcher should stay unchanged.
- On Russian UI, the Free plan badge must read "Для начала" instead of "Попробовать маршрут".
- On Russian UI, the Free plan card must not show English fragments such as "Free", "AI Tutor guided lessons", "Listening и pronunciation", or "voice checks и photo tools".
- Product names "AI Tutor" and "Premium" may remain as branded names.
- The change must be covered by Playwright regression checks and verified visually before deployment.

## Design

Use a narrow CSS adjustment on `.tools-layout-v2` to align the chat/output panel and composer higher inside the Tools surface. Keep the responsive rules intact so mobile retains comfortable spacing.

Use the existing `appCopy` and fallback plan helpers for copy ownership. Update Russian plan labels/features in the localized copy layer and fallback plan definitions so both API-provided plans and local fallback plans render the same Russian strings.

## Verification

- Add a Playwright regression covering the Russian Free card copy and desktop Tools vertical position.
- Run the focused Playwright test before and after implementation for RED/GREEN.
- Run `npm --prefix web-react run build`.
- Run `npm --prefix web-react run e2e`.
- Run the encoding artifact check before deploy if available.
