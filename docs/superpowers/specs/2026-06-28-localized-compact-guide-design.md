# Localized Compact Guide Design

## Goal

The web app guide must use real localized copy for all 35 interface languages and stop wasting modal space on a repeated visible header.

## Requirements

- Keep the existing five-page guide flow and page indicator.
- Remove the visible `DialogTitle` and `DialogDescription` text block from the guide modal body.
- Keep an accessible dialog title and description for screen readers.
- Replace non-Russian guide placeholder copy with localized guide text for every `appLocaleCodes` language.
- Keep branded/product terms stable where useful: `AI Tutor`, `Telegram`, `Free`, `Premium`, `Platinum`, `Listening`, `Review game`, and `Spelling`.
- Do not change backend behavior or plan/payment logic.

## Architecture

The UI change stays in `web-react/src/App.tsx` and `web-react/src/styles/app.css`. The visible modal body starts with the compact icon/page row and then immediately shows the active guide card.

The localization change stays in `web-react/src/lib/i18n.ts`. A dedicated guide copy map will override generated placeholder values for every locale, while the existing `appCopy` fallback and leak checks remain the final safety layer.

## Testing

- Add a failing Playwright/unit assertion that guide copy for non-Russian locales is meaningful and not the generated compact placeholder.
- Add a failing Playwright UI assertion that the guide dialog no longer renders visible title/description text.
- Run the targeted Playwright tests, TypeScript/Vite build, and the project web build before deploy.
