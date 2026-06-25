# Legal Compliance Documents Design

## Goal

Update Poliglot AI public legal pages, Telegram onboarding, and WebApp consent surfaces so the service publishes operator details, a privacy policy, terms, a separate user agreement, a separate personal-data consent, and an opt-in cookie notice.

## Context

The public site currently serves `/privacy.html` and `/terms.html` from `site-react`. Telegram shows only a privacy-policy prompt before onboarding. The WebApp auth and Telegram-link flows already require a personal-data checkbox, but they link only to the privacy policy.

The compliance scan flagged these issues:

- no separate personal-data consent document;
- no visible cookie opt-in banner;
- no public owner/operator details;
- foreign-server/localization risk;
- foreign words without translation risk.

## Operator Details

Use these details in public documents:

- Operator: self-employed person Cheban Denis Igorevich.
- Russian display name: Самозанятый Чебан Денис Игоревич.
- INN: 505017471160.
- Address for notices: Russia, Moscow Region, Shchyolkovo, Sirenevaya street, 9k1, apartment 9.
- Support email: supportpoliglotai@gmail.com.

## Design

### Public Site

Extend the legal page system from two pages to four pages:

- `/privacy.html`: privacy policy and personal-data processing policy.
- `/terms.html`: concise service terms.
- `/agreement.html`: separate user agreement/public offer adapted to Poliglot AI.
- `/consent.html`: separate consent to personal-data processing.

All four pages keep the existing legal-page shell and footer styling. Navigation and footer should expose all four document links.

### Cookies

Add a public-site cookie banner with explicit opt-in actions:

- Necessary cookies are described as required for session/security.
- Optional analytics cookies are enabled only after the user accepts them.
- The choice is stored locally and the banner can be dismissed.
- The banner links to `/privacy.html` and `/consent.html`.

Add the same banner behavior to the WebApp shell so users who open `/app/` directly see the same disclosure.

### Telegram

Keep the pre-onboarding privacy prompt, but update it to link to:

- privacy policy;
- personal-data consent;
- user agreement.

The continue callback remains `privacy_continue`.

### WebApp

Keep the existing explicit checkbox. Expand the text near auth and Telegram-link checkboxes to include links to:

- privacy policy;
- personal-data consent;
- user agreement.

Payloads should continue to send `privacy_consent` and `privacy_policy_url`; add consent and agreement URLs where the backend accepts extra JSON fields without breaking old behavior.

### Localization and Limits

The new legal documents are Russian-first. Existing language selector and generated localization assets remain in place for current privacy/terms behavior, but the legal minimum is the Russian text with clear document links and real operator details.

Server localization cannot be fully fixed by text changes. The policy should disclose current processing and note that localization/infrastructure may require a separate technical migration.

## Acceptance Criteria

- `/privacy.html`, `/terms.html`, `/agreement.html`, and `/consent.html` render without 404 in the public-site build.
- Legal pages show operator name, INN, address, support email, bot contact, and links to all legal documents.
- Cookie banner appears on public site and WebApp until a choice is saved.
- Telegram privacy prompt contains continue plus three legal-document URL buttons.
- WebApp auth and Telegram-link consent blocks link to privacy, consent, and agreement.
- Existing `/privacy.html` and `/terms.html` remain functional.
- Build and focused tests pass.
