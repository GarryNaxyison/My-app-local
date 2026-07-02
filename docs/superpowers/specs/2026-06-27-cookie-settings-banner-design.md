# Cookie Settings Banner Design

## Goal

Add a real `Настроить` action to the NERIVA public-site cookie banner so users can choose optional cookie categories instead of only accepting all cookies or limiting the site to required cookies.

## Context

The current banner in `site-react/src/PublicSiteApp.tsx` stores one of two string values in `localStorage` under `poliglot-cookie-consent`:

- `necessary`
- `accepted`

The banner copy already says optional cookies are used only after consent, but the UI has no settings flow. This makes the banner legally and ergonomically weaker than the reference pattern shown by the user.

## Chosen Approach

Keep the banner compact by default and add a third action:

- `Только необходимые`
- `Настроить`
- `Принять все cookie`

Clicking `Настроить` expands an inline settings panel inside the same fixed banner. The panel should avoid a modal because the current banner is already small, fixed, and self-contained.

## Settings UI

The expanded settings panel contains:

- `Необходимые` as always active and not user-disableable.
- `Аналитические/маркетинговые` as an optional switch.
- `Сохранить выбор` to persist the current category state.

The fast actions remain available from the banner:

- `Только необходимые` saves required-only consent.
- `Принять все cookie` saves all categories enabled.

## Consent Storage

Continue supporting existing stored values so current users are not reprompted:

- `necessary` means only required cookies.
- `accepted` means all categories accepted.

For new settings saves, store a JSON value in `poliglot-cookie-consent` with a version and category booleans, for example:

```json
{"version":1,"necessary":true,"analyticsMarketing":true}
```

The banner should treat any valid existing consent value as dismissing the banner after reload.

## Styling

Use the existing `.cookie-consent-banner` visual language. Add small, stable controls for the settings rows without introducing a separate card inside the banner.

The three default actions must wrap cleanly on mobile. The expanded panel must keep buttons readable and reachable at narrow widths.

## Testing

Add or update focused Playwright coverage in `site-react/e2e/public-site.spec.ts`:

- The banner shows `Настроить`.
- `Настроить` opens settings without saving consent.
- Saving a selected optional state writes the JSON consent and hides the banner.
- Reload keeps the banner hidden after a saved settings choice.
- Existing `necessary` and `accepted` stored values still hide the banner.
- Existing quick actions still save and dismiss the banner.

## Non-Goals

- No new backend endpoint.
- No actual analytics or marketing script installation.
- No legal document rewrite beyond the existing banner wording.
- No modal settings center or account-level preferences page.

## Acceptance Criteria

- Public-site cookie banner has the third `Настроить` action.
- Settings flow provides a real optional category switch and save action.
- Required cookies are clearly always active.
- Existing stored `necessary` and `accepted` values remain compatible.
- Focused Playwright tests pass.
- Browser verification confirms desktop and mobile banner layouts do not overlap or hide actions.
